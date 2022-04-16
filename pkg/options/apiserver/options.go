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

package apiserver

import (
	"fmt"
	"net"

	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/admission"
	genericapiserver "k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/component-base/logs"
	netutils "k8s.io/utils/net"

	"github.com/multi-cluster-platform/mcp/pkg/apiserver"
	"github.com/multi-cluster-platform/mcp/pkg/options/common"
)

type Options struct {
	EnablesLocalDebug bool

	CommonOptions *common.Options
	Log           *logs.Options

	RecommendedOptions *genericoptions.RecommendedOptions
}

func NewOptions() *Options {
	opts := &Options{
		Log:                logs.NewOptions(),
		CommonOptions:      common.NewCommonOptions(),
		RecommendedOptions: genericoptions.NewRecommendedOptions("fake", nil),
	}
	return opts
}

// AddFlags adds flags to the specified FlagSet.
func (o *Options) AddFlags(flags *pflag.FlagSet) {
	utilfeature.DefaultMutableFeatureGate.AddFlag(flags)
	o.RecommendedOptions.AddFlags(flags)
	o.Log.AddFlags(flags)
	o.CommonOptions.AddFlags(flags)

	flags.BoolVar(&o.EnablesLocalDebug, "enable-local-debug", false,
		"Under the local-debug mode the apiserver will allow all access to its resources without authorizing the requests, this flag is only intended for debugging in your workstation.")
}

// Validate checks Options and return a slice of found errs.
func (o *Options) Validate() field.ErrorList {
	return nil
}

// Config fills in fields required to have valid data
func (o *Options) Config() (*apiserver.Config, error) {
	// TODO have a "real" external address
	if err := o.RecommendedOptions.SecureServing.MaybeDefaultWithSelfSignedCerts("localhost", nil, []net.IP{netutils.ParseIPSloppy("127.0.0.1")}); err != nil {
		return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
	}

	if o.EnablesLocalDebug {
		o.RecommendedOptions.Authorization = nil
		o.RecommendedOptions.CoreAPI = nil
		o.RecommendedOptions.Admission = nil
	}

	o.RecommendedOptions.Etcd = nil

	o.RecommendedOptions.ExtraAdmissionInitializers = func(c *genericapiserver.RecommendedConfig) ([]admission.PluginInitializer, error) {
		return []admission.PluginInitializer{}, nil
	}

	serverConfig := genericapiserver.NewRecommendedConfig(apiserver.Codecs)

	if err := o.RecommendedOptions.ApplyTo(serverConfig); err != nil {
		return nil, err
	}

	config := &apiserver.Config{
		GenericConfig: serverConfig,
		ExtraConfig:   &apiserver.ExtraConfig{},
	}
	return config, nil
}
