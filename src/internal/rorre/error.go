// author: Duan Hongxing
// date: 16 Sep, 2024

package rorre

type SoketError struct {
}

func (e *SoketError) Error() string {
	return "socket error"
}
