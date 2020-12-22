package controllers

import (
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func SetupWithManager(mgr ctrl.Manager) error {
	return SetupReconcilerWithManager(
		mgr,
		&ClashProxyReconciler{
			Client: mgr.GetClient(),
			Log:    mgr.GetLogger().WithName("controllers").WithName("ClashProxy"),
			Scheme: mgr.GetScheme(),
		},
	)
}

type Reconciler interface {
	SetupWithManager(mgr ctrl.Manager) error
}

func SetupReconcilerWithManager(mgr manager.Manager, reconcilers ...Reconciler) error {
	for i := range reconcilers {
		if err := reconcilers[i].SetupWithManager(mgr); err != nil {
			return err
		}
	}
	return nil
}
