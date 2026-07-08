package fetchers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	dynamicfake "k8s.io/client-go/dynamic/fake"
)

func Test_FindControllerPodName_ResolvesMatchingPod(t *testing.T) {
	testScheme := runtime.NewScheme()
	require.NoError(t, appsv1.AddToScheme(testScheme))
	require.NoError(t, corev1.AddToScheme(testScheme))

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: "capd-controller-manager", Namespace: "capd-system"},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"control-plane": "controller-manager"}},
		},
	}
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "capd-controller-manager-abc123", Namespace: "capd-system",
			Labels: map[string]string{"control-plane": "controller-manager"},
		},
	}

	gvrToListKind := map[schema.GroupVersionResource]string{
		deploymentGVR: "DeploymentList",
	}
	dyn := dynamicfake.NewSimpleDynamicClientWithCustomListKinds(testScheme, gvrToListKind, deployment, pod)

	podName, err := FindControllerPodName(context.Background(), dyn, "capd-system", "capd-controller-manager")
	require.NoError(t, err)
	assert.Equal(t, "capd-controller-manager-abc123", podName)
}

func Test_FindControllerPodName_NoMatchingPod(t *testing.T) {
	testScheme := runtime.NewScheme()
	require.NoError(t, appsv1.AddToScheme(testScheme))
	require.NoError(t, corev1.AddToScheme(testScheme))

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: "capd-controller-manager", Namespace: "capd-system"},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"control-plane": "controller-manager"}},
		},
	}

	gvrToListKind := map[schema.GroupVersionResource]string{
		deploymentGVR: "DeploymentList",
	}
	dyn := dynamicfake.NewSimpleDynamicClientWithCustomListKinds(testScheme, gvrToListKind, deployment)

	podName, err := FindControllerPodName(context.Background(), dyn, "capd-system", "capd-controller-manager")
	require.NoError(t, err)
	assert.Empty(t, podName)
}
