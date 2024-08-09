#!/bin/bash

mkdir -p typesense-data/

# this is randomly generated, so it's secure, right? ;)
export TYPESENSE_API_KEY=1667b96f-da3c-40f9-a3b5-8b461a78ed68

./bin/typesense-server --data-dir="$(pwd)"/typesense-data --api-key=$TYPESENSE_API_KEY --enable-cors

