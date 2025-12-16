#!/usr/bin/env bash

LOGGER_PORT=11100
COORDINATORS_PORT=11000
COORDINATORS_COUNT=3
WORKERS_PORT=$((${COORDINATORS_PORT} + $COORDINATORS_COUNT))
WORKERS_COUNT=5
SCHEDULERS_PORT=$((${WORKERS_PORT} + $WORKERS_COUNT))
SCHEDULERS_COUNT=2

function get_workers() {
    RES=""
    for i in $(seq -w 1 ${WORKERS_COUNT})
    do
        RES=$RES",http://localhost:$((${WORKERS_PORT} + $i - 1))"
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

function get_coordinators() {
    RES="http://localhost:$((${COORDINATORS_PORT} + $1 - 1))"
    for i in $(seq -w 1 ${COORDINATORS_COUNT})
    do
        if [[ $i != $1 ]]; then
            RES=$RES",http://localhost:$((${COORDINATORS_PORT} + $i - 1))"
        fi
    done
    echo $RES
}

WORKERS=$(get_workers)
WORKERS=${WORKERS:1}
SCHEDULERS=$(get_schedulers)
SCHEDULERS=${SCHEDULERS:1}

SUCCESS=0
go install ../client && go install ../coordinator && go install ../worker && go install ../simple_logger && SUCCESS=1
if [[ $SUCCESS == 1 ]]
then
    simple_logger start -name "logger" -port $LOGGER_PORT -output "requests.log" -override &
    mkdir $MR_PATH
    for i in $(seq -w 1 ${WORKERS_COUNT})
    do
        mkdir "$MR_PATH/worker$i"
        rm -rf "$MR_PATH/worker$i/sources"
        mkdir "$MR_PATH/worker$i/sources"
        cp ../template/build.sh "$MR_PATH/worker$i/sources/"
        cp -r ../template/main "$MR_PATH/worker$i/sources/main"
        if [[ $1 == 0 ]]; then
            $(go env GOPATH)/bin/worker start -name "worker$i" -port "$((${WORKERS_PORT} + $i - 1))" -logger http://localhost:$LOGGER_PORT -override &
        else
            $(go env GOPATH)/bin/worker start -name "worker$i" &
        fi
    done
    for i in $(seq -w 1 ${SCHEDULERS_COUNT})
    do
        if [[ $1 == 0 ]]; then
            $(go env GOPATH)/bin/worker start -name "scheduler$i" -port "$((${SCHEDULERS_PORT} + $i - 1))" -logger http://localhost:$LOGGER_PORT -scheduler -override &
        else
            $(go env GOPATH)/bin/worker start -name "scheduler$i" &
        fi
    done
    for i in $(seq -w 1 ${COORDINATORS_COUNT})
    do
        if [[ $1 == 0 ]]; then
            $(go env GOPATH)/bin/coordinator start -name "coordinator$i" -port "$((${COORDINATORS_PORT} + $i - 1))" -coordinators $(get_coordinators $i) -workers $WORKERS -schedulers $SCHEDULERS -logger http://localhost:$LOGGER_PORT -override &
        else
            $(go env GOPATH)/bin/coordinator start -name "coordinator$i" &
        fi
    done
fi
