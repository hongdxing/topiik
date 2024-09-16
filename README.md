# Topiik
A middleware that for both Key/Value store and Message broker

## How Topiik works

### Cluster
Topiik always start with setup a cluster, no matter playground environment or production environment

### Roles

#### Worker
Maintain Topiik cluster, manage worker groups(partitions), execute commands, memory storage of key/value

#### Persistor
Centralized persistence of binary logs, make scale out/in faster by reduce files copy between worker groups


### Simple local cluster
For setup minimum Topiik environment, need two nodes(or use different port number on same host), one Worker node and one Persistor node

![alt text](src/resource/dev_architecture.png)

### Production cluster

![alt text](src/resource/prod_architecture.png)