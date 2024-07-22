/***
** author: duan hongxing
** data: 14 Jul 2024
** desc:
**
**/

package resp

/*
**	4 bytes of msg length, 1 byte of flag, 1 byte of datatype
** 	datatype: 1:string, 2:integer, 3:string array
 */
const RESPONSE_HEADER_SIZE = 6

const (
	RES_OK                   = "OK"
	RES_NIL                  = "NIL"
	RES_WRONG_ARG            = "WRONG_ARG"
	RES_WRONG_NUMBER_OF_ARGS = "WRONG_NUM_OF_ARGS"
	RES_DATA_TYPE_NOT_MATCH  = "DATA_TYPE_NOT_MATCH"
	RES_SYNTAX_ERROR         = "SYNTAX_ERR"
	RES_KEY_NOT_EXIST        = "KEY_NOT_EXIST"
	RES_KEY_EXIST_ALREADY    = "KEY_EXIST_ALREADY"

	RES_INVALID_OP = "INVALID_OP"

	RES_INVALID_ADDR = "INVALID_ADDR"
)
