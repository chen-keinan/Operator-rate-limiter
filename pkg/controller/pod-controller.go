package controller

import (
	"context"
	"fmt"

	"operator-rate-limiter/pkg/ratelimiter"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type PodReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *PodReconciler) SetupWithManager(mgr ctrl.Manager) error {

	limiter := ratelimiter.NewCustomRateLimiter(
		1,  // base 1s
		60, // max 60s delay
		5,  // max retries before callback
		func(req reconcile.Request, retries int) {
			go handleMaxRetries(mgr.GetClient(), req, retries)
		},
	)
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Pod{}).
		WithOptions(controller.Options{

			RateLimiter:             limiter,
			MaxConcurrentReconciles: 1,
		}).
		Complete(r)
}

func (r *PodReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Reconciling resource", "name", req.Name, "namespace", req.Namespace)
	if true {
		return ctrl.Result{}, fmt.Errorf("temporary error - will retry with backoff")
	}
	return ctrl.Result{}, nil
}

func handleMaxRetries(c client.Client, req reconcile.Request, retries int) {
	ctx := context.Background()
	var cr corev1.Pod
	if err := c.Get(ctx, req.NamespacedName, &cr); err != nil {
		return
	}
	cr.Status.Message = fmt.Sprintf("Max retries (%d) reached", retries)
	_ = c.Status().Update(ctx, &cr)
}
