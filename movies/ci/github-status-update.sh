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

export GITHUB_REF=${GIT_COMMIT_HASH}
export GITHUB_ACTION=update_state
export GITHUB_CONTEXT="kuberik-ci"

exec github-status-updater
