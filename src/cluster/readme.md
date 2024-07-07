# Controller-Workers cluster model #

## Controller ##

Controller is the cluster leader, who in charge of
- Receive client connection and forword to workers
- Maintain health of Captial(CA) and Chief Officer(CO) by sending heartbeat to CO periodically
- Maintain metadata of cluster
    - How many workers
    - Patition info of leader and follower(s)

How first Controller init?
- When create a cluster via command CCSS INIT, the current is assign to Controller
- After Controller created, can join Chief Officer and Worker(s) via command CCSS JOIN --controller host:port


## Workers ##
is is where data stored
- One Workers may have more than one Replica, either Leader of Folower
- Workers must implement vote interface to accept RequestVote from Chief Officer