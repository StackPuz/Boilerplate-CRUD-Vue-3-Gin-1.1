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
)

type BrandController struct {
}

func (con *BrandController) Index(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
    sort := c.DefaultQuery("sort", "Brand.id")
    sortDirection := util.Ternary(c.Query("sort") != "", util.Ternary(c.Query("desc") != "", "desc", "asc"), "asc")
    column := c.Query("sc")
    brands := []map[string]interface{}{}
    query := config.DB.Table("Brand").
        Select("Brand.id as Id, Brand.name as Name").
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
    if err := query.Offset((page - 1) * size).Limit(size).Find(&brands).Error; err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"brands": brands, "last": last})
}

func (con *BrandController) GetCreate(c *gin.Context) {
    c.Status(http.StatusOK)
}

func (con *BrandController) Create(c *gin.Context) {
    var brand models.Brand
    if err := c.ShouldBind(&brand); err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, util.GetErrors(err))
        return
    }
    if err := config.DB.Create(&brand).Error; err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }
    c.JSON(http.StatusOK, brand)
}

func (con *BrandController) Get(c *gin.Context) {
    var brand map[string]interface{}
    config.DB.Table("Brand").
        Select("Brand.id as Id, Brand.name as Name").
        Where("Brand.id = ?", c.Params.ByName("id")).
        Take(&brand)
    brandProducts := []map[string]interface{}{}
    config.DB.Table("Brand").
        Select("Product.name as Name, Product.price as Price").
        Joins("JOIN Product on Brand.id = Product.brand_id").
        Where("Brand.id = ?", c.Params.ByName("id")).
        Find(&brandProducts)
    c.JSON(http.StatusOK, gin.H{"brand": brand, "brandProducts": brandProducts })
}

func (con *BrandController) Edit(c *gin.Context) {
    var brand map[string]interface{}
    config.DB.Table("Brand").
        Select("Brand.id as Id, Brand.name as Name").
        Where("Brand.id = ?", c.Params.ByName("id")).
        Take(&brand)
    brandProducts := []map[string]interface{}{}
    config.DB.Table("Brand").
        Select("Product.name as Name, Product.price as Price, Product.id as Id").
        Joins("JOIN Product on Brand.id = Product.brand_id").
        Where("Brand.id = ?", c.Params.ByName("id")).
        Find(&brandProducts)
    c.JSON(http.StatusOK, gin.H{"brand": brand, "brandProducts": brandProducts })
}

func (con *BrandController) Update(c *gin.Context) {
    var brand models.BrandUpdate
    if err := c.ShouldBind(&brand); err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, util.GetErrors(err))
        return
    }
    brandMap := util.ToMap(brand)
    if err := config.DB.Model(&brand).Updates(&brandMap).Error; err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }
    c.JSON(http.StatusOK, brand)
}

func (con *BrandController) GetDelete(c *gin.Context) {
    var brand map[string]interface{}
    config.DB.Table("Brand").
        Select("Brand.id as Id, Brand.name as Name").
        Where("Brand.id = ?", c.Params.ByName("id")).
        Take(&brand)
    brandProducts := []map[string]interface{}{}
    config.DB.Table("Brand").
        Select("Product.name as Name, Product.price as Price").
        Joins("JOIN Product on Brand.id = Product.brand_id").
        Where("Brand.id = ?", c.Params.ByName("id")).
        Find(&brandProducts)
    c.JSON(http.StatusOK, gin.H{"brand": brand, "brandProducts": brandProducts })
}

func (con *BrandController) Delete(c *gin.Context) {
    var brand models.Brand
    err := config.DB.
        Where("Brand.id = ?", c.Params.ByName("id")).
        Delete(&brand).Error
    if err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }
    c.Status(http.StatusOK)
}
