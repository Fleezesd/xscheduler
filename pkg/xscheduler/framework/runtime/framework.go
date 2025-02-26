package runtime

import "github.com/fleezesd/xscheduler/pkg/xscheduler/framework"

type frameworkImpl struct {
	evictPlugins  []framework.EvictPlugin
	filterPlugins []framework.FilterPlugin
}
