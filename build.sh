go build -o bin/redislibp2pd ./cmd/redislibp2pd


# if arguments are supplied, call them
if [ $# -gt 0 ]; then
    ./bin/redislibp2pd $@
    exit
fi

./bin/redislibp2pd -h