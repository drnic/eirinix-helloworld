# Example EiriniX extension

```plain
kubectl apply -f config/deployment-ns-default.yaml
```

To tail the logs of the `app=eirini-helloworld-extension` pod:

```plain
stern -l app=eirini-helloworld-extension
```

When we deploy an example Eirini pod:

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
