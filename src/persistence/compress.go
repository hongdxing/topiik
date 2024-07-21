/***
** author: duan hongxing
** date: 29 Jun 2024
** desc:
**	To compress db files
**
**/

package persistence


/**
**
** 
**
**
**
**/

/**
** Algorithm:
**	1) maintain a list that record event DEL and command position,
	when event count reach a threshhold like 1000, then trigger 
	compress to remove all operations of relvant keys
**	2) increamennt speed trigger, when db size increased xx(512) MB in xx(15) min,
	then trigger compress, by travelling the db file
	i) if has ttl and epxired already, remove the operation
**/

func compress() {

}
