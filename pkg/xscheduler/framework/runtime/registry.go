package runtime

import (
	"fmt"

	"github.com/fleezesd/xscheduler/pkg/xscheduler/framework"
	"k8s.io/apimachinery/pkg/runtime"
)

type PluginFactory func(args runtime.Object, handle framework.Handle) (framework.Plugin, error)

// Regsitry is a plugin map for out of tree registry
type Registry map[string]PluginFactory

// Register adds a new plugin to the registry. If a plugin with the same name
// exists, it returns an error.
func (r Registry) Register(name string, factory PluginFactory) error {
	if _, ok := r[name]; ok {
		return fmt.Errorf("a plugin named %v already exists", name)
	}
	r[name] = factory
	return nil
}

// Unregister removes an existing plugin from the registry. If no plugin with
// the provided name exists, it returns an error.
func (r Registry) UnRegister(name string) error {
	if _, ok := r[name]; !ok {
		return fmt.Errorf("no plugin named %v exists", name)
	}
	delete(r, name)
	return nil
}

func (r Registry) Merge(in Registry) error {
	for name, factory := range in {
		if err := r.Register(name, factory); err != nil {
			return err
		}
	}
	return nil
}
