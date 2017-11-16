#!/bin/bash

declare -A map=(
    ["api-server-url"]=""
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

# if namespace is set, use the namespace
if [ "" = "${map["namespace"]}" ]
then
    namespace="default"
else
    namespace="${map["namespace"]}"
    # Create the namespace
    createns=$(kubectl --server="${map["api-server-url"]}" create namespace ${namespace}  >/dev/null 2>&1)
fi

yaml=$(echo $YAML |awk '{print}')
echo $yaml | base64 -d > /root/template.yaml

# Before create yaml, clean it
clean=$(kubectl --server="${map["api-server-url"]}" delete -f /root/template.yaml -n ${namespace} >/dev/null 2>&1)

kubectl --server="${map["api-server-url"]}" create -f /root/template.yaml -n  ${namespace}
if [ "$?" -ne "0" ]
then
    printf "[COUT] CO_RESULT = %s\n" "false"
    exit 1
fi
printf "\n[COUT] CO_RESULT = %s\n" "true"
exit






