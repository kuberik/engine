#!/bin/bash
set -e
set -o pipefail

SERVICE_ACCOUNT_DIR="/var/run/secrets/kubernetes.io/serviceaccount"

set_kube_config_values() {
    CLUSTER_NAME=test
    echo -n "Setting cluster '${CLUSTER_NAME}' in kubeconfig..."
    kubectl config set-cluster "${CLUSTER_NAME}" \
    --server="https://kubernetes.default.svc" \
    --certificate-authority="${SERVICE_ACCOUNT_DIR}/ca.crt" \
    --embed-certs=true

    CLUSTER_USER=test
    echo -n "Setting token credentials entry in kubeconfig..."
    kubectl config set-credentials "${CLUSTER_USER}" \
    --token="$(cat ${SERVICE_ACCOUNT_DIR}/token)"

    CLUSTER_CONTEXT="test"
    echo -n "Setting a context '${CLUSTER_CONTEXT}' in kubeconfig..."
    kubectl config set-context "${CLUSTER_CONTEXT}" \
    --cluster="${CLUSTER_NAME}" \
    --user="${CLUSTER_USER}"

    kubectl config use-context "${CLUSTER_CONTEXT}"
}

mkdir -p ${HOME}/.kube
set_kube_config_values

echo -e "\\nAll done!"
