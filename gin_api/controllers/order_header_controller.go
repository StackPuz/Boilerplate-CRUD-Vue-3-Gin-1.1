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

    "github.com/gin-gonic/gin"
)

type OrderHeaderController struct {
}

func (con *OrderHeaderController) Index(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
    sort := c.DefaultQuery("sort", "OrderHeader.id")
    sortDirection := util.Ternary(c.Query("sort") != "", util.Ternary(c.Query("desc") != "", "desc", "asc"), "asc")
    column := c.Query("sc")
    orderHeaders := []map[string]interface{}{}
    query := config.DB.Table("OrderHeader").
        Select("OrderHeader.id as Id, Customer.name as CustomerName, OrderHeader.order_date as OrderDate").
        Joins("LEFT JOIN Customer on OrderHeader.customer_id = Customer.id").
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
        if column == "OrderHeader.order_date" {
            search = types.FormatDateStr((search))
        }
        query.Where(fmt.Sprintf("%s %s ?", column, operator), search)
    }
    var count int64
    query.Count(&count)
    last := math.Ceil(float64(count) / float64(size))
    if err := query.Offset((page - 1) * size).Limit(size).Find(&orderHeaders).Error; err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"orderHeaders": orderHeaders, "last": last})
}

func (con *OrderHeaderController) GetCreate(c *gin.Context) {
    customers := []map[string]interface{}{}
    config.DB.Table("Customer").
        Select("Customer.id as Id, Customer.name as Name").
        Order("Customer.name" + " " + "asc").
        Find(&customers)
    c.JSON(http.StatusOK, gin.H{ "customers": customers })
}

func (con *OrderHeaderController) Create(c *gin.Context) {
    var orderHeader models.OrderHeader
    if err := c.ShouldBind(&orderHeader); err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, util.GetErrors(err))
        return
    }
    if err := config.DB.Create(&orderHeader).Error; err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }
    c.JSON(http.StatusOK, orderHeader)
}

func (con *OrderHeaderController) Get(c *gin.Context) {
    var orderHeader map[string]interface{}
    config.DB.Table("OrderHeader").
        Select("OrderHeader.id as Id, Customer.name as CustomerName, OrderHeader.order_date as OrderDate").
        Joins("LEFT JOIN Customer on OrderHeader.customer_id = Customer.id").
        Where("OrderHeader.id = ?", c.Params.ByName("id")).
        Take(&orderHeader)
    orderHeaderOrderDetails := []map[string]interface{}{}
    config.DB.Table("OrderHeader").
        Select("OrderDetail.no as No, Product.name as ProductName, OrderDetail.qty as Qty").
        Joins("JOIN OrderDetail on OrderHeader.id = OrderDetail.order_id").
        Joins("JOIN Product on OrderDetail.product_id = Product.id").
        Where("OrderHeader.id = ?", c.Params.ByName("id")).
        Find(&orderHeaderOrderDetails)
    c.JSON(http.StatusOK, gin.H{"orderHeader": orderHeader, "orderHeaderOrderDetails": orderHeaderOrderDetails })
}

func (con *OrderHeaderController) Edit(c *gin.Context) {
    var orderHeader map[string]interface{}
    config.DB.Table("OrderHeader").
        Select("OrderHeader.id as Id, OrderHeader.customer_id as CustomerId, OrderHeader.order_date as OrderDate").
        Where("OrderHeader.id = ?", c.Params.ByName("id")).
        Take(&orderHeader)
    orderHeaderOrderDetails := []map[string]interface{}{}
    config.DB.Table("OrderHeader").
        Select("OrderDetail.no as No, Product.name as ProductName, OrderDetail.qty as Qty, OrderDetail.order_id as OrderId").
        Joins("JOIN OrderDetail on OrderHeader.id = OrderDetail.order_id").
        Joins("JOIN Product on OrderDetail.product_id = Product.id").
        Where("OrderHeader.id = ?", c.Params.ByName("id")).
        Find(&orderHeaderOrderDetails)
    customers := []map[string]interface{}{}
    config.DB.Table("Customer").
        Select("Customer.id as Id, Customer.name as Name").
        Order("Customer.name" + " " + "asc").
        Find(&customers)
    c.JSON(http.StatusOK, gin.H{"orderHeader": orderHeader, "orderHeaderOrderDetails": orderHeaderOrderDetails, "customers": customers })
}

func (con *OrderHeaderController) Update(c *gin.Context) {
    var orderHeader models.OrderHeaderUpdate
    if err := c.ShouldBind(&orderHeader); err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, util.GetErrors(err))
        return
    }
    orderHeaderMap := util.ToMap(orderHeader)
    if err := config.DB.Model(&orderHeader).Updates(&orderHeaderMap).Error; err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }
    c.JSON(http.StatusOK, orderHeader)
}

func (con *OrderHeaderController) GetDelete(c *gin.Context) {
    var orderHeader map[string]interface{}
    config.DB.Table("OrderHeader").
        Select("OrderHeader.id as Id, Customer.name as CustomerName, OrderHeader.order_date as OrderDate").
        Joins("LEFT JOIN Customer on OrderHeader.customer_id = Customer.id").
        Where("OrderHeader.id = ?", c.Params.ByName("id")).
        Take(&orderHeader)
    orderHeaderOrderDetails := []map[string]interface{}{}
    config.DB.Table("OrderHeader").
        Select("OrderDetail.no as No, Product.name as ProductName, OrderDetail.qty as Qty").
        Joins("JOIN OrderDetail on OrderHeader.id = OrderDetail.order_id").
        Joins("JOIN Product on OrderDetail.product_id = Product.id").
        Where("OrderHeader.id = ?", c.Params.ByName("id")).
        Find(&orderHeaderOrderDetails)
    c.JSON(http.StatusOK, gin.H{"orderHeader": orderHeader, "orderHeaderOrderDetails": orderHeaderOrderDetails })
}

func (con *OrderHeaderController) Delete(c *gin.Context) {
    var orderHeader models.OrderHeader
    err := config.DB.
        Where("OrderHeader.id = ?", c.Params.ByName("id")).
        Delete(&orderHeader).Error
    if err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }
    c.Status(http.StatusOK)
}
