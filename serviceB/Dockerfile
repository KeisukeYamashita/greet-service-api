FROM golang:1.11 as build
LABEL MAINTAINER=KeisukeYamashita

RUN mkdir -p /go
ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH

WORKDIR /go/src/github.com/KeisukeYamashita/greet-service-api
ADD . .
RUN GO111MODULE=on go build -v -o service

FROM gcr.io/distroless/base
COPY --from=build /go/src/github.com/KeisukeYamashita/greet-service-api/service /

EXPOSE 5100
CMD ["/service"]