# Redis libp2p

Goal: Build a mock redis (client and server) using libp2p

Why:
- Ideally learn channels better
- Build a proper client x server in go (only done in Python)
- Learn libp2p for Gordian


```
sh build.sh start

NODE=/ip4/127.0.0.1/tcp/43459/p2p/XXXXXXXXXXXXXXXXXXX

sh build.sh -p $NODE redis set users name Reece
sh build.sh -p $NODE redis get users name

sh build.sh -p $NODE redis set users other AnotherName
sh build.sh -p $NODE redis set table2 userId 1

sh build.sh -p $NODE redis keys users

sh build.sh -p $NODE redis all
sh build.sh -p $NODE redis del users name
```