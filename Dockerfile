FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY internal/ ./internal/
COPY cmd/ ./cmd/

RUN CGO_ENABLED=0 GOOS=linux go build -o / github.com/luishfonseca/utcc/cmd/...

FROM scratch as wrapper
COPY --from=builder /uTCCWrapper /uTCCWrapper
CMD ["/uTCCWrapper"]

FROM scratch as coordinator
COPY --from=builder /uTCCCoordinator /uTCCCoordinator
CMD ["/uTCCCoordinator"]

FROM wrapper
