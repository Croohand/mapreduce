#!/usr/bin/env bash

MASTER_PORT=11000
SLAVES_PORT=11001
SLAVES_COUNT=5

function get_slaves() {
    for i in $(seq -w 1 ${SLAVES_COUNT})
    do
        RES=$RES",http://localhost:$((${SLAVES_PORT} + $i - 1))"
    done
    echo $RES
}

SLAVES=$(get_slaves)
SLAVES=${SLAVES:1}

SUCCESS=0
go install -i ../client && go install -i ../master && go install -i ../slave && SUCCESS=1
if [[ $SUCCESS == 1 ]]
then
    $(go env GOPATH)/bin/master start -name "master" -port "${MASTER_PORT}" -override -slaves ${SLAVES} &
    for i in $(seq -w 1 ${SLAVES_COUNT})
    do
        $(go env GOPATH)/bin/slave start -name "slave$i" -port "$((${SLAVES_PORT} + $i - 1))" -master "http://localhost:${MASTER_PORT}" -override &
    done
fi
