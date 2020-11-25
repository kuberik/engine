#!/bin/bash

set -e
set -o pipefail

curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl
chmod +x ./kubectl && mv ./kubectl /usr/bin/kubectl

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
${DIR}/generate-kubeconfig.sh
kubectl config set-context --current --namespace kuberik-ci-e2e

RELEASE_VERSION=v0.16.0
curl -LO https://github.com/operator-framework/operator-sdk/releases/download/${RELEASE_VERSION}/operator-sdk-${RELEASE_VERSION}-x86_64-linux-gnu
mv operator-sdk-${RELEASE_VERSION}-x86_64-linux-gnu /usr/bin/operator-sdk
chmod +x /usr/bin/operator-sdk

make -j2 test
