package controllers

import (
    "app/config"
    "app/models"
    "app/util"
    "net/http"

    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
)

type SystemController struct {
}

func (con *SystemController) Profile(c *gin.Context) {
    var userAccount map[string]interface{}
    config.DB.Table("UserAccount").
        Select("UserAccount.name as Name, UserAccount.email as Email").
        Where("UserAccount.id = ?", util.GetUser(c).Id).
        Find(&userAccount)
    c.JSON(http.StatusOK, gin.H{"userAccount": userAccount})
}

func (con *SystemController) UpdateProfile(c *gin.Context) {
    var userAccount models.UserAccount
    c.BindJSON(&userAccount)
    userAccount.Id = util.GetUser(c).Id
    if userAccount.Password != "" {
        password, _ := bcrypt.GenerateFromPassword([]byte(userAccount.Password), 10)
        userAccount.Password = string(password)
    }
    config.DB.Updates(&userAccount)
    c.Status(http.StatusOK)
}

func (con *SystemController) Stack(c *gin.Context) {
    c.String(http.StatusOK, "Vue 3 + Go API 1 + MySQL")
}