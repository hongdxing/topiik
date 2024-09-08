//author: Duan HongXing
//date: 25 Aug, 2024

package internal

import "encoding/json"

type Abytes [][]byte

type Req struct {
	Keys Abytes
	Vals Abytes
	Args string
}

func (r *Req) Build() *Req {

	return r
}

func (r *Req) WithKeys(keys Abytes) *Req {
	r.Keys = keys
	return r
}

func (r *Req) WithKey(key string) *Req {
	r.Keys = Abytes{[]byte(key)}
	return r
}

func (r *Req) WithVals(vals Abytes) *Req {
	r.Vals = vals
	return r
}

func (r *Req) WithVal(val string) *Req {
	r.Vals = Abytes{[]byte(val)}
	return r
}

func (r *Req) WithArgs(args string) *Req {
	r.Args = args
	return r
}

func (r *Req) Marshal() ([]byte, error) {
	buf, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
