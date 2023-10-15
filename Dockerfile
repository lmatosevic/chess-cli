## Install dependencies and build executable file
FROM golang:1.21-bullseye as build

# Create and switch to app directory
WORKDIR /usr/app

# Copy go module definition files
COPY go.mod go.sum ./

# Install go dependencies
RUN go mod download

# Copy application sources
COPY ./ ./

# Configure go cache env variable
ENV GOCACHE=/root/.cache/go-build

# Build executable file using cache from previous builds
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o ./bin/chess-server ./cmd/chess-server


## Run the server
FROM debian:12.1-slim as deploy

# Create and switch to app directory
WORKDIR /usr/app

# Copy server executable file
COPY --from=build /usr/app/bin/chess-server ./

# Copy migrations folder
COPY --from=build /usr/app/migrations ./migrations

# Copy default config file
COPY --from=build /usr/app/config.yaml ./config.yaml

# Expose server port
EXPOSE 64355

# Start the server
ENTRYPOINT ["./chess-server"]
