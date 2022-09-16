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

# Install go
RUN wget http://godeb.s3.amazonaws.com/godeb-amd64.tar.gz && \
    tar zxvf godeb-amd64.tar.gz && \
    ./godeb download ${GO_VERSION} && \
    dpkg -i go_${GO_VERSION}-godeb1_amd64.deb && \
    rm godeb-amd64.tar.gz godeb go_${GO_VERSION}-godeb1_amd64.deb

RUN gem install fpm

WORKDIR /work
