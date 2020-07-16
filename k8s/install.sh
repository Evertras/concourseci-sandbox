#!/bin/bash

cd ${0%/*}

COLOR_OK='\033[0;32m'
COLOR_ERROR='\033[0;31m'
COLOR_RESET='\033[0m'

function log() {
  echo -e "${COLOR_OK}==> ${@}${COLOR_RESET}"
}

function die() {
  echo -e "${COLOR_ERROR}!!!!!!> ${@}${COLOR_RESET}"
  exit 1
}

if ! which helm &> /dev/null || ! helm version | grep v3 &> /dev/null; then
  die 'Helm 3 is required'
fi

log 'Ensuring namespace exists'
kubectl apply -f namespace.yaml || die 'Failed to apply namespace'

log 'Applying Traefik'
kubectl apply -f ./traefik.yaml || die 'Failed to apply Traefik'

log 'Generating keys'
if [ ! -d ../keys ]; then
  ../gen-keys.sh || die 'Failed to generate keys'
else
  log 'Keys already exist, skipping'
fi

log 'Creating web keys secret in k8s'
if ! kubectl get secret web-keys &> /dev/null; then
  kubectl create secret generic web-keys \
    --namespace=concourse \
    --from-file=authorized_worker_keys=../keys/authorized_worker_keys \
    --from-file=session_signing_key=../keys/session_signing_key \
    --from-file=tsa_host_key=../keys/tsa_host_key || die 'Failed to create web-keys'
else
  log Host keys secret already exists, skipping
fi

log 'Creating worker keys secret in k8s'
if ! kubectl get secret worker-keys &> /dev/null; then
  kubectl create secret generic worker-keys \
    --namespace=concourse \
    --from-file=tsa_host_key.pub=../keys/tsa_host_key.pub \
    --from-file=worker_key=../keys/worker_key || die 'Failed to create worker-keys secret'
else
  log 'Worker keys secret already exists, skipping'
fi

log 'Applying main ConcourseCI cluster'
kubectl apply -f concourse.yaml || die 'Failed to apply main ConcourseCI cluster'

