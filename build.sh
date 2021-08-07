#!/bin/bash

set -e

FAIL=1
OK=0

WORKSPACE=$(cd $(dirname $0) && pwd -P)

GOVERSION="go1.14.4"
function check_go_version() {
    local version=$(go version | awk '{print $3}' | sed s/[[:space:]]//g)
    if [[ "${version}" != "${GOVERSION}" ]]; then
      echo "version not matched. expect: " ${GOVERSION} " actual: " ${version}
      return ${FAIL}
    fi
    return ${OK}
}

function reset_env() {
    rm -rf ${OUTPUT}
    mkdir -p ${OUTPUT}/bin/
    mkdir -p ${OUTPUT}/conf/

    export GOPROXY="https://goproxy.io/,direct"
    export GO111MODULE=auto

    echo -e "go env:\n"
    go env
    echo -e "\n"
}

BuildTime=$(date "+%F %T")
GoVersion=$(go version)
AppName="openapi"
AppVersion=${AppName}"_"$(date "+%F %T" | awk '{print $1"_"$2}')
OUTPUT="output"

function build() {
    go build -o ${APPNAME}

    if [[ $? != 0 ]];then
        echo "compile failed"
        return ${FAIL}
    fi

    mkdir ${OUTPUT}/tools
    cd tools
    for d in $(ls); do
      app=$(basename ${d})
      cd ${WORKSPACE}/tools/${app}
      go build -o ${app}
      mv ${app} ${WORKSPACE}/${OUTPUT}/tools/
      cp *.sh ${WORKSPACE}/${OUTPUT}/tools/
    done
    cd ${WORKSPACE}

    cp -rf ./conf/* ./${OUTPUT}/conf
    cp -rf ./control.sh ./${OUTPUT}
    cp -rf ./${APPNAME} ./${OUTPUT}/bin
    cp supervisord/* ./${OUTPUT}/bin
    mkdir -p ${OUTPUT}/logs

    chmod +x supervisord/supervisord
    chmod +x ${OUTPUT}/*.sh
    chmod +x ${OUTPUT}/bin/${APPNAME}
    chmod +x ${OUTPUT}/tools/*

    cd ${OUTPUT}
    tar -zcf ${APPNAME}.tar.gz *

    #find ./ -type d -name .git |xargs -i rm -rf {}
}

function main() {
    check_go_version
    if [[ $? -ne ${OK} ]]; then
        return ${FAIL}
    fi

    reset_env

    build
    return $?
}

main "$@"
