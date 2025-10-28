package controller

import (
	"context"
	"fmt"

	"operator-rate-limiter/pkg/ratelimiter"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/util/workqueue"
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
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Pod{}).
		Named("pod-controller").
		WithOptions(controller.Options{
			// Remove NewQueue, as workaround for type error and r.Queue not existing
			MaxConcurrentReconciles: 1,
			NewQueue:                NewRateLimitingQueue(3, 1, 60),
		}).
		Complete(r)
}

func (r *PodReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var pod corev1.Pod
	err := r.Client.Get(ctx, req.NamespacedName, &pod)
	if err != nil {
		return ctrl.Result{}, err
	}
	if pod.Name != "coredns-7c65d6cfc9-mbl9g" {
		return ctrl.Result{}, nil
	}
	logger.Info("Reconciling resource", "name", req.Name, "namespace", req.Namespace)
	if true {
		return ctrl.Result{}, fmt.Errorf("temporary error - will retry with backoff")
	}
	return ctrl.Result{}, nil
}

func NewRateLimitingQueue(maxRetries, baseDelay, maxDelay int) func(string, workqueue.TypedRateLimiter[reconcile.Request]) workqueue.TypedRateLimitingInterface[reconcile.Request] {
	return func(name string, _ workqueue.TypedRateLimiter[reconcile.Request]) workqueue.TypedRateLimitingInterface[reconcile.Request] {
		rateLimiter := ratelimiter.NewCustomRateLimiter(baseDelay, maxDelay, maxRetries)

		return ratelimiter.NewCustomQueue(rateLimiter)
	}
}
