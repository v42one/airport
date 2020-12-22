package controllerutil

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type contextKeyControllerClient int

func ContextWithControllerClient(ctx context.Context, client client.Client) context.Context {
	return context.WithValue(ctx, contextKeyControllerClient(1), client)
}

func ControllerClientFromContext(ctx context.Context) client.Client {
	if i, ok := ctx.Value(contextKeyControllerClient(1)).(client.Client); ok {
		return i
	}
	return nil
}
