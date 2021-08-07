#!/bin/bash

OK=0
FAIL=1
OK_STR="Sucessful!"
FAIL_STR="Failed!"

APPNAME="openapi"

WORKSPACE=$(cd $(dirname $0) && pwd -P)
CONF=${WORKSPACE}/conf/${APPNAME}.json
SUPERVISORD_CONF=${WORKSPACE}/bin/supervisord.conf

pid=""
supervisord_pid=""

function get_pid() {
    pid=""
    supervisord_pid=""

    PORT=$(grep "http_addr" ${CONF} | awk -F '"' '{print $4}' |awk -F ':' '{print $2}' | sed s/[[:space:]]//g)
    SUPERVISORD_PORT=$(grep "port=" ${SUPERVISORD_CONF} | awk -F ':' '{print $2}' | sed s/[[:space:]]//g)

    # $ss -tlnp
    # LISTEN     0      128         :::8088                    :::*                   users:(("reception",pid=11597,fd=9))
    # LISTEN     0      128         :::9529                    :::*                   users:(("reception",pid=11597,fd=7))
    pid=$(ss -tlnp | grep ${PORT} | grep ${APPNAME} | awk -F ',' '{print $2}' | awk -F '=' '{print $2}')
    if [[ -z ${pid} ]]; then
        return ${FAIL}
    fi

    supervisord_pid=$(ss -tlnp | grep ${SUPERVISORD_PORT} | grep supervisord | awk -F ',' '{print $2}' | awk -F '=' '{print $2}')
    if [[ -z ${supervisord_pid} ]]; then
        return ${FAIL}
    fi
    return ${OK}
}

function start () {
    # support config into multi-env
    bns=""
    for srv in $(get_service_by_host); do
        conf_file="./conf/${srv}.json"
        if [[ -f ${conf_file} ]]; then
            cp "${conf_file}" ./conf/${APPNAME}.json
            bns=${srv}
            break
        fi
    done
    if [[ -z "${bns}" ]]; then
        echo "can not find the target bns on this host"
        return ${FAIL}
    fi

    local start_cmd="cd bin && ./supervisord -c ./supervisord.conf -d"
    
    eval "${start_cmd}"
    if [[ $? -ne 0 ]]; then
        echo "run ${start_cmd} failed"
        return ${FAIL}
    fi
    
    local max_check=15
    local check_cnt=0
    while [[ 1 ]]; do
        get_pid
        if [[ -n ${pid} && -n ${supervisord_pid} ]]; then
            return ${OK}
        fi
        ((check_cnt++))
        if [[ $check_cnt -gt $max_check ]]; then
            break
        fi
        sleep 1
    done
        
    return ${FAIL}
}

function stop() {
    get_pid
    if [[ -z ${pid} && -z ${supervisord_pid} ]]; then
        echo "${APPNAME} has stopped"
        return ${OK}
    fi

    local stop_cmd="cd bin && ./supervisord ctl -u openapi -P openapi -s http://localhost:${SUPERVISORD_PORT} shutdown"
    eval "${stop_cmd}"
    sleep 1 # for wait to release port

    get_pid
    if [[ -z ${pid} && -z ${supervisord_pid} ]]; then
        return ${OK}
    fi

    echo "stop ${APPNAME} failed. using kill -9 to stop..."
    if [[ -n ${supervisord_pid} ]]; then
        kill -9 ${supervisord_pid}
    fi
    if [[ -n ${pid} ]]; then
        kill -9 ${pid}
    fi

    get_pid
    if [[ -z ${pid} && -z ${supervisord_pid} ]]; then
        return ${OK}
    fi

    echo "stop ${APPNAME} failed by kill -9..."
    return ${FAIL}
}

function restart() {
    stop
    if [[ $? -ne ${OK} ]]; then
      echo "stop failed"
      return ${FAIL}
    fi

    start
    if [[ $? -ne ${OK} ]]; then
      echo "start failed"
      return ${FAIL}
    fi

    return ${OK}
}

function reload() {
    return ${OK}
}

function main() {
    if [[ $# -ne 1 ]]; then
      echo "bad param. param must be start|stop|restart|reload"
      return ${FAIL}
    fi

    local param=$1
    param=${param,,}
    local methods=("start" "stop" "restart" "reload")

    for method in "${methods[@]}"; do
        if [[ "${method}" == "${param}" ]]; then
            eval "${method}"
            if [[ $? -eq ${OK} ]]; then
                echo "${method^} ${APPNAME} ${OK_STR}"
                return ${OK}
            else
                echo "${method^} ${APPNAME} ${FAIL_STR}"
                return ${FAIL}
            fi
        fi
    done

    echo "Unknown Command: ${param}"
    return ${FAIL}
}

main "$@"
