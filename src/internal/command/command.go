package command

const (
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

	/*** CLUSTER ***/
	VOTE             = "VOTE"
	APPEND_ENTRY     = "APPENDENTRY"
	CLUSTER          = "CLUSTER"
	CLUSTER_JOIN_ACK = "__CLUSTER_JOIN_ACK"

	/*** OTHERS ***/
	EXPIRE = "EXPIRE"
	DEL    = "DEL"
	TTL    = "TTL"
	QUIT   = "QUIT"
)
