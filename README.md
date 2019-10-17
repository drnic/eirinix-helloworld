# Example Hello World EiriniX extension

Eirini is an extension for Cloud Foundry that runs applications as statefulsets/pods within Kubernetes. This in turn allows operators to extend the behavior of their own Cloud Foundry platform with mutating webhooks - when a pod is created by Eirini, our bespoke webhook can be run to mutate the pod before it commences. 

One convenient starting point for writing mutating webhooks for Eirini-managed Pods is the SUSE [EiriniX](https://github.com/SUSE/eirinix) library (see [blog post](https://www.cloudfoundry.org/blog/introducing-eirinix-how-to-build-eirini-extensions/)).

In this project we have a mutating webhook that adds an environment variable `STICKY_MESSAGE` to each Eirini pod.

Plus, we start by demonstrating how it works without you needing to run Cloud Foundry/Eirini.

This EiriniX extension should work with any Cloud Foundry/Eirini, including:

* [IBM Cloud Foundry Enterprise Edition](https://cloud.ibm.com/docs/cloud-foundry?topic=cloud-foundry-getting-started) with "Eirini Technical Preview" enabled
* [Pivotal Application Service for Kubernetes](https://pivotal.io/platform/pas-on-kubernetes)
* [SUSE Cloud Application Platform](https://www.suse.com/products/cloud-application-platform/cloud-foundry/), with [Eirini Enabled](https://documentation.suse.com/suse-cap/1/html/cap-guides/cha-cap-depl-eirini.html#sec-cap-eirini-enable).
* Any open source Cloud Foundry with [Eirini Release included](https://documentation.suse.com/suse-cap/1/html/cap-guides/cha-cap-depl-eirini.html#sec-cap-eirini-enable).
* Stark & Wayne's [Bootstrap Kubernetes Demos](https://documentation.suse.com/suse-cap/1/html/cap-guides/cha-cap-depl-eirini.html#sec-cap-eirini-enable) with `--scf` flag to run Cloud Foundry/Eirini on your Kubernetes cluster.

This repo is a rewrite of [EiriniX Sample](https://github.com/SUSE/eirinix-sample/) for the benefit of my comprehension of what was going on. I think it is also a simpler implementation and simpler deployment story.

## Demonstration without Eirini on any Kubernetes

We can see this Hello World EiriniX extension running without Cloud Foundry nor Eirini itself. It will deploy the webhook and sample application into the `default` namespace. We will clean up at the end.

### Deploy Hello World Webhook

```plain
kubectl apply -f config/deployment-ns-default.yaml
```

To tail the logs of the `app=eirini-helloworld-extension` pod:

```plain
stern -l app=eirini-helloworld-extension
```

### Deploy Sample Pod

When we deploy an example Eirini-like pod (EiriniX webhooks match to any pod with a label `source_type: APP`):

```plain
kubectl apply -f config/fakes/eirini_app.yaml
```

The logs of the extension will show:

```plain
eirini-helloworld-extension {"level":"info","ts":1571206202.8373592,"logger":"Hello world!","caller":"hello/helloworld.go:33","msg":"Hello from my Eirini extension! Eirini application POD: eirini-fake-app (default)"}
```

If we look at the fake Eirini app pod, we see that it has been mutated to include an additional environment variable `STICKY_MESSAGE`:

```plain
$ kubectl describe pod eirini-fake-app
...
    Environment:
      STICKY_MESSAGE:  Eirinix is awesome!
...
```

### Tear it down

In addition to deleting the fake Eirini app and the webhook/service, we also need to delete a generated `mutatingwebhookconfiguration` and a `secret`:

```plain
kubectl delete -f config/fakes/eirini_app.yaml
kubectl delete -f config/deployment-ns-default.yaml
kubectl delete mutatingwebhookconfigurations eirini-x-drnic-helloworld-mutating-hook-default
kubectl delete secret eirini-x-drnic-helloworld-setupcertificate
```

### Separate namespaces for Webhook and app Pods

It might be common to see the Cloud Foundry/Eirini components running in one namespace, say `scf`, and the Eirini-managed application Pods in another `scf-eirini`.

We can deploy our webhook into `scf`, and watch/mutate pods in `scf-eirini`:

```plain
kubectl apply -f config/deployment-ns-scf.yaml
```

In the `Deployment` spec for the webhook container we can see that we set `POD_NAMESPACE` to `scf-eirini`, and `WEBHOOK_NAMESPACE` to `scf`:

```yaml
containers:
- image: drnic/eirinix-sample:latest
  name: eirini-helloworld-extension
  imagePullPolicy: IfNotPresent
  env:
  - name: WEBHOOK_SERVICE_NAME
    value: eirini-helloworld-extension-service
  - name: WEBHOOK_NAMESPACE
    value: scf
  - name: POD_NAMESPACE
    value: scf-eirini
```

To deploy our fake app into `scf-eirini`:

```plain
kubectl apply -f config/fakes/eirini_app.yaml -n scf-eirini
```

As above, the resulting Pod will be mutated to include a `STICKY_MESSAGE` environment variable.

To clean up our two-namespace demo:

```plain
```plain
kubectl delete -f config/fakes/eirini_app.yaml -n scf-eirini
kubectl delete -f config/deployment-ns-scf.yaml

kubectl delete mutatingwebhookconfigurations eirini-x-drnic-helloworld-mutating-hook-scf-eirini
kubectl delete secret eirini-x-drnic-helloworld-setupcertificate -n scf-eirini
```

## Developers

Our EiriniX webhook above is installed as a Kubernetes deployment, so we need to package this source code as an OCI/Docker image. There is no `Dockerfile`. I'm trying to ween myself off them and the mystery meat they may contain.

Instead we will use [Cloud Native Buildpacks](https://buildpacks.io), and the [`pack` CLI](https://buildpacks.io/docs/install-pack/):

```plain
pack build drnic/eirinix-sample --builder cloudfoundry/cnb:bionic --publish
```

One way to "update" the webhook is to delete and re-apply the service:

```plain
kubectl delete -f config/deployment-ns-default.yaml
kubectl apply -f config/deployment-ns-default.yaml
```

Any previously generated `mutatingwebhookconfiguration` and `secret` will be kept (these must be explicitly deleted to be removed).

See [Tear it down](#tear-it-down) for cleanup of webhook, secret, and `mutatingwebhookconfiguration`.
