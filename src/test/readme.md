### Cluster test scenarios

##### stop controller leader, client should able to still functianal
- init-cluster on 8301
- add-worker 8303
- set some key value
- add-controller 8302
- stop 8301
- client reconnect to 8302
- expect client should able to get values


##### add node that already in a cluster, should return error
    `target node already in cluster:`

