## Service B

### Build

```
$ docker build . -t service-b
```

### Run

You will have to mount `.env` into the image when you run.

```
cp .env.sample .env
```

And fill the `.env` fields. Then run this command.

```
$ docker run -p 5100:5100 -v $(pwd)/.env:/.env service-b
```

## Author 

- KeisukeYamashita