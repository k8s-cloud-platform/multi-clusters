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
	"flag"
	"runtime/debug"

	"github.com/spf13/cobra"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"

	apiserveropts "github.com/multi-cluster-platform/mcp/pkg/options/apiserver"
)

// NewApiServerCommand creates a *cobra.Command object with default parameters
func NewApiServerCommand() *cobra.Command {
	opts := apiserveropts.NewOptions()

	cmd := &cobra.Command{
		Use:  "mcp-apiserver",
		Long: `Multi cluster platform aggregated apiserver.`,
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

			config, err := opts.Config()
			if err != nil {
				return err
			}

			server, err := config.Complete().New()
			if err != nil {
				return err
			}

			ctx := ctrl.SetupSignalHandler()
			return server.Run(ctx, opts)
		},
	}

	fs := cmd.Flags()
	opts.AddFlags(fs)
	fs.AddGoFlagSet(flag.CommandLine)

	return cmd
}
