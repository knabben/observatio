package fetchers

import (
	"context"
	"fmt"
	"sort"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// eventGVR identifies the core Kubernetes Event resource, read generically via the dynamic
// client (research.md R2 in specs/006-day2-ops-dashboard/: Layer 4 "controller reconciliation
// activity" is sourced from Events, not raw controller-pod log tailing).
var eventGVR = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "events"}

// maxDay2OpsEvents bounds how many recent event messages are surfaced per object, keeping the
// evidence list short and readable rather than an unbounded dump.
const maxDay2OpsEvents = 5

// FetchInvolvedObjectEvents returns up to maxDay2OpsEvents recent event messages
// ("Reason: message") for the object identified by namespace/name/kind, most recent first. Used
// only as the debugging path's controller-activity layer, and only when the higher layers
// (conditions, phase, provider resource) are all inconclusive (FR-007).
func FetchInvolvedObjectEvents(ctx context.Context, dyn dynamic.Interface, namespace, name, kind string) ([]string, error) {
	fieldSelector := fmt.Sprintf("involvedObject.name=%s,involvedObject.namespace=%s,involvedObject.kind=%s", name, namespace, kind)
	list, err := dyn.Resource(eventGVR).Namespace(namespace).List(ctx, metav1.ListOptions{FieldSelector: fieldSelector})
	if err != nil {
		return nil, err
	}

	type eventEntry struct {
		lastSeen string
		message  string
	}
	entries := make([]eventEntry, 0, len(list.Items))
	for _, item := range list.Items {
		message, _, _ := unstructured.NestedString(item.Object, "message")
		reason, _, _ := unstructured.NestedString(item.Object, "reason")
		eventType, _, _ := unstructured.NestedString(item.Object, "type")
		lastSeen, _, _ := unstructured.NestedString(item.Object, "lastTimestamp")
		if message == "" {
			continue
		}
		entries = append(entries, eventEntry{
			lastSeen: lastSeen,
			message:  fmt.Sprintf("%s %s: %s", eventType, reason, message),
		})
	}

	sort.Slice(entries, func(i, j int) bool { return entries[i].lastSeen > entries[j].lastSeen })
	if len(entries) > maxDay2OpsEvents {
		entries = entries[:maxDay2OpsEvents]
	}

	messages := make([]string, len(entries))
	for i, e := range entries {
		messages[i] = e.message
	}
	return messages, nil
}
