FROM golang:1.20 AS build
COPY . /src
WORKDIR /src
RUN go build -trimpath -ldflags='-w -s' -o example main.go
FROM busybox
COPY --from=build /src/example /example
CMD [ "/example" ]
