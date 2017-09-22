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
			# 如果有yml文件 直接发送给politage
# 先测试单个流程
# 在测试email 流程
#在测试从hub下载flow流程
# 在测试在yml总加入和修改一些元素测试 比如emial 比如 codata 比如 image name

		#	tar -cvf ./$fullpath.tar -C  $fullpath .
		#	go run main.go --image $imagename --path ./$fullpath.tar 
		curl -i -X POST -H 'Content-type':'application/yaml' --data-binary @$fullpath https://flow.opshub.sh/flow/v1/containerops/python_analysis_coala/flow/latest/yaml

		fi
  done
  }
	read_dir containerops/component
