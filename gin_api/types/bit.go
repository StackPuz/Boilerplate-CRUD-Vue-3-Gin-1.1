package types

import (
    "database/sql/driver"
    "fmt"
    "strconv"
    "strings"
)

type Bit bool

func (v *Bit) UnmarshalJSON(b []byte) error {
    str := string(b)
    if strings.HasPrefix(str, "\"") {
        str, _ = strconv.Unquote(str)
    }
    if str == "true" || str == "1" {
        *v = true
    } else {
        *v = false
    }
    return nil
}

func (b *Bit) Scan(value interface{}) error {
    if value == nil {
        *b = false
    } else {
        switch v := value.(type) {
        case []byte:
            *b = (v[0] == 1)
        case string:
                *b = (v == "1")
        case bool:
            *b = Bit(v)
        default:
            return fmt.Errorf("unsupported type: %T", v)
        }
    }
    return nil
}

func (b Bit) Value() (driver.Value, error) {
    return bool(b), nil
}