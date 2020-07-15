#!/bin/bash

cd ${0%/*}

rm -rf ./keys
mkdir -p ./keys
echo "===> Creating keys..."

docker run -v ${PWD}/keys:/tmp/keys/ --rm concourse/concourse generate-key -t rsa -f /tmp/keys/session_signing_key
docker run -v ${PWD}/keys:/tmp/keys/ --rm concourse/concourse generate-key -t ssh -f /tmp/keys/tsa_host_key
docker run -v ${PWD}/keys:/tmp/keys/ --rm concourse/concourse generate-key -t ssh -f /tmp/keys/worker_key

cp ./keys/worker_key.pub ./keys/authorized_worker_keys

