package watchers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	clusterctlv1 "sigs.k8s.io/cluster-api/cmd/clusterctl/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi"
	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/day2ops"
	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/fetchers"
)

// machineSetGVR identifies the MachineSet resource — first-class CAPI type sitting between
// MachineDeployment and Machine, watched only by the Day-2 Ops aggregator to detect stalled
// rollouts (research.md R5), not exposed as its own list page.
var machineSetGVR = schema.GroupVersionResource{Group: "cluster.x-k8s.io", Version: "v1beta1", Resource: "machinesets"}

// machineHealthCheckGVR identifies the MachineHealthCheck resource — a first-class CAPI
// remediation type, watched only by the Day-2 Ops aggregator to classify self-healing vs
// needs-investigation severity (FR-012, FR-013, research.md R7), not exposed as its own list page.
var machineHealthCheckGVR = schema.GroupVersionResource{Group: "cluster.x-k8s.io", Version: "v1beta1", Resource: "machinehealthchecks"}

// versionSkewCRDs are the CRDs checked for stored-but-no-longer-served versions (research.md R6).
var versionSkewCRDs = []string{"machines.cluster.x-k8s.io", "clusters.cluster.x-k8s.io"}

// day2opsCoreGVRs are always watched: first-class CAPI core types present on every management
// cluster regardless of which infrastructure provider(s) are installed.
var day2opsCoreGVRs = []schema.GroupVersionResource{
	clusterGVR, machineDeploymentGVR, machineGVR, machineSetGVR, machineHealthCheckGVR,
}

// day2opsScheme registers just enough types to read the clusterctl provider inventory (used only
// to detect which infrastructure providers are installed, mirroring the existing pattern in
// clusterapi.GenerateInfrastructureCapability). Kept local to this package (rather than reusing
// system.Scheme) since handlers/system already imports this package — importing it back would be
// a cycle.
var day2opsScheme = func() *runtime.Scheme {
	s := runtime.NewScheme()
	_ = clusterctlv1.AddToScheme(s)
	return s
}()

// day2opsWatchedGVRs returns the core GVRs plus only the provider-infra GVRs for providers
// actually installed in this environment (research.md; a provider's CRD, e.g. VSphereMachine,
// does not exist at all when that provider isn't installed — attempting to watch it fails
// immediately and must not be treated as fatal to the whole dashboard).
func day2opsWatchedGVRs(ctx context.Context) []schema.GroupVersionResource {
	gvrs := append([]schema.GroupVersionResource{}, day2opsCoreGVRs...)

	cli, err := clusterapi.NewClientWithScheme(ctx, day2opsScheme)
	if err != nil {
		return gvrs
	}
	capability, err := clusterapi.GenerateInfrastructureCapability(ctx, cli)
	if err != nil {
		return gvrs
	}
	if capability.Docker.Installed {
		gvrs = append(gvrs, dockerMachineGVR)
	}
	if capability.VSphere.Installed {
		gvrs = append(gvrs, machineInfraGVR)
	}
	return gvrs
}

// day2opsEvent is one fanned-in watch event, tagged with which GVR it came from.
type day2opsEvent struct {
	gvr   schema.GroupVersionResource
	event watch.Event
}

// day2opsStore holds the latest known state per watched kind, keyed by "namespace/name".
type day2opsStore struct {
	mu                  sync.Mutex
	clusters            map[string]clusterv1.Cluster
	machineDeployments  map[string]clusterv1.MachineDeployment
	machines            map[string]clusterv1.Machine
	machineSets         map[string]clusterv1.MachineSet
	machineHealthChecks map[string]clusterv1.MachineHealthCheck
	// providerResources is keyed by the provider-infra object's OWN "namespace/name" (as referenced
	// by a Machine's Spec.InfrastructureRef), not the owning Machine's identity.
	providerResources map[string]day2ops.ProviderResourceStatus
}

func newDay2opsStore() *day2opsStore {
	return &day2opsStore{
		clusters:            map[string]clusterv1.Cluster{},
		machineDeployments:  map[string]clusterv1.MachineDeployment{},
		machines:            map[string]clusterv1.Machine{},
		machineSets:         map[string]clusterv1.MachineSet{},
		machineHealthChecks: map[string]clusterv1.MachineHealthCheck{},
		providerResources:   map[string]day2ops.ProviderResourceStatus{},
	}
}

func objectKey(namespace, name string) string { return namespace + "/" + name }

