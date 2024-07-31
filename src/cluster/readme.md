# Controller-Worker cluster model #

## Controller ##

Controller is the cluster manager, who in charge of
- Receive client connections and forword to workers
- Maintain health of Controllers by sending heartbeat to Followers(Controller follower(s) and workers) periodically
- Maintain metadata of cluster
    - Controllers' information
    - Workers' information
    - Patition information of leader and follower(s)

How first Controller init?
- When create a cluster via command INIT-CLUSTER, the current is assign to Controller
- After Cluster initialized, can add nodes to cluster via command `ADD-NODE host:port controller|worker`


## Workers ##
is is where data stored
- One Workers may have more than one Replica, either Leader of Folower
- Workers must implement vote interface to accept RequestVote from Controller Officer