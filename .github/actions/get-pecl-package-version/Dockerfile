FROM golang:1.26 AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . ./

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /action .

FROM gcr.io/distroless/static-debian13

COPY --from=builder /action /action

ENTRYPOINT ["/action"]