// apply upserts or removes the fanned-in event's object from the store, based on its GVR.
func (s *day2opsStore) apply(evt day2opsEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	unstructuredObj, ok := evt.event.Object.(*unstructured.Unstructured)
	if !ok {
		return fmt.Errorf("unexpected event object type %T", evt.event.Object)
	}

	switch evt.gvr {
	case clusterGVR:
		var cl clusterv1.Cluster
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredObj.UnstructuredContent(), &cl); err != nil {
			return err
		}
		key := objectKey(cl.Namespace, cl.Name)
		if evt.event.Type == watch.Deleted {
			delete(s.clusters, key)
		} else {
			s.clusters[key] = cl
		}
	case machineDeploymentGVR:
		var md clusterv1.MachineDeployment
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredObj.UnstructuredContent(), &md); err != nil {
			return err
		}
		key := objectKey(md.Namespace, md.Name)
		if evt.event.Type == watch.Deleted {
			delete(s.machineDeployments, key)
		} else {
			s.machineDeployments[key] = md
		}
	case machineGVR:
		var m clusterv1.Machine
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredObj.UnstructuredContent(), &m); err != nil {
			return err
		}
		key := objectKey(m.Namespace, m.Name)
		if evt.event.Type == watch.Deleted {
			delete(s.machines, key)
		} else {
			s.machines[key] = m
		}
	case machineSetGVR:
		var ms clusterv1.MachineSet
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredObj.UnstructuredContent(), &ms); err != nil {
			return err
		}
		key := objectKey(ms.Namespace, ms.Name)
		if evt.event.Type == watch.Deleted {
			delete(s.machineSets, key)
		} else {
			s.machineSets[key] = ms
		}
	case machineHealthCheckGVR:
		var mhc clusterv1.MachineHealthCheck
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredObj.UnstructuredContent(), &mhc); err != nil {
			return err
		}
		key := objectKey(mhc.Namespace, mhc.Name)
		if evt.event.Type == watch.Deleted {
			delete(s.machineHealthChecks, key)
		} else {
			s.machineHealthChecks[key] = mhc
		}
	case machineInfraGVR, dockerMachineGVR:
		key := objectKey(unstructuredObj.GetNamespace(), unstructuredObj.GetName())
		if evt.event.Type == watch.Deleted {
			delete(s.providerResources, key)
		} else {
			s.providerResources[key] = day2ops.ExtractProviderResourceStatus(unstructuredObj)
		}
	}
	return nil
}

// providerResourceFor looks up the provider-infra status for a Machine via its
// Spec.InfrastructureRef, returning nil when no matching object has been observed yet.
func (s *day2opsStore) providerResourceFor(m clusterv1.Machine) *day2ops.ProviderResourceStatus {
	s.mu.Lock()
	defer s.mu.Unlock()
	if m.Spec.InfrastructureRef.Name == "" {
		return nil
	}
	namespace := m.Spec.InfrastructureRef.Namespace
	if namespace == "" {
		namespace = m.Namespace
	}
	status, ok := s.providerResources[objectKey(namespace, m.Spec.InfrastructureRef.Name)]
	if !ok {
		return nil
	}
	return &status
}

// snapshot returns copies of the current known objects, safe to compute against without holding
// the store's lock.
func (s *day2opsStore) snapshot() (clusters []clusterv1.Cluster, machineDeployments []clusterv1.MachineDeployment, machines []clusterv1.Machine) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, cl := range s.clusters {
		clusters = append(clusters, cl)
	}
	for _, md := range s.machineDeployments {
		machineDeployments = append(machineDeployments, md)
	}
	for _, m := range s.machines {
		machines = append(machines, m)
	}
	return clusters, machineDeployments, machines
}

// machineSetsFor returns the MachineSets owned by a MachineDeployment, matching on the standard
// CAPI ownership label (the same one the MachineSet/MachineDeployment controllers themselves set).
func (s *day2opsStore) machineSetsFor(md clusterv1.MachineDeployment) []clusterv1.MachineSet {
	s.mu.Lock()
	defer s.mu.Unlock()
	var owned []clusterv1.MachineSet
	for _, ms := range s.machineSets {
		if ms.Namespace == md.Namespace && ms.Labels["cluster.x-k8s.io/deployment-name"] == md.Name {
			owned = append(owned, ms)
		}
	}
	return owned
}

