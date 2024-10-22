package util

import (
    "app/types"
    "errors"
    "fmt"
    "io"
    "mime/multipart"
    "os"
    "path/filepath"
    "strings"
    "time"

    "github.com/fatih/structs"
    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
    "github.com/spf13/viper"
)

var operators = map[string]string{
    "c":  "like",
    "e":  "=",
    "g":  ">",
    "ge": ">=",
    "l":  "<",
    "le": "<=",
}

func GetOperator(oper string) string {
    return operators[oper]
}

func IsInvalidSearch(selectedColumns string, column string) bool {
    if column == "" {
        return false
    }
    for _, s := range strings.Split(selectedColumns, ",") {
        c := strings.Trim(strings.Split(s, " as ")[0], " ")
        if c == column {
            return false
        }
    }
    return true
}

func GetErrors(err error) interface{} {
    var ve validator.ValidationErrors
    errors.As(err, &ve)
    if ve == nil {
        return map[string]string{"message": err.Error()}
    }
    errors := map[string]string{}
    for _, field := range ve {
        errors[field.Field()] = field.Tag()
    }
    return map[string]interface{}{"errors": errors}
}

func GetUser(c *gin.Context) *types.Claims {
    user, _ := c.Get("user")
    return user.(*types.Claims)
}

func GetFile(path string, fileHeader *multipart.FileHeader) string {
    if fileHeader != nil {
        uploadPath := filepath.Join("./uploads", path)
        os.MkdirAll(uploadPath, os.ModePerm)
        ext := filepath.Ext(fileHeader.Filename)
        var filename string
        var filePath string
        for {
            filename = fmt.Sprintf("%x", time.Now().UnixNano()) + ext
            filePath = filepath.Join(uploadPath, filename)
            _, err := os.Stat(filePath)
            if os.IsNotExist(err) {
                break
            }
        }
        file, _ := fileHeader.Open()
        destination, _ := os.Create(filePath)
        io.Copy(destination, file)
        return filename
    }
    return ""
}

func SendMail(mailType, email, token, user string) {
    body := viper.GetString("mail." + mailType)
    body = strings.ReplaceAll(body, "{app_url}", viper.GetString("app.url"))
    body = strings.ReplaceAll(body, "{app_name}", viper.GetString("app.name"))
    body = strings.ReplaceAll(body, "{token}", token)
    if user != "" {
        body = strings.ReplaceAll(body, "{user}", user)
    }
    subject := Ternary(mailType == "welcome", "Login Information", Ternary(mailType == "reset", "Reset Password", viper.GetString("app.name")+" message"))
    body = fmt.Sprintf("From: %s\nSubject: %s\n\n%s", viper.GetString("mail.sender"), subject, body)
    /* You need to complete the SMTP Server configuration before you can sent mail
    auth := smtp.PlainAuth("", viper.GetString("smtp.user"), viper.GetString("smtp.password"), viper.GetString("smtp.host"))
    smtp.SendMail(viper.GetString("smtp.host")+":"+viper.GetString("smtp.port"), auth, viper.GetString("mail.sender"), []string{email}, []byte(body))
    */
}

// GORM will only update non-zero fields; to update zero fields, we need to use a map.
func ToMap(input any) map[string]interface{} {
    return structs.Map(input)
}

func AddressOf[T any](value T) *T {
    return &value
}

func ArrayContains(targets []string, sources []string) bool {
    for _, target := range targets {
        for _, source := range sources {
            if target == source {
                return true
            }
        }
    }
    return false
}

func Ternary[T any](condition bool, trueVal, falseVal T) T {
    if condition {
        return trueVal
    }
    return falseVal
}