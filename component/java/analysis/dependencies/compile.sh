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
    ["version"]=""
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

havedr=`echo gradle -q tasks --all | grep dependencyReport`
if [ "$havedr" = "" ]
then
    echo -e "\nallprojects { apply plugin: 'project-report' }" >> build.gradle
fi
printf "\n[COUT] "
$gradle_version dependencyReport
if [ "$?" -ne "0" ]
then
    printf "[COUT] CO_RESULT = %s\n" "false"
    exit
fi

find ./ -name 'dependencies.txt' | xargs cat
if [ "$?" -eq "0" ]
then
    printf "[COUT] CO_RESULT = %s\n" "true"
else
    printf "[COUT] CO_RESULT = %s\n" "false"
fi
exit