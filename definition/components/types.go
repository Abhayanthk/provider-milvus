// Package components contains parameter types for provider component types.
//
// Each struct here corresponds to a component type defined in versions.yaml
// and is converted to an OpenAPI schema during generation.
// Add fields when a component type accepts parameters beyond
// what the base Instance spec provides.
//
// +k8s:openapi-gen=true
package components

import corev1 "k8s.io/api/core/v1"

// MilvusParameters defines structured parameters for milvus components.
// This struct is converted to OpenAPI schema and served via the /schema endpoint.
// Provider users can specify these fields in the Instance's component Parameters.
type MilvusParameters struct {
	// Configuration is the Milvus engine configuration (YAML). Its contents are
	// parsed and deep-merged into the shared Milvus config file (spec.config).
	// Milvus uses a single global config, so configuration provided on any
	// component is merged into that config; use the config's per-role sections
	// (e.g. proxy, queryNode, dataCoord) to tune individual components.
	Configuration string `json:"configuration,omitempty"`
	// Pod customizes this component's pods beyond the Instance spec, e.g. for
	// secondary networks, RDMA or GPU Direct Storage.
	Pod *PodCustomization `json:"pod,omitempty"`
}

// PodCustomization holds pod-level settings applied to one component's pods.
// Volumes, the security context and init containers take Kubernetes objects
// as-is; they are free-form here to keep the published schema small and are
// validated when the Instance is reconciled.
type PodCustomization struct {
	// Annotations are added to the pods, e.g. k8s.v1.cni.cncf.io/networks to
	// attach SR-IOV secondary networks.
	Annotations map[string]string `json:"annotations,omitempty"`
	// Env adds environment variables to the Milvus container.
	Env []corev1.EnvVar `json:"env,omitempty"`
	// Volumes are added to the pods (Kubernetes Volume objects).
	Volumes []map[string]any `json:"volumes,omitempty"`
	// VolumeMounts mount volumes into the Milvus container.
	VolumeMounts []corev1.VolumeMount `json:"volumeMounts,omitempty"`
	// SecurityContext is the Milvus container's security context (a Kubernetes
	// SecurityContext object), e.g. the IPC_LOCK capability for RDMA.
	SecurityContext map[string]any `json:"securityContext,omitempty"`
	// InitContainers run before Milvus starts (Kubernetes Container objects).
	InitContainers []map[string]any `json:"initContainers,omitempty"`
}
