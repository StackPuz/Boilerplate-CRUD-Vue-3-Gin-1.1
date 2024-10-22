package middleware

import (
    "app/types"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt"
    "github.com/spf13/viper"
)

func Authenticate() gin.HandlerFunc {
    return func(c *gin.Context) {
        allows := []string{
            "/api/login",
            "/api/logout",
            "/api/resetPassword",
            "/api/changePassword",
            "/api/stack",
        }
        for _, path := range allows {
            if c.Request.URL.Path == path || strings.HasPrefix(c.Request.URL.Path, allows[3]) {
                c.Next()
                return
            }
        }
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.Status(http.StatusUnauthorized)
            c.Abort()
            return
        }
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        token, err := jwt.ParseWithClaims(tokenString, &types.Claims{}, func(token *jwt.Token) (interface{}, error) {
            return []byte(viper.GetString("jwtSecret")), nil
        })
        if err != nil || !token.Valid {
            c.Status(http.StatusUnauthorized)
            c.Abort()
            return
        }
        if claims, ok := token.Claims.(*types.Claims); ok {
            c.Set("user", claims)
        } else {
            c.Status(http.StatusUnauthorized)
            c.Abort()
            return
        }
        c.Next()
    }
}