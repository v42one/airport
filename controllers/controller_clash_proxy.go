package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/morlay/clash-proxy/pkg/controllerutil"
	"gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ClashProxyReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (r *ClashProxyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		Complete(r)
}

func (r *ClashProxyReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	d := &appsv1.Deployment{}

	if err := r.Client.Get(ctx, request.NamespacedName, d); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	ls := d.Labels

	if ls == nil || ls["clash-proxy-type"] == "" || d.Annotations == nil {
		return reconcile.Result{}, nil
	}

	proxies, err := r.ListProxies(ctx, d.Namespace)
	if err != nil {
		r.Log.Error(err, "list proxy failed")
		return reconcile.Result{}, nil
	}

	ctx = controllerutil.ContextWithControllerClient(ctx, r.Client)

	c := &corev1.ConfigMap{}
	c.Name = CLASH_PROXY_CONFIG_MAP_NAME
	c.Namespace = d.Namespace

	data, _ := yaml.Marshal(struct {
		Proxies []Proxy `yaml:"proxies"`
	}{proxies})

	c.Data = map[string]string{
		"proxy.yaml": string(data),
	}

	if err := applyConfigMap(ctx, c); err != nil {
		r.Log.Error(err, "apply failed")
		return reconcile.Result{}, nil
	}

	return reconcile.Result{}, nil
}

func (r *ClashProxyReconciler) ListProxies(ctx context.Context, namespace string) ([]Proxy, error) {
	deployments := &appsv1.DeploymentList{}

	s, _ := labels.NewRequirement("clash-proxy-type", selection.Exists, []string{})

	if err := r.Client.List(ctx,
		deployments,
		client.InNamespace(namespace),
		client.MatchingLabelsSelector{Selector: labels.NewSelector().Add(*s)},
	); err != nil {
		return nil, err
	}

	proxies := make([]Proxy, 0)

	for i := range deployments.Items {
		proxies = append(proxies, ProxiesFromDeployment(&deployments.Items[i], publicIPs)...)
	}

	return proxies, nil
}
