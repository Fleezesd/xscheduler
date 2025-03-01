package migration

import (
	"context"
	"sync"

	gocache "github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"golang.org/x/time/rate"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	xschedulerconfig "github.com/fleezesd/xscheduler/pkg/xscheduler/apis/config"
	"github.com/fleezesd/xscheduler/pkg/xscheduler/apis/config/validation"
	"github.com/fleezesd/xscheduler/pkg/xscheduler/controllers/migration/reservation"
	"github.com/fleezesd/xscheduler/pkg/xscheduler/controllers/names"
	"github.com/fleezesd/xscheduler/pkg/xscheduler/controllers/options"
	"github.com/fleezesd/xscheduler/pkg/xscheduler/framework"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/events"
	"k8s.io/utils/clock"
)

const (
	Name = names.MigrationController
)

type Reconciler struct {
	client.Client
	args          *xschedulerconfig.MigrationControllerArgs
	eventRecorder events.EventRecorder

	// todo: interpreter & finder
	reservationInterpreter reservation.Interpreter
	// todo: assumedcache

	// todo: arbitrator

	clock clock.Clock

	// rate limiter
	limiterMap      map[xschedulerconfig.MigrationLimitObjectType]map[string]*rate.Limiter
	limiterCacheMap map[xschedulerconfig.MigrationLimitObjectType]*gocache.Cache
	limiterLock     sync.Mutex
}

func New(args runtime.Object, handle framework.Handle) (framework.Plugin, error) {
	controllerArgs, ok := args.(*xschedulerconfig.MigrationControllerArgs)
	if !ok {
		return nil, errors.Errorf("want args to be of type MigrationControllerArgs, got %T", args)
	}
	if err := validation.ValidateMigrationControllerArgs(nil, controllerArgs); err != nil {
		return nil, err
	}

	_, err := newReconciler(controllerArgs, handle)
	if err != nil {
		return nil, err
	}

	// todo make reconciler
	_, err = controller.New(Name, options.Manager, controller.Options{})
	return nil, nil
}

func newReconciler(args *xschedulerconfig.MigrationControllerArgs, handle framework.Handle) (*Reconciler, error) {
	manager := options.Manager

	// todo: interpreter & finder
	reservationInterpreter := reservation.NewInterpreter(manager)

	r := &Reconciler{
		Client:                 manager.GetClient(),
		args:                   args,
		eventRecorder:          handle.EventRecorder(),
		clock:                  clock.RealClock{},
		reservationInterpreter: reservationInterpreter,
	}
	r.initObjectLimiters()
	if err := manager.Add(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Reconciler) Name() string {
	return Name
}

func (r *Reconciler) Start(ctx context.Context) error {
	return nil
}

func (r *Reconciler) initObjectLimiters() {
	// store actual rate limiter
	r.limiterMap = make(map[xschedulerconfig.MigrationLimitObjectType]map[string]*rate.Limiter)
	// store limter expiration
	r.limiterCacheMap = make(map[xschedulerconfig.MigrationLimitObjectType]*gocache.Cache)

	for limiterType, limiterConfig := range r.args.ObjectLimiters {
		trackExpiration := limiterConfig.Duration.Duration
		if trackExpiration > 0 {
			r.limiterMap[limiterType] = make(map[string]*rate.Limiter)
			// Set cache expiration to 1.5x the tracking duration
			cacheExpiration := trackExpiration + (trackExpiration / 2)

			cache := gocache.New(cacheExpiration, cacheExpiration)
			cache.OnEvicted(func(key string, _ interface{}) {
				r.limiterLock.Lock()
				// delete expiration rate limiter
				delete(r.limiterMap[limiterType], key)
				r.limiterLock.Unlock()
			})
			r.limiterCacheMap[limiterType] = cache
		}
	}
}
