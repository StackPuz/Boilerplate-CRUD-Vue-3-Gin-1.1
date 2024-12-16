package controllers

import (
    "app/config"
    "app/models"
    "app/types"
    "app/util"
    "fmt"
    "net/http"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt"
    "github.com/google/uuid"
    "github.com/spf13/viper"
    "golang.org/x/crypto/bcrypt"
)

type LoginController struct {
}

func getRoles(userId int) []string {
    var userRoles []map[string]interface{}
    config.DB.Table("UserRole").
        Select(`Role.name as "name"`).
        Joins("join Role on Role.id = UserRole.role_id").
        Where("UserRole.user_id = ?", userId).
        Find(&userRoles)
    var roles []string
    for _, userRole := range userRoles {
        roles = append(roles, userRole["name"].(string))
    }
    return roles
}

func getMenu(roles []string) []map[string]interface{} {
    menu := make([]map[string]interface{}, 0)
    for _, e := range config.Menu {
        show := e["show"].(bool)
        rolesStr, exist := e["roles"].(string)
        if show && (!exist || util.ArrayContains(strings.Split(rolesStr, ","), roles)) {
            item := map[string]interface{}{
                "title": e["title"],
                "path":  e["path"],
            }
            menu = append(menu, item)
        }
    }
    return menu
}

func (con *LoginController) GetUser(c *gin.Context) {
    user := util.GetUser(c)
    roles := getRoles(int(user.Id))
    c.JSON(http.StatusOK, gin.H{"name": user.Name, "menu": getMenu(roles)})
}

func (con *LoginController) Login(c *gin.Context) {
    var payload map[string]string
    var user models.UserAccount
    c.BindJSON(&payload)
    if err := config.DB.Where("name = ?", payload["name"]).First(&user).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid credentials"})
        return
    }
    if !user.Active {
        c.JSON(http.StatusBadRequest, gin.H{"message": "User is disabled"})
    } else {
        if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload["password"])) == nil {
            roles := getRoles(int(user.Id))
            token := jwt.NewWithClaims(jwt.SigningMethodHS256, &types.Claims{
                Id:    user.Id,
                Name:  user.Name,
                Roles: roles,
                StandardClaims: jwt.StandardClaims{
                    IssuedAt:  time.Now().Unix(),
                    ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
                },
            })
            tokenString, err := token.SignedString([]byte(viper.GetString("jwtSecret")))
            if err != nil {
                fmt.Println(err.Error())
            }
            c.JSON(http.StatusOK, gin.H{
                "token": tokenString,
                "user": gin.H{
                    "name": user.Name,
                    "menu": getMenu(roles),
                },
            })
        } else {
            c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid credentials"})
        }
    }
}

func (con *LoginController) Logout(c *gin.Context) {
    c.Status(http.StatusOK)
}

func (con *LoginController) ResetPassword(c *gin.Context) {
    var user models.UserAccount
    c.ShouldBind(&user)
    if err := config.DB.Where("email = ?", user.Email).First(&user).Error; err != nil {
        c.Status(http.StatusNotFound)
        return
    }
    token := uuid.New().String()
    config.DB.Model(&user).Update("password_reset_token", token)
    util.SendMail("reset", user.Email, token, "")
    c.Status(http.StatusOK)
}

func (con *LoginController) GetChangePassword(c *gin.Context) {
    var user models.UserAccount
    if err := config.DB.Where("password_reset_token = ?", c.Params.ByName("token")).First(&user).Error; err != nil {
        c.Status(http.StatusNotFound)
        return
    }
    c.Status(http.StatusOK)
}

func (con *LoginController) ChangePassword(c *gin.Context) {
	var payload map[string]string
    c.BindJSON(&payload)
    var user models.UserAccount
    if err := config.DB.Where("password_reset_token = ?", c.Params.ByName("token")).First(&user).Error; err != nil {
        c.Status(http.StatusNotFound)
        return
    }
    password, _ := bcrypt.GenerateFromPassword([]byte(payload["password"]), 10)
    config.DB.Model(&user).Updates(map[string]interface{}{"password": string(password), "password_reset_token": nil})
    c.Status(http.StatusOK)
}