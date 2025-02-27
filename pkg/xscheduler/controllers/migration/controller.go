package migration

import (
	"github.com/fleezesd/xscheduler/pkg/xscheduler/framework"
	"k8s.io/apimachinery/pkg/runtime"
)

// todo: make migration controller
func New(args runtime.Object, handle framework.Handle) (framework.Plugin, error) {
	// controllerArgs, ok := args.(*xschedulerconfig.MigrationControllerArgs)
	return nil, nil
}
