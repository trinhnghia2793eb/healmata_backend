# ==========================================
# STAGE 1: BUILDER
# ==========================================
FROM golang:1.26.4-alpine AS builder

# install git to Alpine environment
RUN apk add --no-cache git

# setup working directory in container
WORKDIR /app

# copy libaries management files & install dependencies
# using Docker Cache --> no need to re-install libraries if no change
COPY go.mod go.sum ./
ENV GOPROXY=direct
RUN go mod download

# copy entire source code into container
COPY . .

# build application into execution file (binary) named 'api-server'
# point to project's entry point (main.go)
RUN CGO_ENABLED=0 GOOS=linux go build -o api-server ./cmd/server/main.go

# ==========================================
# STAGE 2: RUNNER
# ==========================================
FROM alpine:latest

# setup timezone & ca-certificates (need for send email SMTP / call external API)
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# copy binary file which was build from STAGE 1 to STAGE 2
COPY --from=builder /app/api-server .

# expose default port mặc định (viper.SetDefault is 8080)
EXPOSE 8080

# run application
CMD ["./api-server"]