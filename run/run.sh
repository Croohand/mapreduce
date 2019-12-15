#!/usr/bin/env bash

MASTERS_PORT=11000
MASTERS_COUNT=1
SLAVES_PORT=$((${MASTERS_PORT} + $MASTERS_COUNT))
SLAVES_COUNT=5
SCHEDULERS_PORT=$((${SLAVES_PORT} + $SLAVES_COUNT))
SCHEDULERS_COUNT=1

function get_slaves() {
    RES=""
    for i in $(seq -w 1 ${SLAVES_COUNT})
    do
        RES=$RES",http://localhost:$((${SLAVES_PORT} + $i - 1))"
    done
    echo $RES
}

function get_schedulers() {
    RES=""
    for i in $(seq -w 1 ${SCHEDULERS_COUNT})
    do
        RES=$RES",http://localhost:$((${SCHEDULERS_PORT} + $i - 1))"
    done
    echo $RES
}

function get_masters() {
    RES="http://localhost:$((${MASTERS_PORT} + $1 - 1))"
    for i in $(seq -w 1 ${MASTERS_COUNT})
    do
        if [[ $i != $1 ]]; then
            RES=$RES",http://localhost:$((${MASTERS_PORT} + $i - 1))"
        fi
    done
    echo $RES
}

SLAVES=$(get_slaves)
SLAVES=${SLAVES:1}
SCHEDULERS=$(get_schedulers)
SCHEDULERS=${SCHEDULERS:1}

SUCCESS=0
go install -i ../client && go install -i ../master && go install -i ../slave && go install -i ../simple_logger && SUCCESS=1
if [[ $SUCCESS == 1 ]]
then
    simple_logger start -name "logger" -port 11100 -output "requests.log" -override &
    mkdir $MR_PATH
    for i in $(seq -w 1 ${SLAVES_COUNT})
    do
        mkdir "$MR_PATH/slave$i"
        rm -rf "$MR_PATH/slave$i/sources"
        mkdir "$MR_PATH/slave$i/sources"
        cp ../template/build.sh "$MR_PATH/slave$i/sources/"
        cp -r ../template/main "$MR_PATH/slave$i/sources/main"
        if [[ $1 == 0 ]]; then
            $(go env GOPATH)/bin/slave start -name "slave$i" -port "$((${SLAVES_PORT} + $i - 1))" -logger http://localhost:11100 -override &
        else
            $(go env GOPATH)/bin/slave start -name "slave$i" &
        fi
    done
    for i in $(seq -w 1 ${SCHEDULERS_COUNT})
    do
        if [[ $1 == 0 ]]; then
            $(go env GOPATH)/bin/slave start -name "scheduler$i" -port "$((${SCHEDULERS_PORT} + $i - 1))" -logger http://localhost:11100 -scheduler -override &
        else
            $(go env GOPATH)/bin/slave start -name "scheduler$i" &
        fi
    done
    for i in $(seq -w 1 ${MASTERS_COUNT})
    do
        if [[ $1 == 0 ]]; then
            $(go env GOPATH)/bin/master start -name "master$i" -port $MASTERS_PORT -masters $(get_masters $i) -slaves $SLAVES -schedulers $SCHEDULERS -logger http://localhost:11100 -override &
        else
            $(go env GOPATH)/bin/master start -name "master$i" &
        fi
    done
fi
