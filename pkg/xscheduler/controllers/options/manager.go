package options

import (
	ctrl "sigs.k8s.io/controller-runtime"
)

// Manager is set when initializing descheduler.
// Some descheduling scenarios need to be implemented as Controller
// via controller-runtime to simplify the implementation.
var Manager ctrl.Manager
