FROM golang:1.21-alpine AS build
WORKDIR /app

# ensure we can fetch modules and trust TLS
RUN apk add --no-cache git ca-certificates && update-ca-certificates

# use module proxy to avoid git fallback issues
ENV GOPROXY=https://proxy.golang.org,direct

# copy only go.mod (and go.sum if present) to leverage Docker cache
COPY go.mod ./
# run tidy to create go.sum and download modules into cache
RUN go mod tidy

# copy the rest and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./main.go

FROM alpine:3.18
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=build /server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
