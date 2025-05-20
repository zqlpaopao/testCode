package pkg

import "sync"

type GraphHandler[T any, N any, E any] struct {
	Graph     T
	NodeCache map[string]N
	EdgeCache map[string]E
	NodeLock  sync.RWMutex
	EdgeLock  sync.RWMutex
}

// NewGraphHandler -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func NewGraphHandler[T any, N any, E any](
	graph T,
	nodeCacheSize int,
	edgeCacheSize int,
) *GraphHandler[T, N, E] {
	return &GraphHandler[T, N, E]{
		Graph:     graph,
		NodeCache: make(map[string]N, nodeCacheSize),
		EdgeCache: make(map[string]E, edgeCacheSize),
		NodeLock:  sync.RWMutex{},
		EdgeLock:  sync.RWMutex{},
	}
}
