package reservation

import (
	"context"
	"fmt"

	schedulingv1alpha1 "github.com/fleezesd/xscheduler/apis/scheduling/v1alpha1"
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ Interpreter = &interpreterImpl{}

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

func (o *interpreterImpl) GetReservationType() client.Object {
	return &schedulingv1alpha1.Reservation{}
}

func (o *interpreterImpl) Preemption() Preemption {
	return nil
}

func (o *interpreterImpl) CreateReservation(
	ctx context.Context,
	job *schedulingv1alpha1.PodMigrationJob) (Object, error) {
	reservationOptions := job.Spec.ReservationOptions

	reservation := &schedulingv1alpha1.Reservation{
		ObjectMeta: reservationOptions.Template.ObjectMeta,
		Spec:       reservationOptions.Template.Spec,
	}

	err := o.Client.Create(ctx, reservation)
	if err != nil {
		return nil, err
	}
	return &Reservation{Reservation: reservation}, nil
}

func (o *interpreterImpl) GetReservation(ctx context.Context, reservationRef *corev1.ObjectReference) (Object, error) {
	if reservationRef == nil {
		return nil, fmt.Errorf("reservation reference cannot be nil")
	}

	reservation := &schedulingv1alpha1.Reservation{}
	namespacedName := GetReservationNamespacedName(reservationRef)

	// Try getting reservation from cache first
	if err := o.Client.Get(ctx, namespacedName, reservation); err == nil {
		return &Reservation{Reservation: reservation}, nil
	} else if !errors.IsNotFound(err) {
		return nil, fmt.Errorf("failed to get reservation from cache: %w", err)
	}

	// If not found in cache, try getting directly from API server
	klog.Warningf("Reservation %v not found in cache, attempting to fetch from API server", reservationRef.Name)
	if err := o.mgr.GetAPIReader().Get(ctx, namespacedName, reservation); err != nil {
		return nil, fmt.Errorf("failed to get reservation from API server: %w", err)
	}

	return &Reservation{Reservation: reservation}, nil
}
func (o *interpreterImpl) DeleteReservation(ctx context.Context, reservationRef *corev1.ObjectReference) error {
	if lo.IsNil(reservationRef) {
		return nil
	}

	namespacedName := GetReservationNamespacedName(reservationRef)
	logger := klog.V(4)
	logger.Infof("Deleting Reservation %v", namespacedName)

	reservation, err := o.GetReservation(ctx, reservationRef)
	if err != nil {
		return err
	}

	if err := o.Client.Delete(ctx, reservation.OriginObject()); err != nil {
		klog.Errorf("Failed to delete Reservation %v, err: %v", namespacedName, err)
		return err
	}

	logger.Infof("Successfully deleted Reservation %v", namespacedName)
	return nil
}
