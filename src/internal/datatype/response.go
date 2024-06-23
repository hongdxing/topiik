/***
**
**
**
**/

package datatype

type StrResponse struct {
	R bool   // Success or not
	M string // Message: The Value if S is true, else the Error message
}

type IntegerResponse struct {
	R bool
	M int
}

type ListResponse struct {
	R bool
	M []string
}