// machineHealthCheckSnapshot returns a copy of every currently-known MachineHealthCheck.
func (s *day2opsStore) machineHealthCheckSnapshot() []clusterv1.MachineHealthCheck {
	s.mu.Lock()
	defer s.mu.Unlock()
	mhcs := make([]clusterv1.MachineHealthCheck, 0, len(s.machineHealthChecks))
	for _, mhc := range s.machineHealthChecks {
		mhcs = append(mhcs, mhc)
	}
	return mhcs
}

// providerResourceSnapshot returns a copy of every currently-known provider-infra status, for
// checks (like drift) that scan all of them rather than looking one up by owning Machine.
func (s *day2opsStore) providerResourceSnapshot() []day2ops.ProviderResourceStatus {
	s.mu.Lock()
	defer s.mu.Unlock()
	statuses := make([]day2ops.ProviderResourceStatus, 0, len(s.providerResources))
	for _, status := range s.providerResources {
		statuses = append(statuses, status)
	}
	return statuses
}

// assembleData recomputes the full Day2Ops payload from the store's current snapshot. Debugging
// paths are computed only for currently-unhealthy Machines (FR-004), with each layer's evidence
// capped to one line for the live WS push — the full, uncapped path is available on demand via
// GET /api/day2ops/detail (contracts/day2ops-ws-event.md consumer contract). apiext may be nil
// (e.g. in tests), in which case the version-skew check is skipped rather than erroring.
func assembleData(ctx context.Context, dyn dynamic.Interface, apiext *apiextensionsclientset.Clientset, store *day2opsStore, sourceUnavailable bool) day2ops.Data {
	clusters, machineDeployments, machines := store.snapshot()

	debugPaths := make([]day2ops.DebugPath, 0)
	for _, m := range machines {
		if !day2ops.MachineFailed(m) {
			continue
		}
		objectRef := day2ops.ObjectRef{
			Group: machineGVR.Group, Version: machineGVR.Version, Resource: machineGVR.Resource,
			Namespace: m.Namespace, Name: m.Name,
		}
		path := day2ops.ComputeMachineDebugPath(objectRef, m, store.providerResourceFor(m), machineControllerEvents(ctx, dyn, m))
		debugPaths = append(debugPaths, capDebugPathEvidence(path))
	}

	certRisks, caMissingSeverities := clusterCertRisksAndSeverities(ctx, dyn, clusters)
	risks := make([]day2ops.RiskWarning, 0)
	risks = append(risks, certRisks...)
	risks = append(risks, stalledRolloutRisks(store, machineDeployments)...)
	risks = append(risks, driftRisks(store)...)
	risks = append(risks, versionSkewRisks(ctx, apiext)...)

	severities := make([]day2ops.FailureSeverity, 0)
	if severity := day2ops.ComputeManagementClusterSeverity(sourceUnavailable); severity != nil {
		severities = append(severities, *severity)
	}
	severities = append(severities, caMissingSeverities...)
	severities = append(severities, machineHealthCheckSeverities(store)...)
	severities = append(severities, providerControllerSeverities(ctx, dyn)...)

	return day2ops.Data{
		Rollups:           day2ops.ComputeRollups(clusters, machineDeployments, machines),
		DebugPaths:        debugPaths,
		Risks:             risks,
		Severities:        severities,
		SourceUnavailable: sourceUnavailable,
	}
}

// clusterCertRisksAndSeverities fetches each cluster's cert Secrets once and derives both the
// cert-expiry risk warnings (FR-008) and the CA-secret-missing severity (FR-016) from the same
// fetch, rather than reading the Secrets twice.
func clusterCertRisksAndSeverities(ctx context.Context, dyn dynamic.Interface, clusters []clusterv1.Cluster) ([]day2ops.RiskWarning, []day2ops.FailureSeverity) {
	risks := make([]day2ops.RiskWarning, 0)
	severities := make([]day2ops.FailureSeverity, 0)
	if dyn == nil {
		return risks, severities
	}
	for _, cl := range clusters {
		objectRef := day2ops.ObjectRef{
			Group: clusterGVR.Group, Version: clusterGVR.Version, Resource: clusterGVR.Resource,
			Namespace: cl.Namespace, Name: cl.Name,
		}
		expiries, err := fetchers.FetchClusterCertExpiries(ctx, dyn, cl.Namespace, cl.Name)
		risks = append(risks, day2ops.ComputeClusterCertRisks(objectRef, expiries, err, time.Now(), day2ops.DefaultCertExpiryWarningWindow)...)

		if err == nil {
			caFound := false
			for _, e := range expiries {
				if e.SecretName == cl.Name+"-ca" {
					caFound = true
					break
				}
			}
			if severity := day2ops.ComputeCASecretMissingSeverity(objectRef, caFound); severity != nil {
				severities = append(severities, *severity)
			}
		}
	}
	return risks, severities
}

