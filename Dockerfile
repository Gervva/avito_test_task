# FROM golang:1.22-alpine AS builder
# WORKDIR /app

# COPY go.mod ./
# COPY go.sum ./
# RUN go mod download
# RUN go mod verify

# COPY . .
# RUN go build -o ./bin/service ./cmd/service/main.go

# FROM alpine AS runner

# COPY --from=builder /app/bin/service /

# CMD [ "/service" ]






# FROM golang:1.22-alpine
# WORKDIR /app

# COPY go.mod ./
# COPY go.sum ./
# RUN go mod download
# RUN go mod verify

# COPY . .
# RUN go build -o service ./cmd/service/main.go






# FROM golang:1.22

# WORKDIR /app

# COPY go.mod ./
# COPY go.sum ./
# RUN go mod download

# COPY . .

# RUN go build -o app ./cmd/service/main.go

# CMD ["./app"]






# FROM golang:1.22-alpine as build

FROM golang:1.22-alpine
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go mod verify

COPY . .
RUN go build -o service ./cmd/service






# FROM golang:latest AS builder

# WORKDIR /app

# COPY go.mod .
# COPY go.sum .

# RUN go mod download

# COPY . .

# RUN go build -o main ./cmd/service/

# FROM alpine:latest

# WORKDIR /app

# COPY --from=builder /main .

# ENTRYPOINT ["/app/main"]