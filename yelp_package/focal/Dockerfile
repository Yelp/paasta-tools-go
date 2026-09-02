FROM ubuntu:focal

ARG GO_VERSION
ARG GO_SHA256

RUN apt-get -q update && \
    DEBIAN_FRONTEND=noninteractive apt-get -q install -y --no-install-recommends \
        ca-certificates \
        git \
        wget \
        ruby \
        ruby-dev \
        rubygems \
        build-essential \
    && apt-get -q clean

# Install go from go.dev with checksum verification (supports amd64 and arm64)
RUN ARCH=$(dpkg --print-architecture) && \
    wget -q https://go.dev/dl/go${GO_VERSION}.linux-${ARCH}.tar.gz -O /tmp/go.tar.gz && \
    echo "${GO_SHA256}  /tmp/go.tar.gz" | sha256sum -c - && \
    tar -C /usr/local -xzf /tmp/go.tar.gz && \
    rm /tmp/go.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"

RUN gem install fpm

WORKDIR /work
