package command

const (
	/*** CLUSTER ***/
	CLUSTER = "CLUSTER"

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
)

var CmdMap = map[string]int16{
	// Cluster: 1 - 15
	CLUSTER: 1,
	// String: 16 - 47 (32)
	SET:  16,
	GET:  17,
	SETM: 18,
	GETM: 19,
	INCR: 20,

	// List: 48 - 79 (32)
	LPUSH:  48,
	LPOP:   49,
	LPUSHR: 50,
	LPOPR:  51,
	LLEN:   52,

	LPUSHB:  53,
	LPOPB:   54,
	LPUSHRB: 55,
	LPOPRB:  56,

	// Set

	// ZSet

	// Hash

	// Geo

	// Event

	// Key:
	EXISTS: 1001,
	KEYS:   1002,
	EXPIRE: 1003,
	DEL:    1004,
	TTL:    1005,

	// Other:
	QUIT: 2001,
}
