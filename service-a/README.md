## Service A

### Build

```
$ docker build . -t service-a
```

### Run

You will have to pass enviroment variable `ENV`.

```
$ docker run -e ENV=test -p 5050:5050 service-a
```

## Author 

- KeisukeYamashita