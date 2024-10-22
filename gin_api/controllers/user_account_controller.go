package controllers

import (
    "app/config"
    "app/models"
    "app/util"
    "fmt"
    "math"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/gin-gonic/gin/binding"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
)

type UserAccountController struct {
}

func (con *UserAccountController) Index(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
    sort := c.DefaultQuery("sort", "UserAccount.id")
    sortDirection := util.Ternary(c.Query("sort") != "", util.Ternary(c.Query("desc") != "", "desc", "asc"), "asc")
    column := c.Query("sc")
    userAccounts := []map[string]interface{}{}
    query := config.DB.Table("UserAccount").
        Select("UserAccount.id as Id, UserAccount.name as Name, UserAccount.email as Email, UserAccount.active as Active").
        Order(sort + " " + sortDirection)
    if util.IsInvalidSearch(query.Statement.Selects[0], column) {
        c.Status(http.StatusForbidden)
        return
    }
    if c.Query("sw") != "" {
        search := c.Query("sw")
        operator := util.GetOperator(c.Query("so"))
        if operator == "like" {
            search = "%" + search + "%"
        }
        query.Where(fmt.Sprintf("%s %s ?", column, operator), search)
    }
    var count int64
    query.Count(&count)
    last := math.Ceil(float64(count) / float64(size))
    if err := query.Offset((page - 1) * size).Limit(size).Find(&userAccounts).Error; err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"userAccounts": userAccounts, "last": last})
}

func (con *UserAccountController) GetCreate(c *gin.Context) {
    roles := []map[string]interface{}{}
    config.DB.Table("Role").
        Select("Role.id as Id, Role.name as Name").
        Find(&roles)
    c.JSON(http.StatusOK, gin.H{ "roles": roles })
}

func (con *UserAccountController) Create(c *gin.Context) {
    var userAccount models.UserAccount
    if err := c.ShouldBindBodyWith(&userAccount, binding.JSON); err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, util.GetErrors(err))
        return
    }
    token := uuid.New().String()
    userAccount.PasswordResetToken = &token
    password, _ := bcrypt.GenerateFromPassword([]byte(uuid.New().String()[:10]), 10)
    userAccount.Password = string(password)
    if err := config.DB.Create(&userAccount).Error; err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }
    var json map[string]interface{}
    c.ShouldBindBodyWith(&json, binding.JSON)
    roles := json["RoleId"].([]interface{})
    for _, role := range roles {
        userRole := map[string]interface{}{ "user_id": userAccount.Id,"role_id": role }
        config.DB.Table("UserRole").Create(&userRole)
    }
    c.JSON(http.StatusOK, userAccount)
}

func (con *UserAccountController) Get(c *gin.Context) {
    var userAccount map[string]interface{}
    config.DB.Table("UserAccount").
        Select("UserAccount.id as Id, UserAccount.name as Name, UserAccount.email as Email, UserAccount.active as Active").
        Where("UserAccount.id = ?", c.Params.ByName("id")).
        Take(&userAccount)
    userAccountUserRoles := []map[string]interface{}{}
    config.DB.Table("UserAccount").
        Select("Role.name as RoleName").
        Joins("JOIN UserRole on UserAccount.id = UserRole.user_id").
        Joins("JOIN Role on UserRole.role_id = Role.id").
        Where("UserAccount.id = ?", c.Params.ByName("id")).
        Find(&userAccountUserRoles)
    c.JSON(http.StatusOK, gin.H{"userAccount": userAccount, "userAccountUserRoles": userAccountUserRoles })
}

func (con *UserAccountController) Edit(c *gin.Context) {
    var userAccount map[string]interface{}
    config.DB.Table("UserAccount").
        Select("UserAccount.id as Id, UserAccount.name as Name, UserAccount.email as Email, UserAccount.active as Active").
        Where("UserAccount.id = ?", c.Params.ByName("id")).
        Take(&userAccount)
    userAccountUserRoles := []map[string]interface{}{}
    config.DB.Table("UserAccount").
        Select("UserRole.role_id as RoleId").
        Joins("JOIN UserRole on UserAccount.id = UserRole.user_id").
        Where("UserAccount.id = ?", c.Params.ByName("id")).
        Find(&userAccountUserRoles)
    roles := []map[string]interface{}{}
    config.DB.Table("Role").
        Select("Role.id as Id, Role.name as Name").
        Find(&roles)
    c.JSON(http.StatusOK, gin.H{"userAccount": userAccount, "userAccountUserRoles": userAccountUserRoles, "roles": roles })
}

func (con *UserAccountController) Update(c *gin.Context) {
    var userAccount models.UserAccountUpdate
    if err := c.ShouldBindBodyWith(&userAccount, binding.JSON); err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, util.GetErrors(err))
        return
    }
    if userAccount.Password != "" {
        password, _ := bcrypt.GenerateFromPassword([]byte(userAccount.Password), 10)
        userAccount.Password = string(password)
    }
    userAccountMap := util.ToMap(userAccount)
    if userAccountMap["Password"] == "" {
        delete(userAccountMap, "Password")
    }
    if err := config.DB.Model(&userAccount).Updates(&userAccountMap).Error; err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }
    var json map[string]interface{}
    c.ShouldBindBodyWith(&json, binding.JSON)
    roles := json["RoleId"].([]interface{})
    config.DB.Table("UserRole").
        Where("UserRole.user_id = ?", c.Params.ByName("id")).
        Delete(nil)
    for _, role := range roles {
        userRole := map[string]interface{}{ "user_id": c.Params.ByName("id"),"role_id": role }
        config.DB.Table("UserRole").Create(&userRole)
    }
    c.JSON(http.StatusOK, userAccount)
}

func (con *UserAccountController) GetDelete(c *gin.Context) {
    var userAccount map[string]interface{}
    config.DB.Table("UserAccount").
        Select("UserAccount.id as Id, UserAccount.name as Name, UserAccount.email as Email, UserAccount.active as Active").
        Where("UserAccount.id = ?", c.Params.ByName("id")).
        Take(&userAccount)
    userAccountUserRoles := []map[string]interface{}{}
    config.DB.Table("UserAccount").
        Select("Role.name as RoleName").
        Joins("JOIN UserRole on UserAccount.id = UserRole.user_id").
        Joins("JOIN Role on UserRole.role_id = Role.id").
        Where("UserAccount.id = ?", c.Params.ByName("id")).
        Find(&userAccountUserRoles)
    c.JSON(http.StatusOK, gin.H{"userAccount": userAccount, "userAccountUserRoles": userAccountUserRoles })
}

func (con *UserAccountController) Delete(c *gin.Context) {
    var userAccount models.UserAccount
    err := config.DB.
        Where("UserAccount.id = ?", c.Params.ByName("id")).
        Delete(&userAccount).Error
    if err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }
    c.Status(http.StatusOK)
}
