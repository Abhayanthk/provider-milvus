// Package components contains parameter types for provider component types.
//
// Each struct here corresponds to a component type defined in versions.yaml
// and is converted to an OpenAPI schema during generation.
// Add fields when a component type accepts parameters beyond
// what the base Instance spec provides.
//
// +k8s:openapi-gen=true
package components

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
}
