ARG CONSENSUSD_IMAGE
ARG CONSENSUSMINER_IMAGE

FROM ${CONSENSUSD_IMAGE} as consensusd
FROM ${CONSENSUSMINER_IMAGE} as consensusminer

FROM golang:1.19-alpine

RUN mkdir -p /go/src/github.com/consensus-network/consensusd/

WORKDIR /go/src/github.com/consensus-network/consensusd/

RUN apk add bash build-base git

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

COPY --from=consensusd /app/ /app/
COPY --from=consensusminer /app/ /app/
ENV PATH="/app:${PATH}"

WORKDIR /go/src/github.com/consensus-network/consensusd/stability-tests

RUN go install ./...

ENTRYPOINT ["./run/run.sh"]
