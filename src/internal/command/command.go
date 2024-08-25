package command

const (
	/*** CLUSTER ***/
	INIT_CLUSTER           = "INIT-CLUSTER"
	INIT_CLUSTER_I   uint8 = 1
	SHOW_CLUSTER           = "SHOW-CLUSTER"
	SHOW_CLUSTER_I   uint8 = 2
	ADD_WORKER             = "ADD-WORKER"
	ADD_WORKER_I     uint8 = 3
	ADD_CONTROLLER         = "ADD-CONTROLLER"
	ADD_CONTROLLER_I uint8 = 4
	RM_NODE                = "RM-NODE"
	RM_NODE_I        uint8 = 5
	NEW_PARTITION          = "NEW-PARTITION"
	NEW_PARTITION_I  uint8 = 6
	RM_PARTITION           = "RM-PARTITION"
	RM_PARTITION_I   uint8 = 7
	SCALE                  = "SCALE"
	SCALE_I          uint8 = 8
	GET_CTLADDR            = "GET-CTLADDR"
	GET_CTLADDR_I    uint8 = 15

	/*** String ***/
	SET          = "SET"
	SET_I  uint8 = 16
	GET          = "GET"
	GET_I  uint8 = 17
	SETM         = "SETM"
	SETM_I uint8 = 18
	GETM         = "GETM"
	GETM_I uint8 = 19
	INCR         = "INCR"
	INCR_I uint8 = 20

	/*** List ***/
	LPUSH          = "LPUSH"
	LPUSH_I  uint8 = 32
	LPOP           = "LPOP"
	LPOP_I   uint8 = 33
	LPUSHR         = "LPUSHR"
	LPUSHR_I uint8 = 34
	LPOPR          = "LPOPR"
	LPOPR_I  uint8 = 35
	LLEN           = "LLEN"
	LLEN_I   uint8 = 36

	// List blocking
	LPUSHB  = "LPUSHB"
	LPOPB   = "LPOPB"
	LPUSHRB = "LPUSHRB"
	LPOPRB  = "LPOPRB"

	// Set 48

	// ZSet 64

	// Hash 80

	// Geo 96

	// Event 112
	PUB     = "PUB"
	PUB_I   = 112
	POLL    = "POLL"
	POLL_I  = 113
	POLLB   = "POLLB" // Poll Blocking
	POLLB_I = 114

	/*** Keys 128 ***/
	TTL            = "TTL"
	TTL_I    uint8 = 128
	DEL            = "DEL"
	DEL_I    uint8 = 129
	Keys           = "Keys"
	KEYS_I   uint8 = 130
	EXISTS         = "EXISTS"
	EXISTS_I uint8 = 131
	EXPIRE         = "EXPIRE"

	/*** OTHERS  ***/
	QUIT = "QUIT"
)

/*
// command integer
const (
	// Cluster: 1 - 15
	//CLUSTER: 1,
	INIT_CLUSTER int16 = 1
	SHOW_CLUSTER int16 = 2
	ADD_NODE     int16 = 3
	REMOVE_NODE  int16 = 4
	SCALE        int16 = 5

	// String: 16 - 47 (32)
	SET  int16 = 16
	GET  int16 = 17
	SETM int16 = 18
	GETM int16 = 19
	INCR int16 = 20

	// List: 48 - 79 (32)
	LPUSH  int16 = 48
	LPOP   int16 = 49
	LPUSHR int16 = 50
	LPOPR  int16 = 51
	LLEN   int16 = 52

	LPUSHB  int16 = 53
	LPOPB   int16 = 54
	LPUSHRB int16 = 55
	LPOPRB  int16 = 56

	// Set

	// ZSet

	// Hash

	// Geo

	// Event

	// Key:
	EXISTS int16 = 1001
	Keys   int16 = 1002
	EXPIRE int16 = 1003
	DEL    int16 = 1004
	TTL    int16 = 1005

	// Other:
	QUIT                       int16 = 2001
	GET_CONTROLLER_LEADER_ADDR int16 = 3001
)

var CmdCode = map[string]int16{
	// Cluster: 1 - 15
	S_INIT_CLUSTER: 1,
	S_SHOW_CLUSTER: 2,
	S_ADD_NODE:     3,
	S_REMOVE_NODE:  4,
	S_SCALE:        5,

	// String: 16 - 47 (32)
	S_SET:  16,
	S_GET:  17,
	S_SETM: 18,
	S_GETM: 19,
	S_INCR: 20,

	// List: 48 - 79 (32)
	S_LPUSH:  48,
	S_LPOP:   49,
	S_LPUSHR: 50,
	S_LPOPR:  51,
	S_LLEN:   52,

	S_LPUSHB:  53,
	S_LPOPB:   54,
	S_LPUSHRB: 55,
	S_LPOPRB:  56,

	// Set

	// ZSet

	// Hash

	// Geo

	// Event

	// Key:
	S_EXISTS: 1001,
	S_KEYS:   1002,
	S_EXPIRE: 1003,
	S_DEL:    1004,
	S_TTL:    1005,

	// Other:
	S_QUIT: 2001,
}

*/
