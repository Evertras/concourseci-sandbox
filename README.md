# ConcourseCI Sandbox

Playing with [ConcourseCI](https://concourse-ci.org).

As much is from scratch as possible for learning purposes.

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

