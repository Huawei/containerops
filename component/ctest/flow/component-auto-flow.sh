#!/bin/bash
function read_dir(){
    for file in `ls $1`
		    do
		        if  [ -d $1"/"$file ];then
	            read_dir $1"/"$file
		  elif [[ $file =~ "yml" ]] ;then
			fullpath=$1'/'$file
			echo '\n'$1'/'$file
			if [[ $fullpath =~ "images" ]];then
			continue
			fi
		curl -i -X POST -H 'Content-type':'application/yaml' --data-binary @$fullpath https://flow.opshub.sh/flow/v1/containerops/python_analysis_coala/flow/latest/yaml

		fi
  done
  }
	read_dir containerops/component
