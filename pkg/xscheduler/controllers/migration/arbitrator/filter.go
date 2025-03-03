package arbitrator

import (
	"github.com/fleezesd/xscheduler/pkg/xscheduler/apis/config"
	"github.com/fleezesd/xscheduler/pkg/xscheduler/framework"
	"k8s.io/utils/clock"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type filter struct {
	client client.Client
	clock  clock.RealClock

	nonRetryablePodFilter framework.FilterFunc
	retryablePodFilter    framework.FilterFunc
	defaultFilterPlugin   framework.FilterPlugin

	args *config.MigrationControllerArgs
}
