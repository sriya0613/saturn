FROM docker.io/golang:1.25-alpine3.22 AS buildStage

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY types.go ./
COPY timer.go ./
COPY main.go ./

RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o saturn

FROM docker.io/alpine:3.22

WORKDIR /opt
COPY --from=buildStage /app/saturn /opt/saturn

ENTRYPOINT ["/opt/saturn"]
