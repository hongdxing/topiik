package command

const (
	/*** CLUSTER ***/
	CLUSTER      = "CLUSTER"
	INIT_CLUSTER = "INIT-CLUSTER"
	SHOW_CLUSTER = "SHOW-CLUSTER"
	ADD_NODE     = "ADD-NODE"
	REMOVE_NODE  = "REMOVE-NODE"
	SCALE        = "SCALE"

	/*** String ***/
	SET  = "SET"
	GET  = "GET"
	SETM = "SETM"
	GETM = "GETM"
	INCR = "INCR"

	/*** List ***/
	LPUSH  = "LPUSH"
	LPOP   = "LPOP"
	LPUSHR = "LPUSHR"
	LPOPR  = "LPOPR"
	LLEN   = "LLEN"

	// List blocking
	LPUSHB  = "LPUSHB"
	LPOPB   = "LPOPB"
	LPUSHRB = "LPUSHRB"
	LPOPRB  = "LPOPRB"

	/*** Key ***/
	EXISTS = "EXISTS"
	KEYS   = "KEYS"
	EXPIRE = "EXPIRE"
	DEL    = "DEL"
	TTL    = "TTL"

	/*** OTHERS ***/
	QUIT = "QUIT"
	GET_CTL_LEADER_ADDR = "GET-CTL-LEADER-ADDR"
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
	KEYS   int16 = 1002
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
