FROM mcr.microsoft.com/devcontainers/base:bullseye

# Prepare apt
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive

# Install Go
ENV PATH="/usr/local/go/bin:$PATH"
RUN export GO_VERSION="1.21.1" \
    && wget "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" \
    && tar -C /usr/local -xzf "go${GO_VERSION}.linux-amd64.tar.gz" \
    && rm "go${GO_VERSION}.linux-amd64.tar.gz"
