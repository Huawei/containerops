#!/bin/bash


function read_dir(){
    for file in `ls $1`
		    do
		        if  [ -d $1"/"$file ];then
	                # echo $1 "--" $file
			  # echo $1 "--" $file is read_dir
	            read_dir $1"/"$file
		  elif [ $file = "Dockerfile" ] ;then
					
   			echo $1
			fullpath=$1
			echo ${1/containerops\/component\//};
			   
			tmpstr=${1/containerops\/component\//};
			echo ${tmpstr//\//\-};
			imagename=${tmpstr//\//\-};
			#tar -cvf ./$fullpath/imagename.tar -C  $fullpath .
		 # echo error
			tar -cvf ./$fullpath.tar -C  $fullpath .
			go run main.go --image $imagename --path ./$fullpath.tar 
		break
		fi
  done
  }
		  #测试目录 test
	read_dir containerops/component




#for file in ./containerops/component/*
#do 
#if test -d $file
#then
#	echo $file is dir
#fi
#done

#while read line 
#do 
#	echo $line
	#压缩
#	tar -cvf .tar -C ./analysis/checkstyle/ .
	#tar -cvf checkstyle.tar -C ./analysis/checkstyle/ .
	#上传 test-java-gradle-testng
#	go main --images $line
#done <component_dir
