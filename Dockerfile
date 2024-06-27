# Step 1: Modules caching
FROM golang:alpine as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:alpine as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bin/app .


# Step 3: Final
FROM scratch
WORKDIR gopractice
COPY --from=builder /bin/app /app 
ENV GIN_MODE=release
ENV PORT=8080
CMD ["/app"]

