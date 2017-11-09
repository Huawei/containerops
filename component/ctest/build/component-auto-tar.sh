#!/bin/bash
mkdir yml tar

function read_dir(){
    for file in `ls $1`
		    do
			fullpath=$1
		        if  [ -d $1"/"$file ];then
						if [[ $fullpath =~ "images" || $fullpath =~ "ctest" ]];then #"java" "ctest"
									continue
									fi
	            read_dir $1"/"$file
		  elif [ $file = "Dockerfile" ] ;then
					
			tmpstr=${1/containerops\/component\//};
			#echo $tmpstr
			tmpstr=${tmpstr#*../../../}
			#echo $tmpstr
			imagename=${tmpstr//\//\-};
			#echo $imagename;

            if [[ $fullpath =~ "python" ]] ;then
			echo "OK"
			else
			continue
			fi
			echo $fullpath
			echo $imagename
			tar -cvf ../../../component/ctest/build/tar/$imagename.tar -C  $fullpath .
		break
		fi
  done
  }
	read_dir ../../../component	

