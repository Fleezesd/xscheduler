package validation

import (
	"net"
	"reflect"
	"strconv"

	"github.com/fleezesd/xscheduler/pkg/xscheduler/apis/config"
	"github.com/fleezesd/xscheduler/pkg/xscheduler/controllers/names"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation"
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
	// validate healthz & metrics
	for _, msg := range isValidSocketAddr(cc.HealthzBindAddress) {
		errs = append(errs, field.Invalid(field.NewPath("healthzBindAddress"), cc.HealthzBindAddress, msg))
	}
	for _, msg := range isValidSocketAddr(cc.MetricsBindAddress) {
		errs = append(errs, field.Invalid(field.NewPath("metricsBindAddress"), cc.MetricsBindAddress, msg))
	}

	// validate node selector
	if cc.NodeSelector != nil {
		_, err := metav1.LabelSelectorAsSelector(cc.NodeSelector)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return utilerrors.Flatten(utilerrors.NewAggregate(errs))
}

// IsValidSocketAddr checks that string represents a valid socket address
// as defined in RFC 789. (e.g 0.0.0.0:10254 or [::]:10254))
func isValidSocketAddr(value string) []string {
	var errs []string
	ip, port, err := net.SplitHostPort(value)
	if err != nil {
		errs = append(errs, "must be a valid socket address format, (e.g. 0.0.0.0:10254 or [::]:10254)")
		return errs
	}
	portInt, _ := strconv.Atoi(port)
	errs = append(errs, validation.IsValidPortNum(portInt)...)
	errs = append(errs, validation.IsValidIP(field.NewPath("ip"), ip).
		ToAggregate().Error())
	return errs
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
	var errs []error
	m := map[string]interface{}{
		names.MigrationController: ValidateMigrationControllerArgs,
	}
	seenPluginConfig := sets.New[string]()

	for i := range profile.PluginConfig {
		pluginConfigPath := path.Child("pluginConfig").Index(i)
		name := profile.PluginConfig[i].Name
		args := profile.PluginConfig[i].Args

		// check if or not duplicate in seenPluginConfig
		if seenPluginConfig.Has(name) {
			errs = append(errs, field.Duplicate(pluginConfigPath, name))
		} else {
			seenPluginConfig.Insert(name)
		}

		// if has validateFunc for plugin
		if validateFunc, ok := m[name]; ok {
			// check if or not args are true for validatefunc
			// 1 is args and you can see ValidateMigrationControllerArgs this func
			if reflect.TypeOf(args) != reflect.ValueOf(validateFunc).Type().In(1) {
				errs = append(errs, field.Invalid(pluginConfigPath.Child("args"), args, "has to match plugin args"))
			} else {
				// use validatefunc
				in := []reflect.Value{reflect.ValueOf(pluginConfigPath.Child("args")), reflect.ValueOf(args)}
				res := reflect.ValueOf(validateFunc).Call(in)
				if res[0].Interface() != nil {
					errs = append(errs, res[0].Interface().(error))
				}
			}
		}
	}
	return errs
}
