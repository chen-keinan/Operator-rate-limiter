package ratelimiter

import (
	"time"

	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// MaxRetryHandlerCallback is a callback function when max retries reached
type MaxRetryHandler func(req reconcile.Request, retries int)

// CustomRateLimiter wraps a standard limiter and calls MaxRetryHandlerCallback on max retries
type CustomRateLimiter struct {
	inner      workqueue.TypedRateLimiter[reconcile.Request]
	maxRetries int
}

func NewCustomRateLimiter(
	base, maxDelaySeconds int,
	maxRetries int,
) *CustomRateLimiter {

	return &CustomRateLimiter{
		inner: workqueue.NewTypedItemExponentialFailureRateLimiter[reconcile.Request](
			time.Duration(base)*time.Second, time.Duration(maxDelaySeconds)*time.Second,
		),
		maxRetries: maxRetries,
	}
}

func (r *CustomRateLimiter) When(item reconcile.Request) time.Duration {
	return r.inner.When(item)
}

func (r *CustomRateLimiter) Forget(item reconcile.Request) {
	r.inner.Forget(item)
}

func (r *CustomRateLimiter) NumRequeues(item reconcile.Request) int {
	return r.inner.NumRequeues(item)
}
