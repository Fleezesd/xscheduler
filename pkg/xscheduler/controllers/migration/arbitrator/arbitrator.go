package arbitrator

import (
	"sync"
	"time"

	"github.com/fleezesd/xscheduler/apis/scheduling/v1alpha1"
	"github.com/fleezesd/xscheduler/pkg/xscheduler/apis/config"
	"github.com/fleezesd/xscheduler/pkg/xscheduler/framework"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/events"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MigrationFilter interface {
	Filter(pod *corev1.Pod) bool
	PreEvictionFilter(pod *corev1.Pod) bool
}

type Arbitrator interface {
	MigrationFilter
	AddPodMigrationJob(job *v1alpha1.PodMigrationJob)
	DeletePodMigrationJob(job *v1alpha1.PodMigrationJob)
}

// SortFn stably sorts PodMigrationJobs slice based on a certain strategy. Users
// can implement different SortFn according to their needs.
type SortFn func(jobs []*v1alpha1.PodMigrationJob, podOfJob map[*v1alpha1.PodMigrationJob]*corev1.Pod) []*v1alpha1.PodMigrationJob

type arbitratorImpl struct {
	waitingCollection map[types.UID]*v1alpha1.PodMigrationJob
	interval          time.Duration
	sorts             []SortFn
	filter            *filter
	client            client.Client
	eventRecorder     events.EventRecorder
	mu                sync.Mutex
}

type Options struct {
	Client        client.Client
	EventRecorder events.EventRecorder
	Manager       ctrl.Manager
	Handle        framework.Handle
}

func New(args *config.MigrationControllerArgs, options Options) (Arbitrator, error) {
	return nil, nil
}
