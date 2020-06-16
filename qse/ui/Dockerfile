# Build stage
FROM golang:1.13
# Set up user to run container to avoid running container as root
RUN useradd container-runner --create-home --shell /bin/bash
# Set up build dir
WORKDIR /build
# Copy in source files
COPY ./ ./
WORKDIR /build/ui
# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o dist/qse-ui ./cmd/qse-ui
RUN cp -r templates/ ./dist/templates/
# Update owner of binary
RUN chown -R container-runner:container-runner ./dist/*
RUN chown -R container-runner:container-runner ./dist/templates*
# Run container as non-root user
USER container-runner

CMD [ "/build/ui/dist/qse-ui" ]
