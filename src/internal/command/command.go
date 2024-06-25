package command

const (
	/*** String ***/
	SET  = "SET"
	GET  = "GET"
	SETM = "SETM"
	GETM = "GETM"
	INCR = "INCR"

	/*** List ***/
	LPUSH = "LPUSH"
	LPOP  = "LPOP"
	RPUSH = "RPUSH"
	RPOP  = "RPOP"

	/*** Key ***/
	EXISTS = "EXISTS"
	KEYS   = "KEYS"

	/*** CLUSTER ***/
	VOTE         = "VOTE"
	APPEND_ENTRY = "APPENDENTRY"
	CLUSTER      = "CLUSTER"

	/*** OTHERS ***/
	EXPIRE = "EXPIRE"
	DEL    = "DEL"
	TTL    = "TTL"
	QUIT   = "QUIT"
)
