package kured_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trustyou/kured-alert-silencer/pkg/kured"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestExtractNodeIDsFromAnnotation(t *testing.T) {
	const KuredNodeLockAnnotation string = "kured.dev/kured-node-lock"

	fixedTime := time.Date(2024, time.May, 31, 6, 33, 0, 0, time.UTC)
	timeProvider := func() time.Time {
		return fixedTime
	}
	silenceDuration, err := time.ParseDuration("1h")
	if err != nil {
		assert.NoError(t, err)
	}

	tests := []struct {
		name            string
		annotationValue string
		want            []kured.SilenceNode
	}{
		{
			name:            "single node ID",
			annotationValue: `{"nodeID":"kind-control-plane2","metadata":{"unschedulable":false},"created":"2024-05-31T06:31:37.441623522Z","TTL":0}`,
			want: []kured.SilenceNode{
				{
					"kind-control-plane2", time.Date(2024, time.May, 31, 7, 31, 37, 441623522, time.UTC),
				},
			},
		},
		{
			name:            "single node ID not silenced",
			annotationValue: `{"nodeID":"kind-control-plane2","metadata":{"unschedulable":false},"created":"2024-05-30T00:00:00.000000000Z","TTL":0}`,
			want:            []kured.SilenceNode{},
		},
		{
			name:            "multiple node IDs",
			annotationValue: `{"maxOwners":2,"locks":[{"nodeID":"kind-worker2","metadata":{"unschedulable":false},"created":"2024-05-31T06:31:32.735905893Z","TTL":0},{"nodeID":"kind-control-plane","metadata":{"unschedulable":false},"created":"2024-05-31T06:31:49.868231413Z","TTL":0}]}`,
			want: []kured.SilenceNode{
				{
					"kind-worker2", time.Date(2024, time.May, 31, 7, 31, 32, 735905893, time.UTC),
				},
				{
					"kind-control-plane", time.Date(2024, time.May, 31, 7, 31, 49, 868231413, time.UTC),
				},
			},
		},
		{
			name:            "multiple node IDs just one silenced",
			annotationValue: `{"maxOwners":2,"locks":[{"nodeID":"kind-worker2","metadata":{"unschedulable":false},"created":"2024-05-30T00:00:32.735905893Z","TTL":0},{"nodeID":"kind-control-plane","metadata":{"unschedulable":false},"created":"2024-05-31T06:31:49.868231413Z","TTL":0}]}`,
			want: []kured.SilenceNode{
				{
					"kind-control-plane", time.Date(2024, time.May, 31, 7, 31, 49, 868231413, time.UTC),
				},
			},
		},
		{
			name:            "manual lock",
			annotationValue: `{"nodeID":"manual"}`,
			want:            []kured.SilenceNode{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &v1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						KuredNodeLockAnnotation: tt.annotationValue,
					},
				},
			}

			result, err := kured.ExtractNodeIDsFromAnnotation(ds, KuredNodeLockAnnotation, silenceDuration, timeProvider)
			require.NoError(t, err)
			require.Equal(t, tt.want, result)

		})
	}
}

func TestExtractNodeIDsFromAnnotationNoAnnotation(t *testing.T) {
	ds := &v1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{},
		},
	}

	fixedTime := time.Date(2024, time.May, 6, 31, 6, 33, 0, time.UTC)
	timeProvider := func() time.Time {
		return fixedTime
	}
	silenceDuration, err := time.ParseDuration("1h")
	if err != nil {
		assert.NoError(t, err)
	}

	result, err := kured.ExtractNodeIDsFromAnnotation(ds, "weave.works/kured-node-lock", silenceDuration, timeProvider)
	require.NoError(t, err)
	require.Equal(t, []kured.SilenceNode{}, result)

}
