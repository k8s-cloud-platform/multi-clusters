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

package cluster

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/klog/v2"

	"github.com/multi-cluster-platform/mcp/pkg/apis/gateway/v1"
)

// REST implements a RESTStorage for Cluster API
type REST struct{}

var _ rest.Connecter = &REST{}
var _ rest.Redirector = &REST{}

// NewREST returns a RESTStorage object that will work against API services.
func NewREST() *REST {
	return &REST{}
}

func (r *REST) NamespaceScoped() bool {
	return false
}

func (r *REST) New() runtime.Object {
	return &v1.Cluster{}
}

// ConnectMethods returns the list of HTTP methods that can be proxied
func (r *REST) ConnectMethods() []string {
	return []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
}

// NewConnectOptions returns versioned resource that represents proxy parameters
func (r *REST) NewConnectOptions() (runtime.Object, bool, string) {
	return &v1.Cluster{}, true, "path"
}

// Connect returns a handler for the websocket connection
func (r *REST) Connect(ctx context.Context, id string, obj runtime.Object, responder rest.Responder) (http.Handler, error) {
	cluster, ok := obj.(*v1.Cluster)
	if !ok {
		return nil, fmt.Errorf("invalid options object: %#v", obj)
	}
	klog.InfoS("handle for cluster rest", "id", id, "cluster.name", cluster.Name, "cluster.path", cluster.Path)

	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		//user, exist := request.UserFrom(req.Context())
		//if !exist {
		//	responsewriters.InternalError(resp, req, errors.New("no user found for request"))
		//	return
		//}
		//req.Header.Set(authenticationv1.ImpersonateUserHeader, requester.GetName())
		//for _, group := range requester.GetGroups() {
		//	if !skipGroup(group) {
		//		req.Header.Add(authenticationv1.ImpersonateGroupHeader, group)
		//	}
		//}
		//req.Header.Set("Authorization", fmt.Sprintf("bearer %s", impersonateToken))
	}), nil
}

// ResourceLocation returns url for resource redirect to
func (r *REST) ResourceLocation(ctx context.Context, id string) (remoteLocation *url.URL, transport http.RoundTripper, err error) {
	return nil, nil, nil
}
