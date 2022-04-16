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
	"context"
)

type contextKey string

const (
	// ClusterContextKey is the name of cluster using in client http context
	clusterContextKey = contextKey("ClusterName")
	// directContextKey is the command to handle direct http transport, no need to wrapper
	directContextKey = contextKey("Direct")
)

func WithDynamicClusterContext(ctx context.Context, clusterName string) context.Context {
	return context.WithValue(ctx, clusterContextKey, clusterName)
}

func getClusterContext(ctx context.Context) (string, bool) {
	clusterName, ok := ctx.Value(clusterContextKey).(string)
	return clusterName, ok
}

func WithDirectContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, directContextKey, true)
}

func getDirectContext(ctx context.Context) bool {
	direct, ok := ctx.Value(directContextKey).(bool)
	return ok && direct
}
