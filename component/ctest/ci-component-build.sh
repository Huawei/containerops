#!/bin/bash
docker build -t containerops/component/component-ctest-build ./build
docker build -t containerops/component/component-ctest-flow ./flow