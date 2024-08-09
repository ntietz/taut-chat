#!/bin/bash
# Installs Typesense into the local bin directory, if it does not alrady exist.

mkdir -p bin

if [[ ! -f bin/typesense-server ]]; then
    cd bin
    curl -O https://dl.typesense.org/releases/26.0/typesense-server-26.0-linux-amd64.tar.gz
    tar xvf typesense-server-26.0-linux-amd64.tar.gz
    cd -
fi
