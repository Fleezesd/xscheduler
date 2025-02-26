package xscheduler

import (
	"time"

	"github.com/fleezesd/xscheduler/pkg/xscheduler/framework"
	"github.com/fleezesd/xscheduler/pkg/xscheduler/profile"
	corev1informers "k8s.io/client-go/informers/core/v1"
	clientset "k8s.io/client-go/kubernetes"
)

type Xscheduler struct {
	// Profiles are the xscheduling profik
	Profiles profile.Map
	// Close this to shut down the scheduler.
	StopEverything <-chan struct{}

	client       clientset.Interface
	nodeInformer corev1informers.NodeInformer

	xschedulingInterval time.Duration
	nodeSelector        string
	evictionLimiter     framework.Evictor
}
