package types

import (
    "database/sql/driver"
    "encoding/json"
    "fmt"
    "strconv"
)

type Float32 float32
type Float64 float64

func (f Float32) MarshalJSON() ([]byte, error) {
    return json.Marshal(float32(f))
}

func (f Float64) MarshalJSON() ([]byte, error) {
    return json.Marshal(float64(f))
}

func unmarshalFloat[T ~float32 | ~float64](b []byte, f *T) error {
    if len(b) > 0 && b[0] == '"' {
        var str string
        json.Unmarshal(b, &str)
        value, _ := strconv.ParseFloat(str, 64)
        *f = T(value)
    } else {
        var value float64
        json.Unmarshal(b, &value)
        *f = T(value)
    }
    return nil
}

func (f *Float32) UnmarshalJSON(b []byte) error {
    return unmarshalFloat(b, f)
}

func (f *Float64) UnmarshalJSON(b []byte) error {
    return unmarshalFloat(b, f)
}

func scan[T ~float32 | ~float64](value interface{}, f *T) error {
    switch v := value.(type) {
    case []byte:
        float, _ := strconv.ParseFloat(string(v), 64)
        *f = T(float)
    case string:
        float, _ := strconv.ParseFloat(v, 64)
        *f = T(float)
    case float64:
        *f = T(v)
    case float32:
        *f = T(v)
    default:
        return fmt.Errorf("unsupported type: %T", v)
    }
    return nil
}

func (f *Float32) Scan(value interface{}) error {
    return scan(value, f)
}

func (f *Float64) Scan(value interface{}) error {
    return scan(value, f)
}

func (f Float32) Value() (driver.Value, error) {
    return float64(f), nil
}

func (f Float64) Value() (driver.Value, error) {
    return float64(f), nil
}