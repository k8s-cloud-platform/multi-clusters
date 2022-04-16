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

package controllermanager

import (
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/validation/field"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	componentbaseconfig "k8s.io/component-base/config"
	"k8s.io/component-base/logs"

	"github.com/multi-cluster-platform/mcp/pkg/options/common"
)

var (
	resourceNamespace = "mcp-system"
	resourceName      = "controller-manager.mcp.io"
)

type Options struct {
	ProbeAddr   string
	MetricsAddr string

	ConcurrencyManifestWork int

	CommonOptions *common.Options
	Log           *logs.Options

	// LeaderElection defines the configuration of leader election client.
	LeaderElection componentbaseconfig.LeaderElectionConfiguration
}

func NewOptions() *Options {
	return &Options{
		Log:           logs.NewOptions(),
		CommonOptions: common.NewCommonOptions(),
		LeaderElection: componentbaseconfig.LeaderElectionConfiguration{
			ResourceLock:      resourcelock.LeasesResourceLock,
			ResourceNamespace: resourceNamespace,
			ResourceName:      resourceName,
		},
	}
}

// AddFlags adds flags to the specified FlagSet.
func (o *Options) AddFlags(flags *pflag.FlagSet) {
	utilfeature.DefaultMutableFeatureGate.AddFlag(flags)
	o.Log.AddFlags(flags)
	o.CommonOptions.AddFlags(flags)

	flags.StringVar(&o.ProbeAddr, "health-probe-bind-address", "0",
		"The address the probe endpoint binds to.")

	flags.StringVar(&o.MetricsAddr, "metrics-bind-address", "0",
		"The address the metric endpoint binds to.")

	flags.BoolVar(&o.LeaderElection.LeaderElect, "leader-elect", true,
		"Enable leader elect.")

	flags.IntVar(&o.ConcurrencyManifestWork, "concurrency-manifestwork", 10,
		"Concurrency of ManifestWork controller.")
}

// Validate checks Options and return a slice of found errs.
func (o *Options) Validate() field.ErrorList {
	return nil
}
