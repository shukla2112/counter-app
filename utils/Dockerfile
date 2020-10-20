FROM alpine:3.9 
ENV APP=counter
ENV MODE=production

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go install`

# Use the entrypoint to make the container look like binary
ENTRYPOINT ["echo"]

# Use the command to pass the mode and other params`
CMD ["-mode",  "production"]
