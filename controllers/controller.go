package controllers

import (
	"context"

	"github.com/morlay/clash-proxy/pkg/clashproxy"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrlruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Reconciler interface {
	SetupWithManager(mgr ctrlruntime.Manager) error
}

func SetupReconcilerWithManager(mgr manager.Manager, reconcilers ...Reconciler) error {
	for i := range reconcilers {
		if err := reconcilers[i].SetupWithManager(mgr); err != nil {
			return err
		}
	}
	return nil
}

type ClashProxyReconciler struct {
	client.Client
	Provider *clashproxy.Provider
	Log      logr.Logger
	Scheme   *runtime.Scheme
}

func (r *ClashProxyReconciler) SetupWithManager(mgr ctrlruntime.Manager) error {
	return ctrlruntime.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		Complete(r)
}

func (r *ClashProxyReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	d := &appsv1.Deployment{}

	if err := r.Client.Get(ctx, request.NamespacedName, d); err != nil {
		if apierrors.IsNotFound(err) {
			r.Provider.Remove(request.NamespacedName.String())
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	if d.Labels != nil {
		if _, ok := d.Labels[LabelClashProxyType]; !ok {
			return reconcile.Result{}, nil
		}
	}

	r.Provider.Add(request.NamespacedName.String(), ProxyOptionFromDeployment(d))

	return reconcile.Result{}, nil
}
