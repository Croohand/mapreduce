#!/usr/bin/env bash

MASTER_PORT=11000
SLAVES_PORT=11001
SLAVES_COUNT=5
MASTER_ADDR="http://localhost:${MASTER_PORT}"
SCHEDULER_PORT=$((${SLAVES_PORT} + $SLAVES_COUNT))
SCHEDULER_ADDR="http://localhost:${SCHEDULER_PORT}"

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
    $(go env GOPATH)/bin/master start -name "master" -port $MASTER_PORT -override -slaves $SLAVES -schedulers $SCHEDULER_ADDR &
    for i in $(seq -w 1 ${SLAVES_COUNT})
    do
        rm -rf "$MR_PATH/slave$i/sources"
        mkdir "$MR_PATH/slave$i/sources"
        cp ../template/build.sh "$MR_PATH/slave$i/sources/"
        cp -r ../template/main "$MR_PATH/slave$i/sources/main"
        $(go env GOPATH)/bin/slave start -name "slave$i" -port "$((${SLAVES_PORT} + $i - 1))" -master $MASTER_ADDR -override &
    done
    $(go env GOPATH)/bin/slave start -name "scheduler" -port $SCHEDULER_PORT -master $MASTER_ADDR -scheduler -override &
fi
