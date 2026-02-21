package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/viper"
	_ "k8s.io/client-go/plugin/pkg/client/auth" // Import all auth plugins for exec-entrypoint

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/discovery"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwapiv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	cfgatev1alpha1 "cfgate.io/cfgate/api/v1alpha1"
	cfcloudflare "cfgate.io/cfgate/internal/cloudflare"
	"cfgate.io/cfgate/internal/controller"
	"cfgate.io/cfgate/internal/controller/features"
)

var (
	Version   = "0.0.0-dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(cfgatev1alpha1.AddToScheme(scheme))
	utilruntime.Must(gwapiv1.Install(scheme))
	utilruntime.Must(gwapiv1beta1.Install(scheme))

	// Viper configuration
	viper.SetEnvPrefix("CFGATE")
	viper.AutomaticEnv()
	viper.SetDefault("metrics.port", 8080)
	viper.SetDefault("health.port", 8081)
}

func main() {
	var metricsAddr string
	var probeAddr string
	var enableLeaderElection bool
	var secureMetrics bool

	flag.StringVar(&metricsAddr, "metrics-bind-address",
		fmt.Sprintf(":%d", viper.GetInt("metrics.port")),
		"The address the metrics endpoint binds to. Use :8443 for HTTPS or :8080 for HTTP.")
	flag.StringVar(&probeAddr, "health-probe-bind-address",
		fmt.Sprintf(":%d", viper.GetInt("health.port")),
		"The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&secureMetrics, "metrics-secure", false,
		"If set, the metrics endpoint is served securely via HTTPS.")
	opts := zap.Options{
		Development: false,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	// Log startup configuration. Key-value pairs follow logr conventions.
	setupLog.Info("starting cfgate controller manager",
		"version", Version,
		"commit", Commit,
		"buildDate", BuildDate,
		"metricsAddr", metricsAddr,
		"healthProbeAddr", probeAddr,
		"leaderElection", enableLeaderElection,
		"secureMetrics", secureMetrics,
	)

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
		Metrics: metricsserver.Options{
			BindAddress:   metricsAddr,
			SecureServing: secureMetrics,
		},
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "cfgate.io",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// Detect optional Gateway API CRDs for conditional controller behavior.
	dc, err := discovery.NewDiscoveryClientForConfig(mgr.GetConfig())
	if err != nil {
		setupLog.Error(err, "unable to create discovery client")
		os.Exit(1)
	}
	featureGates, err := features.DetectFeatures(dc)
	if err != nil {
		setupLog.Error(err, "unable to detect feature gates")
		os.Exit(1)
	}
	featureGates.LogFeatures(setupLog)

	// Shared credential cache for all Cloudflare-facing reconcilers.
	credCache := cfcloudflare.NewCredentialCache(0) // 0 = default TTL

	if err = (&controller.CloudflareTunnelReconciler{
		Client:          mgr.GetClient(),
		Scheme:          mgr.GetScheme(),
		Recorder:        mgr.GetEventRecorder("cloudflaretunnel-controller"),
		CredentialCache: credCache,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "CloudflareTunnel")
		os.Exit(1)
	}

	if err = (&controller.CloudflareDNSReconciler{
		Client:          mgr.GetClient(),
		Scheme:          mgr.GetScheme(),
		Recorder:        mgr.GetEventRecorder("cloudflaredns-controller"),
		CredentialCache: credCache,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "CloudflareDNS")
		os.Exit(1)
	}

	if err = (&controller.GatewayReconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorder("gateway-controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Gateway")
		os.Exit(1)
	}

	if err = (&controller.GatewayClassReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "GatewayClass")
		os.Exit(1)
	}

	if err = (&controller.HTTPRouteReconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorder("httproute-controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "HTTPRoute")
		os.Exit(1)
	}

	if err = (&controller.CloudflareAccessPolicyReconciler{
		Client:          mgr.GetClient(),
		Scheme:          mgr.GetScheme(),
		Recorder:        mgr.GetEventRecorder("cloudflareaccesspolicy-controller"),
		FeatureGates:    featureGates,
		CredentialCache: credCache,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "CloudflareAccessPolicy")
		os.Exit(1)
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("all controllers registered, starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "manager stopped with error")
		os.Exit(1)
	}
	setupLog.Info("manager shutdown complete")
}
