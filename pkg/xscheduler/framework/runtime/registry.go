package runtime

import (
	"github.com/fleezesd/xscheduler/pkg/xscheduler/framework"
	"k8s.io/apimachinery/pkg/runtime"
)

type PluginFactory func(args runtime.Object, handle framework.Handle) (framework.Plugin, error)

type Registry map[string]PluginFactory
