package internal

type Key struct {
	TheKey string
	/***
	* -1: nerver
	* -2: expired
	* >1: seconds to epxire
	 */
	Expire int
}

type StrValue struct {
	Value  string
	Expire int
}
