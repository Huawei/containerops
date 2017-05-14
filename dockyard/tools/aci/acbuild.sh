#!/usr/bin/env bash
set -e

# Build and statically link acserver (this requires it to already be downloaded)
echo "Building Dockyard ACI File..."
CGO_ENABLED=0 GOOS=linux go build -o dockyard -a -tags netgo -ldflags '-w' github.com/Huawei/dockyard

rmDockyard() { 
  rm dockyard; 
}

acbuildEnd() {
  rmDockyard
  export EXIT=$?
  acbuild --debug end && exit $EXIT
}

# When the script exits, remove the binary we built
trap rmDockyard EXIT

# Start the build with an empty ACI
acbuild --debug begin

# In the event of the script exiting, end the build
trap acbuildEnd EXIT

# Name the ACI
acbuild --debug set-name dockyard.sh/appc/genedna/dockyard 

# Copy the binary and its templates into the ACI
acbuild --debug copy dockyard /dockyard
acbuild --debug copy $GOPATH/src/github.com/Huawei/dockyard/conf /conf
acbuild --debug copy $GOPATH/src/github.com/Huawei/dockyard/cert /cert
acbuild --debug copy $GOPATH/src/github.com/Huawei/dockyard/data /data
acbuild --debug copy $GOPATH/src/github.com/Huawei/dockyard/external /external
acbuild --debug copy $GOPATH/src/github.com/Huawei/dockyard/log /log
acbuild --debug copy $GOPATH/src/github.com/Huawei/dockyard/views /views

# Add a mount point for the ACIs to serve
acbuild --debug mount add acis /acis

# Run acserver
acbuild --debug set-exec -- /dockyard web --address 0.0.0.0

# Save the resulting ACI
acbuild --debug write --overwrite dockyard-latest-linux-amd64.aci
