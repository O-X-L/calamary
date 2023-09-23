#!/bin/bash

cd "$(dirname "$0")/../docs"

rm -rf build
mkdir build

sphinx-build -b html source/ build/
