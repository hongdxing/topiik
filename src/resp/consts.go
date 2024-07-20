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
