// Package common defines shared constants used across the provider.
package common

const (
	// ProviderName is this provider's identity.
	ProviderName = "milvus"

	ComponentStandalone = "standalone"
	ComponentProxy      = "proxy"
	ComponentMixCoord   = "mixCoord"
	ComponentRootCoord  = "rootCoord"
	ComponentIndexCoord = "indexCoord"
	ComponentDataCoord  = "dataCoord"
	ComponentQueryCoord = "queryCoord"
	ComponentIndexNode  = "indexNode"
	ComponentDataNode   = "dataNode"
	ComponentQueryNode  = "queryNode"
	ComponentStreaming  = "streamingNode"

	ComponentTypeMilvus = "milvus"
)
