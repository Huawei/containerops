#!/bin/bash
docker build -t containerops/component/component-ctest-build ./build
docker build -t containerops/component/component-ctest-flow ./flow
#docker run \
#	    --rm \
#		    --env CO_DATA=" \
#			    version=gradle3 \
#				    git-url=https://github.com/vanniuner/gradle-demo.git \
#					    out-put-type=json" \
#						    containerops/analysis/java_gradle_checkstyle \`
