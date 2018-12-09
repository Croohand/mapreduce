function get_slaves() {
    RES="http://localhost:11001"
    for i in {2..5}
    do
        RES=$RES",http://localhost:1100$i"
    done
    echo $RES
}

SLAVES=$(get_slaves)

SUCCESS=0
go install -i ../client && go install -i ../master && go install -i ../slave && SUCCESS=1
if [[ $SUCCESS == 1 ]]
then
    $(go env GOPATH)/bin/master start -name "master" -port "11000" -override -slaves ${SLAVES} &
    for i in {1..5}
    do
        $(go env GOPATH)/bin/slave start -name "slave$i" -port "1100$i" -master http://localhost:11000 -override &
    done
fi
