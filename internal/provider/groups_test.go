package provider

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"

	corev1alpha1 "github.com/openeverest/openeverest/v2/api/core/v1alpha1"

	"github.com/openeverest/provider-milvus/definition/components"
	"github.com/openeverest/provider-milvus/internal/common"
)

func groupParams(t *testing.T, groups ...components.DeploymentGroup) *runtime.RawExtension {
	t.Helper()
	raw, err := json.Marshal(components.MilvusParameters{Groups: groups})
	require.NoError(t, err)
	return &runtime.RawExtension{Raw: raw}
}

func TestBuildMilvusSpecDeploymentGroups(t *testing.T) {
	gpuTaint := []corev1.Toleration{{Key: "nvidia.com/gpu", Operator: corev1.TolerationOpExists}}
	c := newTestContext(t, corev1alpha1.InstanceSpec{
		Topology: &corev1alpha1.TopologySpec{Type: "cluster"},
		Components: map[string]corev1alpha1.ComponentSpec{
			common.ComponentQueryNode: {Parameters: groupParams(t,
				components.DeploymentGroup{
					Name: "l40s", Replicas: 4,
					NodeSelector: map[string]string{"nvidia.com/gpu.product": "NVIDIA-L40S"},
					Tolerations:  gpuTaint,
					Annotations:  map[string]string{"k8s.v1.cni.cncf.io/networks": "sriov-l40s"},
				},
				components.DeploymentGroup{
					Name: "h200", Replicas: 2,
					NodeSelector: map[string]string{"nvidia.com/gpu.product": "NVIDIA-H200"},
					Affinity: map[string]any{"podAntiAffinity": map[string]any{
						"preferredDuringSchedulingIgnoredDuringExecution": []any{map[string]any{
							"weight":          100,
							"podAffinityTerm": map[string]any{"topologyKey": "kubernetes.io/hostname"},
						}},
					}},
				},
			)},
		},
	})
	spec, err := BuildMilvusSpec(c)
	require.NoError(t, err)

	queryNode := spec.Com.QueryNode
	assert.Equal(t, ptr.To(int32(6)), queryNode.Replicas)
	require.Len(t, queryNode.Groups, 2)
	l40s := queryNode.Groups[0]
	assert.Equal(t, "l40s", l40s.Name)
	assert.Equal(t, ptr.To(int32(4)), l40s.Replicas)
	assert.Equal(t, &map[string]string{"nvidia.com/gpu.product": "NVIDIA-L40S"}, l40s.NodeSelector)
	assert.Equal(t, &gpuTaint, l40s.Tolerations)
	assert.Equal(t, "sriov-l40s", l40s.Annotations["k8s.v1.cni.cncf.io/networks"])
	assert.Nil(t, queryNode.Groups[1].Tolerations, "unset placement inherits the component's")
	assert.Nil(t, l40s.Affinity)
	require.NotNil(t, queryNode.Groups[1].Affinity)
	assert.Equal(t, "kubernetes.io/hostname",
		queryNode.Groups[1].Affinity.PodAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution[0].PodAffinityTerm.TopologyKey)

	assert.Empty(t, spec.Com.DataNode.Groups)
}

func TestValidateDeploymentGroups(t *testing.T) {
	tests := map[string]struct {
		component string
		replicas  *int32
		groups    []components.DeploymentGroup
		wantErr   string
	}{
		"replicas match the groups": {
			component: common.ComponentQueryNode, replicas: ptr.To(int32(3)),
			groups: []components.DeploymentGroup{{Name: "a", Replicas: 1}, {Name: "b", Replicas: 2}},
		},
		"replicas differ from the groups": {
			component: common.ComponentQueryNode, replicas: ptr.To(int32(1)),
			groups:  []components.DeploymentGroup{{Name: "a", Replicas: 1}, {Name: "b", Replicas: 2}},
			wantErr: "must equal the sum of its groups' replicas (3)",
		},
		"invalid name": {
			component: common.ComponentDataNode,
			groups:    []components.DeploymentGroup{{Name: "L40S", Replicas: 1}},
			wantErr:   `groups[0].name "L40S" is invalid`,
		},
		"duplicate name": {
			component: common.ComponentProxy,
			groups:    []components.DeploymentGroup{{Name: "a", Replicas: 1}, {Name: "a", Replicas: 1}},
			wantErr:   `duplicate group "a"`,
		},
		"negative replicas": {
			component: common.ComponentStreaming,
			groups:    []components.DeploymentGroup{{Name: "a", Replicas: -1}},
			wantErr:   "replicas must be >= 0",
		},
		"component without groups support": {
			component: common.ComponentMixCoord,
			groups:    []components.DeploymentGroup{{Name: "a", Replicas: 1}},
			wantErr:   `component "mixCoord" does not support groups`,
		},
		"invalid affinity": {
			component: common.ComponentQueryNode,
			groups:    []components.DeploymentGroup{{Name: "a", Replicas: 1, Affinity: map[string]any{"nodeAffinty": map[string]any{}}}},
			wantErr:   `group "a" affinity is invalid`,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			c := newTestContext(t, corev1alpha1.InstanceSpec{
				Topology: &corev1alpha1.TopologySpec{Type: "cluster"},
				Components: map[string]corev1alpha1.ComponentSpec{
					tt.component: {Replicas: tt.replicas, Parameters: groupParams(t, tt.groups...)},
				},
			})
			err := validateInstance(c)
			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}
