#!/bin/sh

if [ -z ${KUBERIK_SCREENPLAY_RESULT+x} ]; then
    export GITHUB_STATE=pending
    export GITHUB_DESCRIPTION="kuberik CI pipeline running"
else
    if [ "${KUBERIK_SCREENPLAY_RESULT}" == "success" ]; then
        export GITHUB_STATE=success
    else
        export GITHUB_STATE=failure
    fi
    export GITHUB_DESCRIPTION="kuberik CI pipeline finished with ${GITHUB_STATE}"
fi

OWNER_REPO=$(echo ${GIT_URL} | sed -r 's/https:\/\/github.com\/(.+)\/([0-9a-zA-Z_-]+).+/\1\/\2/')
export GITHUB_OWNER=$(echo ${OWNER_REPO} | cut -d'/' -f 1)
export GITHUB_REPO=$(echo ${OWNER_REPO} | cut -d'/' -f 2)
export GITHUB_REF=${GIT_OID}
export GITHUB_ACTION=update_state
export GITHUB_CONTEXT="kuberik-ci"

github-status-updater
