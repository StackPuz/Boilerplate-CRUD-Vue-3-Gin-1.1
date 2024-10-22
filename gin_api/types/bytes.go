package types

import (
    "encoding/json"
    "strconv"
    "strings"
)

type Bytes []byte

func (v Bytes) MarshalJSON() ([]byte, error) {
    return json.Marshal(strings.TrimRight(string(v), "\x00"))
}

func (v *Bytes) UnmarshalJSON(b []byte) error {
    str, _ := strconv.Unquote(string(b))
    *v = Bytes(str)
    return nil
}
