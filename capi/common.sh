#!/bin/bash
#
## Copyright 2025 Amim Knabben
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

set -o errexit
set -o nounset
set -o pipefail

CLUSTER_NAME=${CLUSTER_NAME:-kind}

function check_command() {
    if ! [ -x "$(command -v $1)" ]; then
        echo "The $1 binary is not installed."
        exit 1
    fi
}

function run_docker() {
    docker exec ${CLUSTER_NAME}-control-plane $@
}



