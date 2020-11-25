#!/bin/sh

rm -rf /src/*;
set -e;
cd /src;
git clone ${GIT_URL} .;
git fetch origin ${GIT_REF};
git checkout ${GIT_OID};
