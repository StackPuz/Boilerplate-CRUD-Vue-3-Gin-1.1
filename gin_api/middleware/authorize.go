package middleware

import (
    "app/config"
    "app/util"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
)

func Authorize() gin.HandlerFunc {
    return func(c *gin.Context) {
        path := strings.Split(c.Request.URL.Path, "/")[2]
        var roles string
        for _, m := range config.Menu {
            if m["api"] == path && m["roles"] != nil {
                roles = m["roles"].(string)
                break
            }
        }
        if roles == "" {
            c.Next()
            return
        }
        menuRoles := strings.Split(roles, ",")
        userRoles := util.GetUser(c).Roles
        if util.ArrayContains(menuRoles, userRoles) {
            c.Next()
            return
        }
        c.AbortWithStatus(http.StatusForbidden)
    }
}