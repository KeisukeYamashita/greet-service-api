# Greet Service API

## Diagram of the service

![Diagram of architecture](./img/architecture.png)

## NOTE: Create Secret before make deploy

Create this `secret.yaml` file.

```
apiVersion: v1
kind: Secret
metadata:
  name: mysecret
type: Opaque
stringData:
  .env: |-
    SECRET_MESSAGE_PREFIX=Yeah
```

Then apply.

```
kubectl apply -f secret.yaml
```

Now you can deploy by Helm.

## Author

- KeisukeYamashita