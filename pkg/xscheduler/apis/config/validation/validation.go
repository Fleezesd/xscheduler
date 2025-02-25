package validation

import (
	"github.com/fleezesd/xscheduler/pkg/xscheduler/apis/config"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
	componentbasevalidation "k8s.io/component-base/config/validation"
)

// ValidateKubeSchedulerConfiguration ensures validation of the XschedulerConfiguration struct
func ValidateXschedulerConfiguration(cc *config.XschedulerConfiguration) utilerrors.Aggregate {
	var errs []error
	errs = append(errs, componentbasevalidation.ValidateClientConnectionConfiguration(&cc.ClientConnection, field.NewPath("clientConnection")).ToAggregate())
	errs = append(errs, componentbasevalidation.ValidateLeaderElectionConfiguration(&cc.LeaderElection, field.NewPath("leaderElection")).ToAggregate())

	if cc.LeaderElection.LeaderElect && cc.LeaderElection.ResourceLock != "leases" {
		leaderElection := field.NewPath("leaderElection")
		errs = append(errs, field.Invalid(leaderElection.Child("resourceLock"), cc.LeaderElection.ResourceLock, `resourceLock value must be "leases"`))
	}

	profilesPath := field.NewPath("profiles")
	if len(cc.Profiles) == 0 {
		errs = append(errs, field.Required(profilesPath, ""))
	} else {
		existingProfiles := make(map[string]int, len(cc.Profiles))
		for i := range cc.Profiles {
			profile := &cc.Profiles[i]
			path := profilesPath.Index(i)
			errs = append(errs, validateXschedulerProfile(path, profile)...)
			if idx, ok := existingProfiles[profile.Name]; ok {
				errs = append(errs, field.Duplicate(path.Child("name"), profilesPath.Index(idx).Child("name")))
			}
			existingProfiles[profile.Name] = i
		}
	}
	// todo: another validation & know path usage
	return utilerrors.Flatten(utilerrors.NewAggregate(errs))
}

func validateXschedulerProfile(path *field.Path, profile *config.XschedulerProfile) []error {
	var errs []error
	if len(profile.Name) == 0 {
		errs = append(errs, field.Required(path.Child("name"), ""))
	}
	errs = append(errs, validatePluginConfig(path, profile)...)
	return errs
}

func validatePluginConfig(path *field.Path, profile *config.XschedulerProfile) []error {
	// todo validate plugin config
	return []error{}
}
