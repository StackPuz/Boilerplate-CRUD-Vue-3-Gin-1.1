package types

import (
    "database/sql/driver"
    "fmt"
    "strconv"
    "time"
)

type Date time.Time

var formats = map[string]string{
    "date":     "01/02/2006",
    "time":     "15:04:05",
    "datetime": "01/02/2006 15:04:05",
}

func (d Date) MarshalJSON() ([]byte, error) {
    return []byte(strconv.Quote(FormatDate(time.Time(d)))), nil
}

func (d *Date) UnmarshalJSON(b []byte) error {
    str, _ := strconv.Unquote(string(b))
    *d = Date(GetDate(str))
    return nil
}

func (d *Date) Scan(value interface{}) error {
    switch v := value.(type) {
    case []byte:
        time, _ := time.Parse("15:04:05", string(v))
        *d = Date(time)
    case string:
        time, _ := time.Parse("15:04:05", v)
        *d = Date(time)
    case time.Time:
        *d = Date(v)
    default:
        return fmt.Errorf("unsupported type: %T", v)
    }
    return nil
}

func (d Date) Value() (driver.Value, error) {
    value := time.Time(d)
    if value.IsZero() {
        return nil, nil
    }
    return value, nil
}

func GetDate(value string) time.Time {
    if value == "" {
        return time.Time{}
    }
    if len(value) == len(formats["date"]) {
        date, _ := time.Parse(formats["date"], value)
        return date
    } else if len(value) == len(formats["time"]) {
        date, _ := time.Parse(formats["time"], value)
        return date.AddDate(1, 0, 0)
    } else {
        date, _ := time.Parse(formats["datetime"], value)
        return date
    }
}

func FormatDate(t time.Time) string {
    if t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0 {
        return t.Format(formats["date"])
    } else if t.Year() == 0 && t.Month() == 1 && t.Day() == 1 {
        return t.Format(formats["time"])
    } else {
        return t.Format(formats["datetime"])
    }
}

func FormatDateStr(value string) string {
    if len(value) == len(formats["date"]) {
        return GetDate(value).Format("2006-01-02")
    } else if len(value) == len(formats["time"]) {
        return GetDate(value).Format("15:04:05")
    } else {
        return GetDate(value).Format("2006-01-02 15:04:05")
    }
}