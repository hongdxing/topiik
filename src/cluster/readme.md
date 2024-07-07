# Controller-Chief-Workers (CCSS) cluster model #

## Controller(CA) ##

Controller is the cluster leader, who in charge of
- Receive client connection and forword to workers
- Maintain health of Captial(CA) and Chief Officer(CO) by sending heartbeat to CO periodically
- Maintain metadata of cluster
    - How many workers
    - Patition info of leader and follower(s)

How first Controller init?
- When create a cluster via command CCSS INIT, the current is assign to Controller
- After Controller created, can join Chief Officer and Worker(s) via command CCSS JOIN --controller host:port


## Chief officer(CO) ##
Chief Officer is standby of Controller, who in charge of
- Receive heartbeat of Controller
- Receive metadata received from Controller
- Ready to take place Controller anytime
- if not receive heartbeat from Controller in a specific time perior, then CO send RequestVote rpc to all Workers(SA) to promote himself to Controller

## Workers ##
is is where data stored
- One Workers may have more than one Replica, either Leader of Folower
- Workers must implement vote interface to accept RequestVote from Chief Officer