package ratelimiter

import (
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type CustomQueue struct {
	workqueue.TypedRateLimitingInterface[reconcile.Request]
	rateLimiter workqueue.TypedRateLimiter[reconcile.Request]
	maxRetries  int
}

func NewCustomQueue(rateLimiter workqueue.TypedRateLimiter[reconcile.Request], maxRetries int) *CustomQueue {
	return &CustomQueue{
		TypedRateLimitingInterface: workqueue.NewTypedRateLimitingQueue(rateLimiter),
		rateLimiter:                rateLimiter,
		maxRetries:                 maxRetries,
	}
}

func (q *CustomQueue) AddRateLimited(item reconcile.Request) {
	if q.NumRequeues(item) >= q.maxRetries {
		q.done(item)
		return
	}
	q.TypedRateLimitingInterface.AddRateLimited(item)
}

func (q *CustomQueue) NumRequeues(item reconcile.Request) int {
	return q.TypedRateLimitingInterface.NumRequeues(item)
}

func (q *CustomQueue) Forget(item reconcile.Request) {
	q.TypedRateLimitingInterface.Forget(item)
}

func (q *CustomQueue) done(item reconcile.Request) {
	q.TypedRateLimitingInterface.Forget(item)
	q.TypedRateLimitingInterface.Done(item)
}
