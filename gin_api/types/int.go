package types

import (
    "database/sql/driver"
    "encoding/json"
    "fmt"
    "strconv"
)

type Int32 int32
type Int64 int64

func (i Int32) MarshalJSON() ([]byte, error) {
    return json.Marshal(int32(i))
}

func (i Int64) MarshalJSON() ([]byte, error) {
    return json.Marshal(int64(i))
}

func unmarshalInt[T ~int32 | ~int64](b []byte, i *T) error {
    if len(b) > 0 && b[0] == '"' {
        var str string
        json.Unmarshal(b, &str)
        value, _ := strconv.Atoi(str)
        *i = T(value)
    } else {
        var value int64
        json.Unmarshal(b, &value)
        *i = T(value)
    }
    return nil
}

func (i *Int32) UnmarshalJSON(b []byte) error {
    return unmarshalInt(b, i)
}

func (i *Int64) UnmarshalJSON(b []byte) error {
    return unmarshalInt(b, i)
}

func (i *Int32) Scan(value interface{}) error {
    switch v := value.(type) {
    case string:
        vi, _ := strconv.ParseInt(v, 10, 32)
        *i = Int32(vi)
    case int64:
        *i = Int32(v)
    default:
        return fmt.Errorf("unsupported type: %T", v)
    }
    return nil
}

func (i *Int64) Scan(value interface{}) error {
    switch v := value.(type) {
    case string:
        vi, _ := strconv.ParseInt(v, 10, 64)
        *i = Int64(vi)
    case int64:
        *i = Int64(v)
    default:
        return fmt.Errorf("unsupported type: %T", v)
    }
    return nil
}

func (i Int32) Value() (driver.Value, error) {
    return int64(i), nil
}

func (i Int64) Value() (driver.Value, error) {
    return int64(i), nil
}