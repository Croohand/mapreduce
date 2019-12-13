#!/usr/bin/env bash

MASTERS_PORT=11000
MASTERS_COUNT=3
SLAVES_PORT=$((${MASTERS_PORT} + $MASTERS_COUNT))
SLAVES_COUNT=5
SCHEDULER_PORT=$((${SLAVES_PORT} + $SLAVES_COUNT))
SCHEDULER_ADDR="http://localhost:${SCHEDULER_PORT}"

function get_slaves() {
    RES=""
    for i in $(seq -w 1 ${SLAVES_COUNT})
    do
        RES=$RES",http://localhost:$((${SLAVES_PORT} + $i - 1))"
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

SUCCESS=0
go install -i ../client && go install -i ../master && go install -i ../slave && SUCCESS=1
if [[ $SUCCESS == 1 ]]
then
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
    if [[ $1 == 0 ]]; then
        $(go env GOPATH)/bin/slave start -name "scheduler1" -port $SCHEDULER_PORT -logger http://localhost:11100 -scheduler -override &
        $(go env GOPATH)/bin/slave start -name "scheduler2" -port $(($SCHEDULER_PORT + 1)) -logger http://localhost:11100 -scheduler -override &
        $(go env GOPATH)/bin/master start -name "master1" -port $MASTERS_PORT -masters $(get_masters 1) -slaves $SLAVES -schedulers $SCHEDULER_ADDR -logger http://localhost:11100 -override &
        $(go env GOPATH)/bin/master start -name "master2" -port $(($MASTERS_PORT+1)) -masters $(get_masters 2) -slaves $SLAVES -schedulers $SCHEDULER_ADDR -logger http://localhost:11100 -override &
        $(go env GOPATH)/bin/master start -name "master3" -port $(($MASTERS_PORT+2)) -masters $(get_masters 3) -slaves $SLAVES -schedulers $SCHEDULER_ADDR -logger http://localhost:11100 -override &
    else
        $(go env GOPATH)/bin/slave start -name "scheduler1" &
        $(go env GOPATH)/bin/slave start -name "scheduler2" &
        $(go env GOPATH)/bin/master start -name "master1" &
        $(go env GOPATH)/bin/master start -name "master2" &
        $(go env GOPATH)/bin/master start -name "master3" &
    fi
fi
