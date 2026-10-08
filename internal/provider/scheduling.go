package provider

import (
	commonv1alpha1 "github.com/openeverest/openeverest/v2/api/common/v1alpha1"
	"github.com/openeverest/openeverest/v2/provider-runtime/controller"

	"github.com/openeverest/provider-milvus/internal/milvusapi"
)

// applySchedulingPolicy places one component's pods. The operator sets no
// placement of its own, so an unset policy leaves the scheduler's defaults.
func applySchedulingPolicy(spec *milvusapi.ComponentSpec, policy *commonv1alpha1.SchedulingPolicy) {
	spec.TopologySpreadConstraints = controller.TopologySpreadConstraints(policy, spec.PodLabels)
	if policy == nil {
		return
	}
	spec.SchedulerName = policy.SchedulerName
	spec.NodeSelector = policy.NodeSelector
	spec.Affinity = policy.Affinity
	spec.Tolerations = policy.Tolerations
}
