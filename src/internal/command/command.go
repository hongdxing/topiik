package command

const (
	/*** String ***/
	SET   = "SET"
	SETEX = "SETEX"
	GET   = "GET"

	/*** List ***/
	LPUSH = "LPUSH"
	LPOP  = "LPOP"
	RPUSH = "RPUSH"
	RPOP  = "RPOP"

	/*** CLUSTER ***/
	VOTE = "VOTE"

	/*** OTHERS ***/
	EXPIRE = "EXPIRE"
	DEL    = "DEL"
	TTL    = "TTL"
	QUIT   = "QUIT"
)
