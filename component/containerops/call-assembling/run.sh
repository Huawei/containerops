#!/bin/bash
declare -A map=(
    ["git-url"]=""
    ["assembling-url"]=""
    ["registry-url"]=""
    ["namespace"]=""
    ["image"]=""
    ["tag"]=""
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

git clone ${map["git-url"]}
if [ "$?" -ne "0" ]
then
    printf "[COUT] CO_RESULT = %s\n" "false"
    exit
fi

printf "======finish git clone=====\n"

pdir=`echo ${map["git-url"]} | awk -F '/' '{print $NF}' | awk -F '.' '{print $1}'`
cd ./$pdir
tar -cf $pdir.tar *
if [ "$?" -ne "0" ]
then
    printf "[COUT] CO_RESULT = %s\n" "false"
    exit
fi

printf "====Finish compress file, Begin to call assembling=====\n"

resp_code=$(curl -s -w "%{http_code}" -o /dev/null -k -X POST  --data-binary @./$pdir.tar -H 'Content-Type:application/x-tar' "${map['assembling-url']}/assembling/build?registry=${map['registry-url']}&namespace=${map['namespace']}&image=${map['image']}&tag=${map['tag']}")
if [ "$?" -ne "0" ]
then
    printf "[COUT] CO_RESULT = %s\n" "false"
    exit
fi
if [ "$resp_code" -ne 200 ]
then
    printf "[COUT] CO_RESULT = %s\n" "false"
    exit
fi


printf "\n[COUT] CO_RESULT = %s\n" "true"
exit