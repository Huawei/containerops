 #!/bin/bash
curl -i -X POST -H 'Content-type':'application/yaml' --data-binary @ci-component.yml  https://flow.opshub.sh/flow/v1/containerops/component-ctest-flow/flow/latest/yaml