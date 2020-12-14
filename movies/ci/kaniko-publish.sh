#!/bin/bash

cat > /kaniko/.docker/config.json <<EOF
{
  "auths": {
    "${DOCKER_REGISTRY_SERVER}": {
      "auth": "$(echo -n ${DOCKER_REGISTRY_USERNAME}:${DOCKER_REGISTRY_PASSWORD} | base64)",
      "email": ""
    }
  }
}
EOF

exec /kaniko/executor \
    --context=dir://$(pwd) \
    --destination="${DOCKER_REGISTRY_SERVER}/${GITHUB_OWNER}/${GITHUB_REPO}:$(echo ${GITHUB_COMMIT_HASH} | cut -c-7)"
