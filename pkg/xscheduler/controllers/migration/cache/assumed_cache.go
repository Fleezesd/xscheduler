package cache

import (
	"sync"

	"k8s.io/apimachinery/pkg/types"
)

// todo: make assumed cache for xscheduler
type assumedCache struct {
	lock  sync.Mutex
	items map[types.UID]any
}
