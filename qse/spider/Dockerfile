# Build stage
FROM golang:1.13
# Set up user to run container to avoid running container as root
RUN useradd container-runner --create-home --shell /bin/bash
# Set up build dir
WORKDIR /build
# Copy in source files
COPY ./ ./
WORKDIR /build/spider
# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o dist/qse-spider ./cmd/qse-spider
# Update owner of binary
RUN chown -R container-runner:container-runner ./dist/*
# Run container as non-root user
USER container-runner

CMD [ "/build/spider/dist/qse-spider" ]
