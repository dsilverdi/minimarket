FROM golang:1.15 as builder
# Define build env
ENV GOOS linux
ENV CGO_ENABLED 0
# Add a work directory
WORKDIR /app
# Cache and install dependencies
COPY go.mod go.sum ./
RUN go mod download

ARG SERVICE
# Copy app files
COPY . .
# Build app
RUN go build -o ./bin/${SERVICE} ./cmd/${SERVICE}

FROM alpine:3.14 as production
# Add certificates
ENV TZ=Asia/Jakarta

ARG SERVICE

RUN apk add --no-cache --update  tzdata \
  && cp /usr/share/zoneinfo/${TZ} /etc/localtime \
  && echo ${TZ} > /etc/timezone
  
RUN apk add --no-cache ca-certificates
# Copy built binary from builder
COPY --from=builder /app/bin/${SERVICE} ./service
# Expose portzzz 
EXPOSE 4000
# Exec built binary
ENTRYPOINT ["./service"]