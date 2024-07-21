package command

const (
	/*** CLUSTER ***/
	S_CLUSTER      = "CLUSTER"
	S_INIT_CLUSTER = "INIT-CLUSTER"
	S_SHOW_CLUSTER = "SHOW-CLUSTER"
	S_ADD_NODE     = "ADD-NODE"
	S_REMOVE_NODE  = "REMOVE-NODE"
	S_SCALE        = "SCALE"

	/*** String ***/
	S_SET  = "SET"
	S_GET  = "GET"
	S_SETM = "SETM"
	S_GETM = "GETM"
	S_INCR = "INCR"

	/*** List ***/
	S_LPUSH  = "LPUSH"
	S_LPOP   = "LPOP"
	S_LPUSHR = "LPUSHR"
	S_LPOPR  = "LPOPR"
	S_LLEN   = "LLEN"

	// List blocking
	S_LPUSHB  = "LPUSHB"
	S_LPOPB   = "LPOPB"
	S_LPUSHRB = "LPUSHRB"
	S_LPOPRB  = "LPOPRB"

	/*** Key ***/
	S_EXISTS = "EXISTS"
	S_KEYS   = "KEYS"
	S_EXPIRE = "EXPIRE"
	S_DEL    = "DEL"
	S_TTL    = "TTL"

	/*** OTHERS ***/
	S_QUIT = "QUIT"
)

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
