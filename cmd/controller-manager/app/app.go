/*
Copyright 2022 The MultiClusterPlatform Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package app

import (
	"context"
	"flag"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgodiscovery "k8s.io/client-go/discovery"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"
	workv1 "open-cluster-management.io/api/work/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	appsv1alpha1 "github.com/multi-cluster-platform/mcp/pkg/apis/apps/v1alpha1"
	"github.com/multi-cluster-platform/mcp/pkg/controllers"
	"github.com/multi-cluster-platform/mcp/pkg/discovery"
	controllermanageropts "github.com/multi-cluster-platform/mcp/pkg/options/controller-manager"
	// +kubebuilder:scaffold:imports
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	// +kubebuilder:scaffold:scheme
	utilruntime.Must(appsv1alpha1.AddToScheme(scheme))
	utilruntime.Must(workv1.Install(scheme))
}

// NewControllerManagerCommand creates a *cobra.Command object with default parameters
func NewControllerManagerCommand() *cobra.Command {
	opts := controllermanageropts.NewOptions()

	cmd := &cobra.Command{
		Use:  "controller-manager",
		Long: `Multi cluster platform controller controller-manager.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Log.ValidateAndApply(); err != nil {
				return err
			}

			cliflag.PrintFlags(cmd.Flags())
			buildInfo, ok := debug.ReadBuildInfo()
			if ok {
				klog.Infof("build info: \n%s", buildInfo)
			}

			if errs := opts.Validate(); len(errs) != 0 {
				return errs.ToAggregate()
			}

			ctx := ctrl.SetupSignalHandler()
			return run(ctx, opts)
		},
	}

	fs := cmd.Flags()
	opts.AddFlags(fs)
	fs.AddGoFlagSet(flag.CommandLine)

	return cmd
}

func run(ctx context.Context, opts *controllermanageropts.Options) error {
	config := ctrl.GetConfigOrDie()

	if opts.CommonOptions.EnableCRDCheck {
		dclient, err := clientgodiscovery.NewDiscoveryClientForConfig(config)
		if err != nil {
			klog.ErrorS(err, "unable to new discovery client")
			os.Exit(1)
		}
		if !discovery.CRDsInstalled(dclient) {
			klog.Error("crd not installed")
			os.Exit(1)
		}
	}

	mgr, err := ctrl.NewManager(config, ctrl.Options{
		Scheme:                     scheme,
		MetricsBindAddress:         opts.MetricsAddr,
		HealthProbeBindAddress:     opts.ProbeAddr,
		LeaderElection:             opts.LeaderElection.LeaderElect,
		LeaderElectionResourceLock: opts.LeaderElection.ResourceLock,
		LeaderElectionNamespace:    opts.LeaderElection.ResourceNamespace,
		LeaderElectionID:           opts.LeaderElection.ResourceName,
	})
	if err != nil {
		klog.ErrorS(err, "unable to start controller-manager")
		os.Exit(1)
	}

	if err = (&controllers.ManifestWorkController{
		Client: mgr.GetClient(),
		Reader: mgr.GetAPIReader(),
	}).SetupWithManager(mgr, controller.Options{
		MaxConcurrentReconciles: opts.ConcurrencyManifestWork,
	}); err != nil {
		klog.ErrorS(err, "unable to create manifestwork controller")
		os.Exit(1)
	}

	// +kubebuilder:scaffold:builder

	klog.Info("starting controller-manager")
	if err := mgr.Start(ctx); err != nil {
		klog.ErrorS(err, "unable to run controller-manager")
		os.Exit(1)
	}

	// never reach here
	return nil
}
