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

package wrapper

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/multi-cluster-platform/mcp/pkg/apis/gateway"
)

const (
	version = "v1"

	pathSeparator = "/"
	pathPrefix    = "/apis"

	pathCluster = "clusters"
	pathShadow  = "shadow"
)

// mcpTransport is a transport for gateway cluster and gateway shadow
type mcpTransport struct {
	delegate http.RoundTripper

	clusterName string
}

var _ http.RoundTripper = &mcpTransport{}

func NewTransport() *mcpTransport {
	return &mcpTransport{}
}

func (t *mcpTransport) NewRoundTrip(delegate http.RoundTripper) http.RoundTripper {
	t.delegate = delegate
	return t
}

func (t *mcpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if getDirectContext(req.Context()) {
		return t.delegate.RoundTrip(req)
	}

	if !t.isFallBack() {
		clusterName, exists := getClusterContext(req.Context())
		if !exists {
			return nil, fmt.Errorf("missing cluster name in the request context")
		}
		t.clusterName = clusterName
	}

	req.URL.Path = t.formatURL(req.URL.Path)
	return t.delegate.RoundTrip(req)
}

// shadow request, send req to hub cluster: http://localhost/apis/gateway.mcp.io/v1/shadow/api/v1/nodes
// cluster request, send req to spoke cluster: http://localhost/apis/gateway.mcp.io/v1/clusters/{name}/api/v1/nodes
func (t *mcpTransport) formatURL(reqPath string) string {
	originalPath := strings.TrimPrefix(reqPath, "/")
	if t.isFallBack() {
		return strings.Join([]string{pathPrefix, gateway.GroupName, version, pathShadow, originalPath}, pathSeparator)
	}
	return strings.Join([]string{pathPrefix, gateway.GroupName, version, pathCluster, t.clusterName, originalPath}, "/")
}

func (t *mcpTransport) isFallBack() bool {
	return t.clusterName == ""
}
