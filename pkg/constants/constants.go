package constants

// finalizers
const (
	DeployableFinalizer = "deployable/apps.mcp.io"
)

// labels
const (
	DeployableLabelNamespace = "deployable.apps.mcp.io/namespace"
	DeployableLabelName      = "deployable.apps.mcp.io/name"

	/*
	 * apiVersion: apps.mcp.io/v1alpha1
	 * kind: Manifest
	 * metadata:
	 *   name: apps-v1-deployment-my-nginx
	 *   namespace: default
	 *   labels:
	 *     manifest.apps.mcp.io/apiGroup: apps
	 *     manifest.apps.mcp.io/apiVersion: v1
	 *     manifest.apps.mcp.io/kind: Deployment
	 *     manifest.apps.mcp.io/namespace: default
	 *     manifest.apps.mcp.io/name: my-nginx
	 */
	ManifestLabelAPIGroup   = "manifest.apps.mcp.io/apiGroup"
	ManifestLabelAPIVersion = "manifest.apps.mcp.io/apiVersion"
	ManifestLabelKind       = "manifest.apps.mcp.io/kind"
	ManifestLabelNamespace  = "manifest.apps.mcp.io/namespace"
	ManifestLabelName       = "manifest.apps.mcp.io/name"
)
