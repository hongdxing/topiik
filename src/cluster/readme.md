# Controller-Workers cluster model #

## Controller ##

Controller is the cluster leader, who in charge of
- Receive client connection and forword to workers
- Maintain health of Controllers by sending heartbeat to Follower periodically
- Maintain metadata of cluster
    - Controllers' information
    - Workers' information
    - Patition information of leader and follower(s)

How first Controller init?
- When create a cluster via command CCSS INIT, the current is assign to Controller
- After Controller created, can join Controller and Worker(s) via command `CLUSTER JOIN host:port controller|worker`


## Workers ##
is is where data stored
- One Workers may have more than one Replica, either Leader of Folower
- Workers must implement vote interface to accept RequestVote from Controller Officer