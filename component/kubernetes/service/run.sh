#!/bin/bash
declare -A map=(
    ["api-server-url"]=""
    ["base64-yaml-content"]=""
    ["namespace"]=""
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

# TODO if there are "=" in base64 string?

if [ "" = "${map["api-server-url"]}" ]
then
    printf "[COUT] Handle input error: %s\n" "api-server-url"
    printf "[COUT] CO_RESULT = %s\n" "false"
    exit
fi


# TODO support get yaml content from url

yaml=$(echo $YAML |awk '{print}')


echo $yaml | base64 -d > /root/service.yaml

#
if [ "" = "${map["namespace"]}" ]
then
    namespace="default"
   else
    namespace="${map["namespace"]}"
fi

# TODO allow specific namespace
resp_code=$(curl -s -w "%{http_code}"  -o ./response.log -k -X POST  --data-binary @/root/service.yaml -H 'Content-Type:application/yaml' "${map['api-server-url']}/api/v1/namespaces/${namespace}/services")
if [ "$resp_code" -ne 201 ]
then
    cat ./response.log
    printf "[COUT] CO_RESULT = %s\n" "false"
    exit
fi


printf "\n[COUT] CO_RESULT = %s\n" "true"
exit







