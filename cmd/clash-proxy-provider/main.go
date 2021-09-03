package main

import (
	"context"
	"flag"
	"os"

	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/morlay/clash-proxy/controllers"
	"github.com/morlay/clash-proxy/pkg/version"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrlruntime "sigs.k8s.io/controller-runtime"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
}

func parseFlags(ctrlOpt *controllers.Manager) {
	flag.StringVar(&ctrlOpt.Namespace, "watch-namespace", os.Getenv("WATCH_NAMESPACE"), "watch namespace")
	flag.StringVar(&ctrlOpt.MetricsBindAddress, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&ctrlOpt.LeaderElection, "enable-leader-election", false, "Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flag.Parse()
}

func main() {
	mgr := controllers.Manager{}
	mgr.Scheme = scheme
	mgr.Port = 9443
	mgr.LeaderElectionID = "74b83f88.clash-proxy-operator"
	mgr.Logger = ctrlruntime.Log.WithValues("clash-proxy-provider", version.Version)

	parseFlags(&mgr)

	ctrlruntime.SetLogger(zap.New(zap.UseDevMode(true)))

	ctx := context.Background()

	if err := mgr.Init(ctx); err != nil {
		ctrlruntime.Log.WithName("init").Error(err, "")
	}

	if err := mgr.Start(ctx); err != nil {
		ctrlruntime.Log.WithName("start").Error(err, "")
	}
}
