# Builder container
FROM registry.access.redhat.com/ubi9/go-toolset AS builder

ARG GCFLAGS=""

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

USER 0
RUN make build

# Final runtime image
FROM registry.access.redhat.com/ubi9/ubi-minimal

WORKDIR /app

COPY --from=builder /app/bin/dcm-placement-api .

# TODO: remove when using service providers 
RUN mkdir -p /app/.kube && chown -R 1001:0 /app
USER 1001

# Run the server
EXPOSE 8080
ENTRYPOINT ["/app/dcm-placement-api"]
CMD ["run"]
