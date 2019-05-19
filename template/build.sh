export GOPATH="$1" && rm -rf $1/bin && mkdir $1/bin && go build -o $1/bin/main $1/src/main
