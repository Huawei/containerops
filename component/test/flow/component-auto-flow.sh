#!/bin/bash
curl -XGET -O  https://hub.opshub.sh/binary/v1/containerops/component/binary/v0.1/yml.tar
mkdri yml
tar -zxvf yml.tar -C yml
function read_dir(){
    for file in `ls $1`
		    do
				fullpath=$1'/'$file
				echo $fullpath
				if [[ $file =~ "yml" ]] ;then
				echo $file
				curl -i -X POST -H 'Content-type':'application/yaml' --data-binary @$fullpath https://flow.opshub.sh/flow/v1/containerops/python_analysis_coala/flow/latest/yaml
				fi
			
  done
  }
	read_dir yml