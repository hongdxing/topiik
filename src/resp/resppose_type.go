/*
* @author: Duan Hongxing
* @date: 21 Aug, 2024
* @desc:
*
 */

package resp

type RespType int8

const (
	String       RespType = 1
	StringArray  RespType = 2
	Integer      RespType = 3 // int64
	IntegerArray RespType = 4 // int64 list
	Double       RespType = 5
	DoubleArray  RespType = 6
	Map          RespType = 7
	Set          RespType = 8
	//Byte         RespType = 9
	//ByteArray    RespType = 10

	Redirect RespType = 32
)

type RespFlag int8

const (
	Success RespFlag = 1
	Fail    RespFlag = 0
)
