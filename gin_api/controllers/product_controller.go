package controllers

import (
    "app/config"
    "app/models"
    "app/types"
    "app/util"
    "fmt"
    "math"
    "net/http"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/gin-gonic/gin/binding"
)

type ProductController struct {
}

func (con *ProductController) Index(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
    sort := c.DefaultQuery("sort", "Product.id")
    sortDirection := util.Ternary(c.Query("sort") != "", util.Ternary(c.Query("desc") != "", "desc", "asc"), "asc")
    column := c.Query("sc")
    products := []map[string]interface{}{}
    query := config.DB.Table("Product").
        Select("Product.id as Id, Product.image as Image, Product.name as Name, Product.price as Price, Brand.name as BrandName, UserAccount.name as UserAccountName").
        Joins("LEFT JOIN Brand on Product.brand_id = Brand.id").
        Joins("LEFT JOIN UserAccount on Product.create_user = UserAccount.id").
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
    if err := query.Offset((page - 1) * size).Limit(size).Find(&products).Error; err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"products": products, "last": last})
}

func (con *ProductController) GetCreate(c *gin.Context) {
    brands := []map[string]interface{}{}
    config.DB.Table("Brand").
        Select("Brand.id as Id, Brand.name as Name").
        Order("Brand.name" + " " + "asc").
        Find(&brands)
    c.JSON(http.StatusOK, gin.H{ "brands": brands })
}

func (con *ProductController) Create(c *gin.Context) {
    var product models.Product
    if err := c.ShouldBindWith(&product, binding.FormMultipart); err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, util.GetErrors(err))
        return
    }
    image, _ := c.FormFile("ImageFile")
    if image != nil {
        product.Image = util.AddressOf(util.GetFile("products", image))
    }
    product.CreateUser = &util.GetUser(c).Id
    product.CreateDate = util.AddressOf(types.Date(time.Now()))
    if err := config.DB.Create(&product).Error; err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }
    c.JSON(http.StatusOK, product)
}

func (con *ProductController) Get(c *gin.Context) {
    var product map[string]interface{}
    config.DB.Table("Product").
        Select("Product.id as Id, Product.name as Name, Product.price as Price, Brand.name as BrandName, UserAccount.name as UserAccountName, Product.image as Image").
        Joins("LEFT JOIN Brand on Product.brand_id = Brand.id").
        Joins("LEFT JOIN UserAccount on Product.create_user = UserAccount.id").
        Where("Product.id = ?", c.Params.ByName("id")).
        Take(&product)
    c.JSON(http.StatusOK, gin.H{"product": product })
}

func (con *ProductController) Edit(c *gin.Context) {
    var product map[string]interface{}
    config.DB.Table("Product").
        Select("Product.id as Id, Product.name as Name, Product.price as Price, Product.brand_id as BrandId, Product.image as Image").
        Where("Product.id = ?", c.Params.ByName("id")).
        Take(&product)
    brands := []map[string]interface{}{}
    config.DB.Table("Brand").
        Select("Brand.id as Id, Brand.name as Name").
        Order("Brand.name" + " " + "asc").
        Find(&brands)
    c.JSON(http.StatusOK, gin.H{"product": product, "brands": brands })
}

func (con *ProductController) Update(c *gin.Context) {
    var product models.ProductUpdate
    if err := c.ShouldBindWith(&product, binding.FormMultipart); err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, util.GetErrors(err))
        return
    }
    image, _ := c.FormFile("ImageFile")
    if image != nil {
        product.Image = util.AddressOf(util.GetFile("products", image))
    }
    productMap := util.ToMap(product)
    if err := config.DB.Model(&product).Updates(&productMap).Error; err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }
    c.JSON(http.StatusOK, product)
}

func (con *ProductController) GetDelete(c *gin.Context) {
    var product map[string]interface{}
    config.DB.Table("Product").
        Select("Product.id as Id, Product.name as Name, Product.price as Price, Brand.name as BrandName, UserAccount.name as UserAccountName, Product.image as Image").
        Joins("LEFT JOIN Brand on Product.brand_id = Brand.id").
        Joins("LEFT JOIN UserAccount on Product.create_user = UserAccount.id").
        Where("Product.id = ?", c.Params.ByName("id")).
        Take(&product)
    c.JSON(http.StatusOK, gin.H{"product": product })
}

func (con *ProductController) Delete(c *gin.Context) {
    var product models.Product
    err := config.DB.
        Where("Product.id = ?", c.Params.ByName("id")).
        Delete(&product).Error
    if err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }
    c.Status(http.StatusOK)
}
