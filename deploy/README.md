# How to Install



## OCM

two ways to install OCM

### 1. use clusteradm to install

- VERSION=v0.1.0 && curl -L -o clusteradm_linux_amd64.tar.gz https://github.com/open-cluster-management-io/clusteradm/releases/download/$VERSION/clusteradm_linux_amd64.tar.gz && tar -zxvf clusteradm_linux_amd64.tar.gz

- init cluster in hub cluster, use command: clusteradm init
- join cluster in spoke cluster, use command: clusteradm join --hub-token <hub-token> --hub-apiserver <hub-apiserver> --cluster-name <cluster-name>

- approve request, use command in hub cluster: clusteradm accept --clusters <cluster-name>

### 2. direct to install

https://github.com/open-cluster-management-io/registration/blob/main/README.md

