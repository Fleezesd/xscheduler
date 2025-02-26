package validation

import (
	xschedulerconfiguration "github.com/fleezesd/xscheduler/pkg/xscheduler/apis/config"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateMigrationControllerArgs(path *field.Path, args *xschedulerconfiguration.MigrationControllerArgs) error {
	return nil
}
