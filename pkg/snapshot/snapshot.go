package snapshot

import (
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	cachev3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
)

func GenerateSnapshot() cachev3.Snapshot {
	return cachev3.NewSnapshot(
		"1",
		[]types.Resource{}, // endpoints
		[]types.Resource{}, // clusters
		[]types.Resource{}, // routes
		[]types.Resource{}, // listeners
		[]types.Resource{}, // runtimes
		[]types.Resource{},
	)
}
