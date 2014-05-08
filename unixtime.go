package stripe

import (
	"errors"
	"strconv"
	"time"
)

var UnixTimeUnmarshalError = errors.New("stripe: invalid timestamp")

type UnixTime struct{ time.Time }

func (t UnixTime) MarshalJSON() ([]byte, error) {
	return strconv.AppendInt(nil, t.UnixNano()/int64(time.Millisecond), 10), nil
}

func (t *UnixTime) UnmarshalJSON(data []byte) error {
	str := string(data)
	var i int64
	if str != "null" {
		var err error
		i, err = strconv.ParseInt(str, 10, 64)
		if err != nil {
			return UnixTimeUnmarshalError
		}
	}
	t.Time = time.Unix(i, 0)
	return nil
}

func (t *UnixTime) Scan(src interface{}) error {
	ts, ok := src.(time.Time)
	if !ok {
		return UnixTimeUnmarshalError
	}
	t.Time = ts
	return nil
}
