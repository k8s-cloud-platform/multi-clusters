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

package main

import (
	"context"
	"log"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/multi-cluster-platform/mcp/pkg/wrapper"
)

func main() {
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return
	}

	// ***** This is the place you need to wrap *****
	config.Wrap(wrapper.NewTransport().NewRoundTrip)

	// now we could create and visit all the resources
	client := kubernetes.NewForConfigOrDie(config)

	// list for test
	listNamespace(client, "cluster1")
	listDeployments(client, "cluster1")
}

func listNamespace(client *kubernetes.Clientset, clusterName string) {
	nss, err := client.CoreV1().Namespaces().List(
		wrapper.WithDynamicClusterContext(context.TODO(), clusterName),
		metav1.ListOptions{},
	)
	if err != nil {
		log.Printf("error listing Namespaces: %v", err)
	}

	for _, ns := range nss.Items {
		log.Printf("Namespace: %s\n", ns.Name)
	}
}

func listDeployments(client *kubernetes.Clientset, clusterName string) {
	deploys, err := client.AppsV1().Deployments("").List(
		wrapper.WithDynamicClusterContext(context.TODO(), clusterName),
		metav1.ListOptions{},
	)
	if err != nil {
		log.Printf("error listing Deployments: %v", err)
	}
	for _, deploy := range deploys.Items {
		log.Printf("Deployment: %s\n", deploy.Name)
	}
}
