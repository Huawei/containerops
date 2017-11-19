#!/bin/bash
declare -A map=(
    ["git-url"]=""
    ["assembling-url"]=""
    ["registry-url"]=""
    ["namespace"]=""
    ["image"]=""
    ["tag"]=""
    ["username"]=""
    ["password"]=""
    ["authstr"]=""
    ["insecure-registry"]=""
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

if [ "" != "${map["username"]}" -a "" != "${map["password"]}"  ]
then
    link=$(echo -e "{\n\t\"username\": \"${map["username"]}\",\n\t\"password\": \"${map["password"]}\"\n}\n")
    map["authstr"]=`echo "$link" | base64 `
    echo "${map["authstr"]}"
fi

if [ "" = "${map["git-url"]}" ]
then
    printf "[COUT] Handle input error: %s\n" "git-url"
    printf "[COUT] CO_RESULT = %s\n" "false"
    exit 1
fi

git clone ${map["git-url"]}
if [ "$?" -ne "0" ]
then
    printf "[COUT] CO_RESULT = %s\n" "false"
    exit 1
fi

printf "[COUT] Finish git clone, Begin to compress file\n"

pdir=`echo ${map["git-url"]} | awk -F '/' '{print $NF}' | awk -F '.' '{print $1}'`
cd ./$pdir

tar -cf $pdir.tar *
if [ "$?" -ne "0" ]
then
    printf "[COUT] CO_RESULT = %s\n" "false"
    exit 1
fi

printf "[COUT] Finish compress file, Begin to call assembling\n"


if [ "" = "${map["authstr"]}" ]
then
    if  [ "" = "${map["insecure-registry"]}" ]
    then
        resp_code=$(curl -s -w "%{http_code}" -o ./response.log -k -X POST  --data-binary @./$pdir.tar -H 'Content-Type:application/x-tar' "${map['assembling-url']}/assembling/build?registry=${map['registry-url']}&namespace=${map['namespace']}&image=${map['image']}&tag=${map['tag']}")
    else
        resp_code=$(curl -s -w "%{http_code}" -o ./response.log -k -X POST  --data-binary @./$pdir.tar -H 'Content-Type:application/x-tar' "${map['assembling-url']}/assembling/build?registry=${map['registry-url']}&namespace=${map['namespace']}&image=${map['image']}&tag=${map['tag']}&insecure_registry=${map['insecure-registry']}")
    fi
else
    if  [ "" = "${map["insecure-registry"]}" ]
    then
        resp_code=$(curl -s -w "%{http_code}" -o ./response.log -k -X POST  --data-binary @./$pdir.tar -H 'Content-Type:application/x-tar' "${map['assembling-url']}/assembling/build?registry=${map['registry-url']}&namespace=${map['namespace']}&image=${map['image']}&tag=${map['tag']}&authstr=${map['authstr']}")
    else
        resp_code=$(curl -s -w "%{http_code}" -o ./response.log -k -X POST  --data-binary @./$pdir.tar -H 'Content-Type:application/x-tar' "${map['assembling-url']}/assembling/build?registry=${map['registry-url']}&namespace=${map['namespace']}&image=${map['image']}&tag=${map['tag']}&authstr=${map['authstr']}&insecure_registry=${map['insecure-registry']}")
    fi
fi

if [ "$?" -ne "0" ]
then
    printf "[COUT] CO_RESULT = %s\n" "false"
    exit 1
fi
if [ "$resp_code" -ne 200 ]
then
    cat ./response.log
    printf "\n[COUT] RESPONSE_CODE = %s\n" "$resp_code"
    printf "\n[COUT] CO_RESULT = %s\n" "false"
    exit 1
fi


printf "\n[COUT] CO_RESULT = %s\n" "true"
exit