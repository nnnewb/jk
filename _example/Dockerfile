FROM golang:1.20 AS build
ENV GOPROXY=https://goproxy.cn/
RUN mkdir -p "/src"
COPY . /src
WORKDIR /src
RUN go build -trimpath -ldflags="-w -s" -o example main.go

FROM scratch
COPY --from=build /src/example /example
CMD ["/example"]
