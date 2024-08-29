/*
*
*
*
 */

package util

import "time"

func GetUtcEpoch() int64 {
	return time.Now().UTC().Unix()
}
