/***
** author: duan hongxing
** date: 23 Jun 2024
** desc:
**
**/

package executor

import (
	"encoding/json"
	"topiik/internal/datatype"
)

/***String response***/
func marshalResponseError(err error) []byte {
	return marshalResponse(err.Error(), false)
}

func marshalResponseSuccess(response string) []byte {
	return marshalResponse(response, true)
}

func marshalResponse(response string, success bool) []byte {
	b, _ := json.Marshal(&datatype.StrResponse{R: success, M: response})
	return b
}

/***Integer response***/
/*
func marshalIntegerResponseError(err error) []byte {
	return marshalIntegerResponse(-1, false)
}*/

func marshalIntegerResponseSuccess(response int) []byte {
	return marshalIntegerResponse(response, true)
}

func marshalIntegerResponse(response int, success bool) []byte {
	b, _ := json.Marshal(&datatype.IntegerResponse{R: success, M: response})
	return b
}

/***String response***/
/*
func marshalListResponseError(err error) []byte {
	return marshalListResponse(err.Error(), false)
}*/

func marshalListResponseSuccess(response []string) []byte {
	return marshalListResponse(response, true)
}

func marshalListResponse(response []string, success bool) []byte {
	b, _ := json.Marshal(&datatype.ListResponse{R: success, M: response})
	return b
}
