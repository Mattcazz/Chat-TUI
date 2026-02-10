FROM golang:1.25

WORKDIR /app

COPY server/go.mod server/go.sum ./server/
COPY pkg/go.mod ./pkg/go.mod

WORKDIR /app/server
RUN go mod download

# Copy source code
WORKDIR /app
COPY server ./server
COPY pkg ./pkg

WORKDIR /app/server
RUN go mod tidy

CMD ["make"]
