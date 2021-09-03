package controllers

import (
	"context"
	"os"
	"strings"

	"github.com/morlay/clash-proxy/pkg/clashproxy"
	"golang.org/x/sync/errgroup"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	pkgerrors "github.com/pkg/errors"
	ctrlruntimeruntime "sigs.k8s.io/controller-runtime"
)

type Manager struct {
	ctrlruntimeruntime.Options
	provider *clashproxy.Provider
	mgr      manager.Manager
}

var CLUSTER_PUBLIC_IPS = strings.TrimSpace(os.Getenv("CLUSTER_PUBLIC_IPS"))

func (m *Manager) Init(ctx context.Context) error {
	restConfig := ctrlruntimeruntime.GetConfigOrDie()

	mgr, err := ctrlruntimeruntime.NewManager(restConfig, m.Options)
	if err != nil {
		return pkgerrors.Wrap(err, "unable to start manager")
	}

	m.mgr = mgr

	if CLUSTER_PUBLIC_IPS != "" {
		m.provider = clashproxy.ParsePublicIPs(CLUSTER_PUBLIC_IPS)
	} else {
		m.provider = &clashproxy.Provider{}

		if err := m.provider.LoadPublicIPsFromKube(ctx, mgr.GetAPIReader()); err != nil {
			return pkgerrors.Wrap(err, "load public ips failed")
		}
	}

	if err := SetupReconcilerWithManager(
		mgr,
		&ClashProxyReconciler{
			Provider: m.provider,
			Client:   mgr.GetClient(),
			Log:      mgr.GetLogger().WithName("controllers").WithName("ClashProxy"),
			Scheme:   mgr.GetScheme(),
		},
	); err != nil {
		return pkgerrors.Wrap(err, "unable to create controller")
	}

	return nil
}

func (m *Manager) Start(ctx context.Context) error {

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return m.mgr.Start(ctrlruntimeruntime.SetupSignalHandler())
	})

	g.Go(func() error {
		return m.provider.Start(ctx)
	})

	return g.Wait()
}
