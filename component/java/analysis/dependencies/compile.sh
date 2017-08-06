#!/bin/bash

exec 2>/tmp/error_out
function STOUT(){
    echo "[COUT] $@"
    $@ | awk '{print "[COUT]", $0}'
    if [ "`cat /tmp/error_out`" = "" ]
    then
        return 0
    else
        awk '{print "[COUT]", $0}' /tmp/error_out
        echo '' > /tmp/error_out
        return 1
    fi
}

declare -A map=(
    ["git-url"]="" 
)
data=$(echo $CO_DATA |awk '{print}')
for i in ${data[@]}
do
    temp=$(echo $i |awk -F '=' '{print $1}')
    value=$(echo $i |awk -F '=' '{print $2}')
    for key in ${!map[@]}
    do
        if [ "$temp" = "$key" ]
        then
            map[$key]=$value
        fi
    done
done
if [ "$?" -ne "0" ]
then
    printf "[COUT] CO_RESULT = %s\n" "false"
fi

if [ "" = "${map["git-url"]}" ]
then
    printf "[COUT] Handle input error: %s\n" "git-url"
    printf "[COUT] CO_RESULT = %s\n" "false"
    exit
fi

STOUT git clone ${map["git-url"]}
pdir=`echo ${map["git-url"]} | awk -F '/' '{print $NF}' | awk -F '.' '{print $1}'`
cd ./$pdir
if [ ! -f "build.gradle" ]
then
    printf "[COUT] CO_RESULT = file build.gradle not found! \n"
    printf "[COUT] CO_RESULT = %s\n" "false"
    exit
fi

STOUT gradle dependencies

if [ "$?" -eq "0" ]
then
    printf "[COUT] CO_RESULT = %s\n" "true"
fi
