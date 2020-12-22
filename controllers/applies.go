package controllers

import (
	"context"

	"github.com/morlay/clash-proxy/pkg/controllerutil"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func applyConfigMap(ctx context.Context, configMap *corev1.ConfigMap) error {
	c := controllerutil.ControllerClientFromContext(ctx)

	current := &corev1.ConfigMap{}

	err := c.Get(ctx, types.NamespacedName{Name: configMap.Name, Namespace: configMap.Namespace}, current)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}
		return c.Create(ctx, configMap)
	}

	if !controllerutil.IsControllerGenerationEqual(current, configMap) {
		return c.Patch(ctx, configMap, controllerutil.JSONPatch(types.StrategicMergePatchType))
	}

	return nil
}
