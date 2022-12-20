FROM golang:1.19 AS build

WORKDIR /workspace

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

ENV CGO_ENABLED=1
RUN go build -o main -ldflags='-s -w -extldflags "-static"' main.go

FROM scratch

COPY --from=build /workspace/main /main

ENTRYPOINT [ "/main" ]