// machineHealthCheckSeverities classifies every known MachineHealthCheck's remediation state
// (FR-012, FR-013).
func machineHealthCheckSeverities(store *day2opsStore) []day2ops.FailureSeverity {
	severities := make([]day2ops.FailureSeverity, 0)
	for _, mhc := range store.machineHealthCheckSnapshot() {
		objectRef := day2ops.ObjectRef{
			Group: machineHealthCheckGVR.Group, Version: machineHealthCheckGVR.Version, Resource: machineHealthCheckGVR.Resource,
			Namespace: mhc.Namespace, Name: mhc.Name,
		}
		status := day2ops.MachineHealthCheckStatus{
			Name: mhc.Name, ExpectedMachines: mhc.Status.ExpectedMachines,
			CurrentHealthy: mhc.Status.CurrentHealthy, RemediationsAllowed: mhc.Status.RemediationsAllowed,
		}
		if severity := day2ops.ComputeMachineHealthCheckSeverity(objectRef, status); severity != nil {
			severities = append(severities, *severity)
		}
	}
	return severities
}

// providerControllerSeverities scans the well-known controller namespaces for not-ready Pods
// (FR-014). A namespace that doesn't exist (provider not installed) yields no error and no
// severities — it's simply skipped.
func providerControllerSeverities(ctx context.Context, dyn dynamic.Interface) []day2ops.FailureSeverity {
	severities := make([]day2ops.FailureSeverity, 0)
	if dyn == nil {
		return severities
	}
	for _, namespace := range fetchers.ControllerNamespaces {
		pods, err := fetchers.FetchControllerPodStatuses(ctx, dyn, namespace)
		if err != nil {
			continue
		}
		for _, pod := range pods {
			if severity := day2ops.ComputeProviderControllerSeverity(pod); severity != nil {
				severities = append(severities, *severity)
			}
		}
	}
	return severities
}

func stalledRolloutRisks(store *day2opsStore, machineDeployments []clusterv1.MachineDeployment) []day2ops.RiskWarning {
	risks := make([]day2ops.RiskWarning, 0)
	for _, md := range machineDeployments {
		objectRef := day2ops.ObjectRef{
			Group: machineDeploymentGVR.Group, Version: machineDeploymentGVR.Version, Resource: machineDeploymentGVR.Resource,
			Namespace: md.Namespace, Name: md.Name,
		}
		machineSets := store.machineSetsFor(md)
		if risk := day2ops.ComputeStalledRolloutRisk(objectRef, machineSets, nil, time.Now()); risk != nil {
			risks = append(risks, *risk)
		}
	}
	return risks
}

func driftRisks(store *day2opsStore) []day2ops.RiskWarning {
	risks := make([]day2ops.RiskWarning, 0)
	for _, provider := range store.providerResourceSnapshot() {
		objectRef := day2ops.ObjectRef{Name: provider.Name}
		if risk := day2ops.ComputeDriftRisk(objectRef, provider); risk != nil {
			risks = append(risks, *risk)
		}
	}
	return risks
}

func versionSkewRisks(ctx context.Context, apiext *apiextensionsclientset.Clientset) []day2ops.RiskWarning {
	risks := make([]day2ops.RiskWarning, 0)
	if apiext == nil {
		return risks
	}
	for _, crdName := range versionSkewCRDs {
		info, err := fetchers.FetchCRDVersionInfo(ctx, apiext, crdName)
		if err != nil {
			continue
		}
		if risk := day2ops.ComputeVersionSkewRisk(day2ops.ObjectRef{Name: crdName}, info); risk != nil {
			risks = append(risks, *risk)
		}
	}
	return risks
}

// machineControllerEvents fetches recent Events for a Machine only when its higher-priority
// debugging layers (conditions, phase) give no explanation — avoiding an extra API call for the
// common case where conditions/phase already answer the question (research.md R2, FR-007).
func machineControllerEvents(ctx context.Context, dyn dynamic.Interface, m clusterv1.Machine) []string {
	if !day2ops.ShouldFetchControllerActivityEvents(m) {
		return nil
	}
	events, err := fetchers.FetchInvolvedObjectEvents(ctx, dyn, m.Namespace, m.Name, "Machine")
	if err != nil {
		return nil
	}
	return events
}

