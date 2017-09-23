#!/bin/bash
function read_dir(){
    for file in `ls $1`
		    do
		        if  [ -d $1"/"$file ];then
	            read_dir $1"/"$file
		  elif [ $file = "Dockerfile" ] ;then
					
   			echo $1
			fullpath=$1
			echo ${1/containerops\/component\//};
			   
			tmpstr=${1/containerops\/component\//};
			echo ${tmpstr//\//\-};
			imagename=${tmpstr//\//\-};
			if [[ $fullpath =~ "images" ]];then
			continue
			fi
			tar -cvf ./$fullpath.tar -C  $fullpath .
			go run main.go --image $imagename --path ./$fullpath.tar 
		break
		fi
  done
  }
	read_dir containerops/component
