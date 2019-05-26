## Service A

### Build

```
$ docker build . -t service-a
```

### Run

You will have to pass enviroment variable `ENV`.

```
$ docker run -e ENV=test service-a
```

## Author 

- KeisukeYamashita