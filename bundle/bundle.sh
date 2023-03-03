#!/bin/bash

pushd bundle

opa build -b apikeys --optimize=1 --entrypoint "apikeys/allow"

popd