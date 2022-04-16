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
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"

	"github.com/multi-cluster-platform/mcp/pkg/apis/gateway"
	gatewayregistry "github.com/multi-cluster-platform/mcp/pkg/registry/gateway/cluster"
	"github.com/multi-cluster-platform/mcp/pkg/registry/gateway/shadow"
)

// ExtraConfig holds custom apiserver config
type ExtraConfig struct {
	// TODO Place you custom config here.
}

// Config defines the config for the apiserver
type Config struct {
	GenericConfig *genericapiserver.RecommendedConfig
	ExtraConfig   *ExtraConfig
}

// CompletedConfig embeds a private pointer that cannot be instantiated outside of this package.
type CompletedConfig struct {
	*completedConfig
}

type completedConfig struct {
	GenericConfig genericapiserver.CompletedConfig
	ExtraConfig   *ExtraConfig
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (cfg *Config) Complete() CompletedConfig {
	c := completedConfig{
		cfg.GenericConfig.Complete(),
		cfg.ExtraConfig,
	}

	c.GenericConfig.Version = &version.Info{
		Major: "1",
		Minor: "0",
	}

	return CompletedConfig{&c}
}

// New returns a new instance of MCPServer from the given config.
func (c completedConfig) New() (*MCPServer, error) {
	genericServer, err := c.GenericConfig.New("mcp-apiserver", genericapiserver.NewEmptyDelegate())
	if err != nil {
		return nil, err
	}

	s := &MCPServer{
		GenericAPIServer: genericServer,
	}

	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(gateway.GroupName, Scheme, ParameterCodec, Codecs)

	v1alpha1storage := map[string]rest.Storage{}
	v1alpha1storage["shadow"] = shadow.NewREST()
	v1alpha1storage["clusters"] = gatewayregistry.NewREST()
	apiGroupInfo.VersionedResourcesStorageMap["v1alpha1"] = v1alpha1storage

	if err := s.GenericAPIServer.InstallAPIGroup(&apiGroupInfo); err != nil {
		return nil, err
	}

	return s, nil
}
