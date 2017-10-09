#!/bin/bash
function read_dir(){
    for file in `ls $1`
		    do
		        if  [ -d $1"/"$file ];then
	            read_dir $1"/"$file
		  elif [ $file = "Dockerfile" ] ;then
					
			fullpath=$1
			tmpstr=${1/containerops\/component\//};
			#echo $tmpstr
			tmpstr=${tmpstr#*../../../}
			#echo $tmpstr
			imagename=${tmpstr//\//\-};
			#echo $imagename;

			if [[ $fullpath =~ "images" ]];then
			continue
			fi
			#if [[ $fullpath =~ "ctest" ]];then
			#continue
			#fi
			#if [[ $fullpath =~ "java" ]];then
			#continue
			#fi

            if [[ $fullpath =~ "python" ]] ;then
			echo "OK"
			else
			continue
			fi
			echo $fullpath
			echo $imagename
			tar -cvf ./$fullpath.tar -C  $fullpath .
			echo ----------------;
			go run main.go --image $imagename --path ./$fullpath
		break
		fi
  done
  }
	#read_dir containerops/component
	read_dir ../../../component	
