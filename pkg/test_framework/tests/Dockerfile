FROM ubuntu:18.04

ENV VER_K3D=v1.6.0
ENV VER_PORT=v1.1.1

RUN DEBIAN_FRONTEND=noninteractive apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    python3 \
    && rm -rf /var/lib/apt/lists/*

RUN curl -sL https://raw.githubusercontent.com/rancher/k3d/$VER_K3D/install.sh \
    -o k3d_install.sh \
    && bash k3d_install.sh \
    && rm k3d_install.sh
RUN curl -sL https://raw.githubusercontent.com/Yelp/ephemeral-port-reserve/$VER_PORT/ephemeral_port_reserve.py \
    -o port_reserve.py
