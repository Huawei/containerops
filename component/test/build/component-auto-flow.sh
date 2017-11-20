#!/bin/bash
mkdir yml
tar -cvf yml.tar -C  yml .
curl -XPUT --data-binary @yml.tar https://hub.opshub.sh/binary/v1/containerops/component/binary/v0.1/yml.tar -i
	
