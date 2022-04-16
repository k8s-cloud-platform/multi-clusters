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
	"context"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	genericapiserver "k8s.io/apiserver/pkg/server"
	clientgodiscovery "k8s.io/client-go/discovery"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"

	gatewayinstall "github.com/multi-cluster-platform/mcp/pkg/apis/gateway/install"
	"github.com/multi-cluster-platform/mcp/pkg/discovery"
	apiserveropts "github.com/multi-cluster-platform/mcp/pkg/options/apiserver"
)

var (
	// Scheme defines methods for serializing and deserializing API objects.
	Scheme = runtime.NewScheme()
	// Codecs provides methods for retrieving codecs and serializers for specific
	// versions and content types.
	Codecs = serializer.NewCodecFactory(Scheme)
	// ParameterCodec handles versioning of objects that are converted to query parameters.
	ParameterCodec = runtime.NewParameterCodec(Scheme)
)

func init() {
	gatewayinstall.Install(Scheme)

	// we need to add the options to empty v1
	// TODO fix the server code to avoid this
	metav1.AddToGroupVersion(Scheme, schema.GroupVersion{Version: "v1"})

	// TODO: keep the generic API server from wanting this
	unversioned := schema.GroupVersion{Group: "", Version: "v1"}
	Scheme.AddUnversionedTypes(unversioned,
		&metav1.Status{},
		&metav1.APIVersions{},
		&metav1.APIGroupList{},
		&metav1.APIGroup{},
		&metav1.APIResourceList{},
	)
}

// MCPServer contains state for a Kubernetes cluster master/api server.
type MCPServer struct {
	GenericAPIServer *genericapiserver.GenericAPIServer
}

func (server *MCPServer) Run(ctx context.Context, opts *apiserveropts.Options) error {
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

	//mgr, err := ctrl.NewManager(config, ctrl.Options{
	//	Scheme:                  scheme,
	//	MetricsBindAddress:      opts.MetricsAddr,
	//	HealthProbeBindAddress:  opts.ProbeAddr,
	//	LeaderElection:          opts.LeaderElection.LeaderElect,
	//	LeaderElectionNamespace: opts.LeaderElection.ResourceNamespace,
	//	LeaderElectionID:        opts.LeaderElection.ResourceName,
	//})
	//if err != nil {
	//	klog.ErrorS(err, "unable to start controller-manager")
	//	os.Exit(1)
	//}

	// TODO add post start hook
	//server.GenericAPIServer.AddPostStartHookOrDie("test", func(context genericapiserver.PostStartHookContext) error {
	//	klog.Info("hook for post start...")
	//	return nil
	//})

	return server.GenericAPIServer.PrepareRun().Run(ctx.Done())
}
