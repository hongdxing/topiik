package command

const (
	/*** String ***/
	SET = "SET"
	GET = "GET"

	/*** List ***/
	LPUSH = "LPUSH"
	LPOP  = "LPOP"
	RPUSH = "RPUSH"
	RPOP  = "RPOP"

	/*** CLUSTER ***/
	VOTE         = "VOTE"
	APPEND_ENTRY = "APPENDENTRY"

	/*** OTHERS ***/
	EXPIRE = "EXPIRE"
	DEL    = "DEL"
	TTL    = "TTL"
	QUIT   = "QUIT"
)
