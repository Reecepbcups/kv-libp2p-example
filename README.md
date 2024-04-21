# KV libp2p Example

A mock Key-Value client and server using libp2p (like redis).

Why:
- Learn [LibP2P](https://libp2p.io/) for [Gordian](https://github.com/rollchains/gordian)
- Build a proper Client x Server in Go (only done in Python)

Manual Testing:
```
sh build.sh start

NODE=/ip4/127.0.0.1/tcp/43459/p2p/XXXXXXXXXXXXXXXXXXX

sh build.sh -p $NODE redis set users name Reece
sh build.sh -p $NODE redis get users name

sh build.sh -p $NODE redis set users other AnotherName
sh build.sh -p $NODE redis set table2 userId 1

sh build.sh -p $NODE redis keys users
sh build.sh -p $NODE redis values users

sh build.sh -p $NODE redis all
sh build.sh -p $NODE redis del users name
```
