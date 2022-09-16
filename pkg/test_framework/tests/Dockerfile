FROM python:3.7-alpine

ENV VER_K3D=v5.4.6
ENV VER_PORT=v1.1.4

RUN apk --no-cache add ca-certificates bash curl

RUN curl -sL https://raw.githubusercontent.com/rancher/k3d/$VER_K3D/install.sh \
    -o k3d_install.sh \
    && bash k3d_install.sh \
    && rm k3d_install.sh
RUN curl -sL https://raw.githubusercontent.com/Yelp/ephemeral-port-reserve/$VER_PORT/ephemeral_port_reserve.py \
    -o port_reserve.py
