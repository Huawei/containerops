#!/bin/bash

function STOUT(){
    $@ 1>/tmp/standard_out 2>/tmp/error_out
    if [ "$?" -eq "0" ]
    then
        cat /tmp/standard_out | awk '{print "[COUT]", $0}'
        cat /tmp/error_out | awk '{print "[COUT]", $0}' >&2
        return 0
    else
        cat /tmp/standard_out | awk '{print "[COUT]", $0}'
        cat /tmp/error_out | awk '{print "[COUT]", $0}' >&2
        return 1
    fi
}
function STOUT2(){
    $@ 1>/dev/null 2>/tmp/error_out
    if [ "$?" -eq "0" ]
    then
        cat /tmp/error_out | awk '{print "[COUT]", $0}' >&2
        return 0
    else
        cat /tmp/error_out | awk '{print "[COUT]", $0}' >&2
        return 1
    fi
}

declare -A map=(
    ["git-url"]=""
    ["target"]=""
    ["version"]=""
    ["build-path"]=""
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

if [ "" = "${map["git-url"]}" ]
then
    printf "[COUT] Handle input error: %s\n" "git-url"
    printf "[COUT] CO_RESULT = %s\n" "false"
    exit
fi

if [ "${map["version"]}" = "gradle3" ]
then
    gradle_version=$gradle3/gradle
elif [ "${map["version"]}" = "gradle4" ]
then
    gradle_version=$gradle4/gradle
else
    printf "[COUT] Handle input error: %s\n" "version should be one of gradle3,gradle4"
    printf "[COUT] CO_RESULT = %s\n" "false"
    exit
fi

if [ "" = "${map["target"]}" ]
then
    printf "[COUT] Handle input error: %s\n" "target"
    printf "[COUT] CO_RESULT = %s\n" "false"
    exit
fi

STOUT git clone ${map["git-url"]}
if [ "$?" -ne "0" ]
then
    printf "[COUT] CO_RESULT = %s\n" "false"
    exit
fi
pdir=`echo ${map["git-url"]} | awk -F '/' '{print $NF}' | awk -F '.' '{print $1}'`
cd ./$pdir
if [ ! -f "build.gradle" ]
then
    printf "[COUT] CO_RESULT = file build.gradle not found! \n"
    printf "[COUT] CO_RESULT = %s\n" "false"
    exit
fi

if [ "" = "${map["build-path"]}" ]
then
    map["build-path"]="./"
fi
cd ${map["build-path"]}
cat /root/javadoc.conf >> build.gradle

$gradle_version javadoc
if [ "$?" -ne "0" ]
then
    printf "[COUT] CO_RESULT = %s\n" "false"
    exit
fi

cd build/docs
tar -zcvf /root/javadoc.tar javadoc 1>/dev/null 2>&1
if [ "$?" -ne "0" ]
then
    printf "[COUT] CO_RESULT = %s\n" "false"
    exit
fi

curl -i -X PUT -T /root/javadoc.tar ${map["target"]} 2>/dev/null

if [ "$?" -eq "0" ]
then
    printf "\n[COUT] download javadoc url : %s\n" ${map["target"]}
    printf "[COUT] CO_RESULT = %s\n" "true"
else
    printf "\n[COUT] upload javadoc to %s fail %s\n" ${map["target"]}
    printf "[COUT] CO_RESULT = %s\n" "false"
fi
exit