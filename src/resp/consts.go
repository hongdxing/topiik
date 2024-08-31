/***
** author: duan hongxing
** data: 14 Jul 2024
** desc:
**
**/

package resp

/*
*	4 bytes of msg length, 1 byte of flag, 1 byte of datatype
* 	datatype: 1:string, 2:integer, 3:string array
 */
const RESPONSE_HEADER_SIZE = 6

const (
	RES_OK                   = "OK"
	RES_NIL                  = "NIL"
	RES_WRONG_ARG            = "WRONG_ARG"
	RES_WRONG_NUMBER_OF_ARGS = "WRONG_NUM_OF_ARGS"
	RES_DATA_TYPE_MISMATCH   = "DT_MISMATCH"
	RES_SYNTAX_ERROR         = "SYNTAX_ERR"
	RES_EMPTY_KEY            = "EMPTY_KEY"
	RES_KEY_NOT_EXIST        = "KEY_NOT_EXIST"
	RES_KEY_EXIST_ALREADY    = "KEY_EXIST_ALR"

	RES_NO_ENOUGH_WORKER = "NO_ENOUGH_WORKER"
	RES_NO_PARTITION     = "NO_PARTITION"
	RES_NO_CLUSTER       = "NO_CLUSTER" // if command run on node that not in cluster yet
	RES_NODE_NA          = "NODE_NA"    // if node not accessible

	RES_NO_CTL               = "NO_CONTROLLER" // if no controller leader available
	RES_INVALID_ADDR         = "INVALID_ADDR"
	RES_CONN_RESET           = "CONN_RESET" // Client need to reconnect
	RES_REJECTED             = "REJECTED"   // if node not allow to be removed
	RES_INVALID_PARTITION_ID = "RES_INVALID_PARTITION_ID"
)
