# Capital-Chief-Salors (CCSS) cluster model #

## Capital(CA) ##

Capital is the cluster leader, who in charge of
- Receive client connection and forword to salors
- Maintain health of Captial(CA) and Chief Officer(CO) by sending heartbeat to CO periodically
- Maintain metadata of cluster
    - How many salors
    - Patition info of leader and follower(s)

How first Capital init?
- When create a cluster via command CCSS INIT, the current is assign to Capital
- After Capital created, can join Chief Officer and Salor(s) via command CCSS JOIN --capital host:port


## Chief officer(CO) ##
Chief Officer is standby of Capital, who in charge of
- Receive heartbeat of Capital
- Receive metadata received from Capital
- Ready to take place Capital anytime
- if not receive heartbeat from Capital in a specific time perior, then CO send RequestVote rpc to all Salors(SA) to promote himself to Capital

## Salors ##
Salors is where data stored
- One Salors may have more than one Replica, either Leader of Folower
- Salors must implement vote interface to accept RequestVote from Chief Officer