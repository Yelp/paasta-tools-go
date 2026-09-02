FROM ubuntu:jammy

ARG GO_VERSION

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

# Install go (supports amd64 and arm64)
RUN ARCH=$(dpkg --print-architecture) && \
    wget https://go.dev/dl/go${GO_VERSION}.linux-${ARCH}.tar.gz && \
    tar -C /usr/local -xzf go${GO_VERSION}.linux-${ARCH}.tar.gz && \
    rm go${GO_VERSION}.linux-${ARCH}.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"

RUN gem install fpm

WORKDIR /work
