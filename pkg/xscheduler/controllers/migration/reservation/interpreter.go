package reservation

import (
	"context"

	schedulingv1alpha1 "github.com/fleezesd/xscheduler/apis/scheduling/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type interpreterImpl struct {
	mgr ctrl.Manager
	client.Client
}

func newInterpreter(mgr ctrl.Manager) Interpreter {
	return &interpreterImpl{
		mgr:    mgr,
		Client: mgr.GetClient(),
	}
}

func (o *interpreterImpl) GetReservationType() client.Object

func (o *interpreterImpl) Preemption()

func (o *interpreterImpl) CreateReservation(ctx context.Context, job *schedulingv1alpha1.PodMigrationJob) (Object, error)

func (o *interpreterImpl) GetReservation(ctx context.Context, reservationRef *corev1.ObjectReference) (Object, error)

func (o *interpreterImpl) DeleteReservation(ctx context.Context, reservationRef *corev1.ObjectReference) error
