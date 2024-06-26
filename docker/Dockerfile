# -- multistage docker build: stage #1: build stage
FROM golang:1.19-alpine AS build

RUN mkdir -p /go/src/github.com/consensus-network/consensusd/

WORKDIR /go/src/github.com/consensus-network/consensusd/

RUN apk add --no-cache curl git openssh binutils gcc musl-dev

COPY go.mod .
COPY go.sum .


# Cache consensusd dependencies
RUN go mod download

COPY . .

RUN go build $FLAGS -o consensusd .

# --- multistage docker build: stage #2: runtime image
FROM alpine
WORKDIR /app

RUN apk add --no-cache ca-certificates tini

COPY --from=build /go/src/github.com/consensus-network/consensusd/consensusd /app/
COPY --from=build /go/src/github.com/consensus-network/consensusd/infrastructure/config/sample-consensusd.conf /app/

USER root
ENTRYPOINT [ "/sbin/tini", "--" ]
