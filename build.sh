
NAME=kvlibp2pd
go build -o bin/$NAME ./cmd/$NAME


# if arguments are supplied, call them
if [ $# -gt 0 ]; then
    ./bin/$NAME $@
    exit
fi

./bin/$NAME -h