// capDebugPathEvidence truncates every layer's evidence to at most one line, per the live WS
// payload's size contract; the uncapped version is only ever returned by the REST detail endpoint.
func capDebugPathEvidence(path day2ops.DebugPath) day2ops.DebugPath {
	capped := make([]day2ops.DebugLayer, len(path.Layers))
	for i, l := range path.Layers {
		if len(l.Evidence) > 1 {
			l.Evidence = l.Evidence[:1]
		}
		capped[i] = l
	}
	path.Layers = capped
	return path
}

// WatchDay2Ops opens the Day-2 Ops dashboard's own set of Kubernetes watches (Cluster,
// MachineDeployment, Machine, and — as later user stories add them — MachineSet and
// MachineHealthCheck), fans their events into one loop, recomputes the consolidated Day2OpsEvent
// on every change, and streams it to this connection. Unlike every other resource watcher in this
// package (which relay exactly one GVR 1:1 into one connection via WatchResourceViaWebSocket),
// this one fans in several GVRs because the dashboard's rollup is itself a cross-resource view
// (research.md R9 in specs/006-day2-ops-dashboard/).
func WatchDay2Ops(ctx context.Context, conn *websocket.Conn, objType string) error {
	dynamicClient, err := clusterapi.NewDynamicClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create dynamic client: %w", err)
	}
	// Best-effort: the version-skew check is skipped (not the whole connection failed) if this
	// client can't be constructed, since it's a secondary risk check, not core functionality.
	apiextClient, _ := clusterapi.NewAPIExtensionsClient(ctx)

	watchedGVRs := day2opsWatchedGVRs(ctx)
	events := make(chan day2opsEvent)
	watchErrs := make(chan error, len(watchedGVRs))
	watchCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	logger := log.FromContext(ctx)
	var wg sync.WaitGroup
	opened := 0
	for _, gvr := range watchedGVRs {
		watcher, err := dynamicClient.Resource(gvr).Namespace("").Watch(watchCtx, metav1.ListOptions{ResourceVersion: "0"})
		if err != nil {
			// Not every GVR is guaranteed to exist on every management cluster (e.g. the
			// VSphereMachine CRD is absent on a Docker-only install) - skip it rather than
			// aborting the whole dashboard over one optional/provider-specific resource kind.
			logger.Error(err, "skipping day2ops watch for unavailable resource", "gvr", gvr)
			continue
		}
		opened++
		wg.Add(1)
		go func(gvr schema.GroupVersionResource, watcher watch.Interface) {
			defer wg.Done()
			defer watcher.Stop()
			for event := range watcher.ResultChan() {
				select {
				case events <- day2opsEvent{gvr: gvr, event: event}:
				case <-watchCtx.Done():
					return
				}
			}
			// ResultChan closed without ctx being cancelled: the underlying watch broke
			// (e.g. API server became unreachable) rather than this handler shutting down.
			select {
			case watchErrs <- fmt.Errorf("watch for %v closed unexpectedly", gvr):
			default:
			}
		}(gvr, watcher)
	}
	if opened == 0 {
		cancel()
		return fmt.Errorf("no day2ops watches could be opened against the management cluster")
	}

	store := newDay2opsStore()

	// Seed an initial snapshot immediately, so the dashboard doesn't wait for the first change.
	if err = conn.WriteJSON(EventResponse{Type: "MODIFIED", Event: objType, Data: assembleData(ctx, dynamicClient, apiextClient, store, false)}); err != nil {
		cancel()
		wg.Wait()
		return fmt.Errorf("failed to write initial day2ops snapshot: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			cancel()
			wg.Wait()
			return ctx.Err()
		case werr := <-watchErrs:
			// A contributing watch died: tell the client the data source is unavailable, then
			// close the connection so the frontend's own reconnect/backoff takes over (FR-017).
			_ = conn.WriteJSON(EventResponse{Type: "MODIFIED", Event: objType, Data: assembleData(ctx, dynamicClient, apiextClient, store, true)})
			cancel()
			wg.Wait()
			return werr
		case evt := <-events:
			if err = store.apply(evt); err != nil {
				cancel()
				wg.Wait()
				return fmt.Errorf("failed to apply day2ops event: %w", err)
			}
			if err = conn.WriteJSON(EventResponse{Type: "MODIFIED", Event: objType, Data: assembleData(ctx, dynamicClient, apiextClient, store, false)}); err != nil {
				cancel()
				wg.Wait()
				return fmt.Errorf("failed to write day2ops event: %w", err)
			}
		}
	}
}
