#!/usr/bin/bash -x

# Copyright 2025 Amim Knabben
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#
# Usage:
#    CLUSTER_NAME=kind ./clusterapi.sh
#

# Import common functions
PWD="$(dirname "$0")"
. "$PWD/common.sh"

# Check for binaries installed on PATH
check_command kind
check_command docker
check_command clusterctl

# The vCenter Username
export VSPHERE_USERNAME=${VSPHERE_USERNAME:-'administrator@vsphere.local'}
# The vCenter Password
export VSPHERE_PASSWORD=${VSPHERE_PASSWORD}
# The vCenter server IP or FQDN
export VSPHERE_SERVER=${VSPHERE_SERVER}
# The vSphere datacenter to deploy the management cluster on
export VSPHERE_DATACENTER=${VSPHERE_DATACENTER}
# The vSphere datastore to deploy the management cluster on
export VSPHERE_DATASTORE=${VSPHERE_DATASTORE}
# The VM network to deploy the management cluster on
export VSPHERE_NETWORK=${VSPHERE_NETWORK}
# The vSphere resource pool for your VMs
export VSPHERE_RESOURCE_POOL=${VSPHERE_RESOURCE_POOL}
# The VM folder for your VMs.
export VSPHERE_FOLDER=${VSPHERE_FOLDER}
# The VM template to use for your VMs already created
export VSPHERE_TEMPLATE=${VSPHERE_TEMPLATE}
# The public SSH-authorized key on all machines
export VSPHERE_SSH_AUTHORIZED_KEY=`cat ~/.ssh/id_rsa.pub`
# The certificate thumbprint for the vCenter server
export VSPHERE_TLS_THUMBPRINT=${VSPHERE_TLS_THUMBPRINT}
# The IP address reserved and used for the control plane endpoint
export CONTROL_PLANE_ENDPOINT_IP=${CONTROL_PLANE_ENDPOINT_IP}
# Set the CPI K8s Version to be installed via CAPI
export CPI_IMAGE_K8S_VERSION=${CPI_IMAGE_K8S_VERSION}
# Set the CAPV default ClusterClass name
export CLUSTER_CLASS_NAME=template

export KUBERNETES_VERSION=${KUBERNETES_VERSION:-1.30.0}

# Initialize the management cluster into kind
clusterctl init --infrastructure vsphere

# Generate the infrastructure manifests for CAPV
clusterctl generate cluster wl-beer 	\
  	--kubernetes-version v1.30.0 		\
  	--control-plane-machine-count=3 	\
  	--worker-machine-count=3			\ 
	--flavor=topology				    \
  > wl-beer.yaml

# Create the cluster object and get the kubeconfig
kubectl apply -f wl-beer.yaml

# Store the kubeconfig on .kube
clusterctl get kubeconfig wl-beer > ${HOME}/.kube/beer-kubeconfig

