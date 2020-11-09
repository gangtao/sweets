# sweets
a configuration service with rest api interface and etcd/zk backend

# how to run it
1. install swag by running `go get -u github.com/swaggo/swag/cmd/swag`
2. run `make doc` or `swag init` to generate all the rest api document
3. run `make docker` will build the docker image locally.
4. run `docker-compose up` will start the sweets server with both `zookeeper` and `etcd` backend. etcd is the default backend, set `KV_TYPE` to `ZK` to switch to zookeeper backend.
5. in `client` folder, run `go build` to build the client and run `./client` to run the test.

# perfornance test
k6 is used to run the performance test. to test the server performance.
1. start the server in `tests/single-node-etcd` or `tests/single-node-zk` with etcd or zookeeper backend by `docker-compose up`
2. in `tests` folder, run `k6 run -i iterations -u concurrent_users test.js` to run the performance test 

# related 3rd party tools
- Gin GO rest framwork refer to https://godoc.org/github.com/gin-gonic/gin
- gin-swagger for REST API document generation refer to  https://github.com/swaggo/gin-swagger
- Zookeeper client refer to https://github.com/samuel/go-zookeeper
- ETCD client refer to https://godoc.org/gopkg.in/coreos/etcd.v2/clientv3 
- K6 a performance test tool refer to https://github.com/loadimpact/k6