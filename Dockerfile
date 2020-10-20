FROM golang:1.15-alpine AS build_base

ENV APP=counter

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /tmp/${APP}

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Unit tests
RUN CGO_ENABLED=0 go test -v

# Build the Go app
RUN  CGO_ENABLED=0 go build -o ./out/${APP} .

# Start fresh from a smaller image
FROM alpine:3.9 
ENV APP=counter
RUN apk add ca-certificates

# Adding the appuser to have the security
RUN adduser -S -D -H -h /app appuser
WORKDIR /app

COPY --from=build_base /tmp/${APP}/out/${APP} /app/${APP}
COPY config/ /app/config
RUN chown -R appuser /app
USER appuser

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go install`

# Use the entrypoint to make the container look like binary
#ENTRYPOINT ./${APP}
ENTRYPOINT ["./counter"]

# Use the command to pass the mode and other params`
CMD ["-mode",  "dev"]
