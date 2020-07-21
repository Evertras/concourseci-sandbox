# ConcourseCI Sandbox

Playing with [ConcourseCI](https://concourse-ci.org).

As much is from scratch as possible for learning purposes.

These instructions probably won't work out of the box, this is mostly self-reference.

## Run it

```bash
./gen-keys.sh
docker-compose up
```

## Run it in local Kubernetes

Assumes Docker for Desktop, which lets us hit localhost for access rather than
something like minikube which will require some IP shenanigans.

And really this should be using the [actual Helm chart](https://github.com/concourse/concourse-chart)
but because this is for learning funsies we're doing it more from scratch.  This
only targets a local kubernetes cluster with Docker for Desktop.

```bash
# Add 127.0.0.1 for concourse.localhost
sudo vim /etc/hosts

./k8s/install.sh
```

## Credentials

Github and Dockerhub credentials are required for some of the pipelines.  Copy
the sample file and fill it in with actual credentials; this file is in gitignore!

```bash
cp credentials-sample.yaml credentials.yaml
```

## Add SmeeProxy

This repository contains a custom Kubernetes controller for [Smee proxies](https://smee.io).

The controller is built with [kubebuilder](https://kubebuilder.io).

Reference implementation for a controller [found in jetstack on github](https://github.com/jetstack/kubebuilder-sample-controller).

```bash
# Install the Smee proxy resource
cd smeeproxy
make install
```

This allows webhooks to send to Smee.io and then get pushed to our internal infra.

## Crudlib repository

Crudlib is a useless Go library that writes and reads some stuff to Redis.  The
code itself is useless but it does have testable elements we can use to construct
a CI pipe around it.

Requires the SmeeProxy resource type to be installed in order to function.  This
is because the PR resource webhook can't be reached by the outside world, so we
run a Smee proxy to trigger webhooks.  If you don't do this, PRs won't be built
automatically more than once every 24 hours.

```bash
# Install the Smee.io proxy
kubectl apply -f ./crudlib/ci/smee.yaml

# Install the CI pipeline - may need different -t value
fly -t ks set-pipeline -c crudlib/ci/pipeline.yaml -p crudlib -l credentials.yaml
```

