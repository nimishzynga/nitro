package plasma

import (
	"fmt"
	"github.com/t3rm1n4l/nitro/skiplist"
	"math/rand"
	"runtime"
	"unsafe"
)

type Config struct {
	MaxDeltaChainLen int
	MaxPageItems     int
	MinPageItems     int
	Compare          skiplist.CompareFn
	ItemSize         ItemSizeFn

	MaxSize             int64
	File                string
	FlushBufferSize     int
	NumPersistorThreads int
	NumEvictorThreads   int

	LSSCleanerThreshold int
	AutoLSSCleaning     bool
	AutoSwapper         bool

	EnableShapshots bool

	// TODO: Remove later
	MaxMemoryUsage int
	shouldSwap     func() bool
}

func applyConfigDefaults(cfg Config) Config {
	if cfg.NumPersistorThreads == 0 {
		cfg.NumPersistorThreads = runtime.NumCPU()
	}

	if cfg.NumEvictorThreads == 0 {
		cfg.NumEvictorThreads = runtime.NumCPU()
	}

	// TODO: Remove later
	if cfg.shouldSwap == nil && cfg.MaxMemoryUsage > 0 {
		cfg.shouldSwap = func() bool {
			return ProcessRSS() >= int(0.7*float32(cfg.MaxMemoryUsage))
		}
	}
	return cfg
}

func DefaultConfig() Config {
	return Config{
		MaxDeltaChainLen: 200,
		MaxPageItems:     400,
		MinPageItems:     25,
		Compare:          cmpItem,
		ItemSize: func(itm unsafe.Pointer) uintptr {
			return uintptr((*item)(itm).Size())
		},
		FlushBufferSize:     1024 * 1024,
		LSSCleanerThreshold: 10,
		AutoLSSCleaning:     true,
		AutoSwapper:         false,
		EnableShapshots:     true,

		MaxMemoryUsage: 1024 * 1024 * 1024 * 512,
		MaxSize:        1024 * 1024 * 1024 * 100,
		File:           fmt.Sprintf("/tmp/plasma-data-%d.db", rand.Int()),
	}
}