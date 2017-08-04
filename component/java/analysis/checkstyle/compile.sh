#!/bin/bash

declare -A map=(
    ["git-url"]="" 
    ["out-put-type"]=""
    ["report-path"]=""
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

if [[ "${map["out-put-type"]}" =~ ^(xml|json|yaml)$ ]]
then
    printf "[COUT] out-put-type: %s\n" "${map["out-put-type"]}"
else
    printf "[COUT] Handle input error: %s\n" "out-put-type should be one of xml,json,yaml"
    printf "[COUT] CO_RESULT = %s\n" "false"
    exit
fi

if [ "" = "${map["report-path"]}" ]
then
    map["report-path"]="build/reports/checkstyle"
fi

git clone ${map["git-url"]} 2>&1 | awk '{print "[COUT] CO_RESULT =", $0}'
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

havecheckstyle=`echo gradle -q tasks --all | grep checkstyle`
if [ "$havecheckstyle" = "" ]
then
    echo -e "\napply plugin: 'checkstyle'" >> build.gradle 2>&1 | awk '{print "[COUT] CO_RESULT =", $0}'
    mkdir -p ./config/checkstyle 2>&1 | awk '{print "[COUT] CO_RESULT =", $0}'
    cp /root/checkstyle.xml ./config/checkstyle/ 2>&1 | awk '{print "[COUT] CO_RESULT =", $0}'
fi

gradle checkstyleMain 2>&1 | awk '{print "[COUT] CO_RESULT =", $0}'
gradle checkstyleTest 2>&1 | awk '{print "[COUT] CO_RESULT =", $0}'

if [ "${map["out-put-type"]}" = "xml" ]
then
    cat ${map["report-path"]}/main.xml 2>&1 | awk '{print "[COUT] CO_RESULT =", $0}'
    cat ${map["report-path"]}/test.xml 2>&1 | awk '{print "[COUT] CO_RESULT =", $0}'
else
    java -jar /root/convert.jar ${map["report-path"]}/main.xml ${map["out-put-type"]} 2>&1 | awk '{print "[COUT] CO_RESULT =", $0}'
    java -jar /root/convert.jar ${map["report-path"]}/test.xml ${map["out-put-type"]} 2>&1 | awk '{print "[COUT] CO_RESULT =", $0}'
fi

if [ "$?" -eq "0" ]
then
    printf "[COUT] CO_RESULT = %s\n" "true"
fi
