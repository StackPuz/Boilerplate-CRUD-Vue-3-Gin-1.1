package config

import (
    "app/types"
    "fmt"
    "regexp"
    "strings"
    "time"

    "github.com/spf13/viper"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/schema"
)

var DB *gorm.DB
var Menu []map[string]interface{}
var timeExp = regexp.MustCompile(`^([01][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$`)

func Setup() {
    viper.SetConfigFile("config/config.json")
    viper.ReadInConfig()
    viper.UnmarshalKey("menu", &Menu)
    connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true", viper.GetString("db.user"), viper.GetString("db.password"), viper.GetString("db.host"), viper.GetString("db.port"), viper.GetString("db.database"))
    db, _ := gorm.Open(mysql.Open(connection), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
    db.Callback().Query().After("gorm:query").Register("my_plugin:after_query", afterQuery)
    DB = db
}

func afterQuery(db *gorm.DB) {
    formatData(db.Statement.Dest)
}

func formatData(data any) {
    switch d := data.(type) {
    case *[]map[string]interface{}:
        for _, e := range *d {
            formatData(e)
        }
    case *map[string]interface{}:
        formatData(*d)
    case map[string]interface{}:
        for key, value := range d {
            if value != nil {
                switch v := value.(type) {
                case string:
                    if v == "\x00" {
                        d[key] = false
                    } else if v == "\x01" {
                        d[key] = true
                    } else if timeExp.MatchString(v) {
                        time, _ := time.Parse("15:04:05", v)
                        d[key] = types.FormatDate(time)
                    } else if strings.HasSuffix(v, "\x00") {
                        d[key] = strings.TrimRight(v, "\x00")
                    }
                case time.Time:
                    d[key] = types.FormatDate(v)
                case []byte:
                    d[key] = strings.TrimRight(string(v), "\x00")
                }
            }
        }
    }
}