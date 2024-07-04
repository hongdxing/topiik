/***
** author: duan hongxing
** data: 4 Jul 2024
** desc:
**
**/

package ccss

import "fmt"

func JoinACK(id string, salorAddr string) (result string, err error) {
	salor := Salor{
		Id:      id,
		Address: salorAddr,
	}
	salors = append(salors, salor)
	fmt.Println(salors)
	return "OK", nil
}
