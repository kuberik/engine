#!/bin/sh

set -e;
rm -rf ./*;

depth=50

echo "machine github.com login ${GITHUB_USERNAME} password ${GITHUB_TOKEN}" > ~/.netrc

GIT_URL="https://github.com/${GITHUB_OWNER}/${GITHUB_REPO}.git"

git init
git remote add origin ${GIT_URL}
git fetch --no-tags --depth ${depth} ${GIT_URL} ${GIT_COMMIT_HASH}:refs/remotes/origin/${GIT_BRANCH}
git checkout ${GIT_BRANCH}
