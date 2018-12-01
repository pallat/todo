FROM golang:1.11 AS build-env

ENV GO111MODULE=on
ENV GOOS=linux
ENV GOARCH=386
ADD . /src
RUN cd /src && go build -o goapp main.go

FROM alpine:3.8
WORKDIR /app
COPY --from=build-env /src/goapp /app/
CMD /app/goapp