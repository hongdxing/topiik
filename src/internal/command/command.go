package command

const (
	// CLUSTER(1-15)
	INIT_CLUSTER           = "INIT-CLUSTER"
	SHOW                   = "SHOW"
	ADD_WORKER             = "ADD-WORKER"
	ADD_CONTROLLER         = "ADD-CONTROLLER"
	REMOVE_NODE            = "REMOVE-NODE"
	NEW_PARTITION          = "NEW-PARTITION"
	REMOVE_PARTITION       = "REMOVE-PARTITION"
	SCALE                  = "SCALE" // to be deleted
	GET_CTLADDR            = "GET-CTLADDR"
	INIT_CLUSTER_I   uint8 = 1
	SHOW_I           uint8 = 2
	ADD_WORKER_I     uint8 = 3
	ADD_CONTROLLER_I uint8 = 4
	REMOVE_NODE_I    uint8 = 5
	NEW_PARTITION_I  uint8 = 6
	RM_PARTITION_I   uint8 = 7
	SCALE_I          uint8 = 8
	GET_CTLADDR_I    uint8 = 15

	// String(16-31)
	SET          = "SET"
	GET          = "GET"
	SETM         = "SETM"
	GETM         = "GETM"
	INCR         = "INCR"
	SET_I  uint8 = 16
	GET_I  uint8 = 17
	SETM_I uint8 = 18
	GETM_I uint8 = 19
	INCR_I uint8 = 20

	// List(32-55)
	LPUSH  = "LPUSH"
	LPOP   = "LPOP"
	LPUSHR = "LPUSHR"
	LPOPR  = "LPOPR"
	LLEN   = "LLEN"
	LRANGE = "LRANGE"
	LSET   = "LSET"

	LPUSH_I  uint8 = 32
	LPOP_I   uint8 = 33
	LPUSHR_I uint8 = 34
	LPOPR_I  uint8 = 35
	LLEN_I   uint8 = 36
	LRANGE_I uint8 = 37
	LSET_I   uint8 = 38

	// List blocking
	LPUSHB  = "LPUSHB"
	LPOPB   = "LPOPB"
	LPUSHRB = "LPUSHRB"
	LPOPRB  = "LPOPRB"

	// Set 56-79

	// ZSet 80-103

	// Hash 104-127

	// Geo 128-151

	// Event 152-175
	PUB     = "PUB"
	PUB_I   = 152
	POLL    = "POLL"
	POLL_I  = 153
	POLLB   = "POLLB" // Poll Blocking
	POLLB_I = 154

	/*** Keys 176 ***/
	TTL            = "TTL"
	DEL            = "DEL"
	Keys           = "Keys"
	EXISTS         = "EXISTS"
	EXPIRE         = "EXPIRE"
	TTL_I    uint8 = 176
	DEL_I    uint8 = 177
	KEYS_I   uint8 = 178
	EXISTS_I uint8 = 179

	/*** OTHERS  ***/
	QUIT = "QUIT"
)
