package controllers

import (
    "app/config"
    "app/models"
    "app/util"
    "net/http"

    "github.com/gin-gonic/gin"
)

type OrderDetailController struct {
}

func (con *OrderDetailController) GetCreate(c *gin.Context) {
    products := []map[string]interface{}{}
    config.DB.Table("Product").
        Select("Product.id as Id, Product.name as Name").
        Order("Product.name" + " " + "asc").
        Find(&products)
    c.JSON(http.StatusOK, gin.H{ "products": products })
}

func (con *OrderDetailController) Create(c *gin.Context) {
    var orderDetail models.OrderDetail
    if err := c.ShouldBind(&orderDetail); err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, util.GetErrors(err))
        return
    }
    if err := config.DB.Create(&orderDetail).Error; err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }
    c.JSON(http.StatusOK, orderDetail)
}

func (con *OrderDetailController) Edit(c *gin.Context) {
    var orderDetail map[string]interface{}
    config.DB.Table("OrderDetail").
        Select("OrderDetail.order_id as OrderId, OrderDetail.no as No, OrderDetail.product_id as ProductId, OrderDetail.qty as Qty").
        Where("OrderDetail.order_id = ?", c.Params.ByName("orderId")).
        Where("OrderDetail.no = ?", c.Params.ByName("no")).
        Take(&orderDetail)
    products := []map[string]interface{}{}
    config.DB.Table("Product").
        Select("Product.id as Id, Product.name as Name").
        Order("Product.name" + " " + "asc").
        Find(&products)
    c.JSON(http.StatusOK, gin.H{"orderDetail": orderDetail, "products": products })
}

func (con *OrderDetailController) Update(c *gin.Context) {
    var orderDetail models.OrderDetailUpdate
    if err := c.ShouldBind(&orderDetail); err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, util.GetErrors(err))
        return
    }
    orderDetailMap := util.ToMap(orderDetail)
    if err := config.DB.Model(&orderDetail).Updates(&orderDetailMap).Error; err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }
    c.JSON(http.StatusOK, orderDetail)
}

func (con *OrderDetailController) GetDelete(c *gin.Context) {
    var orderDetail map[string]interface{}
    config.DB.Table("OrderDetail").
        Select("OrderDetail.order_id as OrderId, OrderDetail.no as No, Product.name as ProductName, OrderDetail.qty as Qty").
        Joins("LEFT JOIN Product on OrderDetail.product_id = Product.id").
        Where("OrderDetail.order_id = ?", c.Params.ByName("orderId")).
        Where("OrderDetail.no = ?", c.Params.ByName("no")).
        Take(&orderDetail)
    c.JSON(http.StatusOK, gin.H{"orderDetail": orderDetail })
}

func (con *OrderDetailController) Delete(c *gin.Context) {
    var orderDetail models.OrderDetail
    err := config.DB.
        Where("OrderDetail.order_id = ?", c.Params.ByName("orderId")).
        Where("OrderDetail.no = ?", c.Params.ByName("no")).
        Delete(&orderDetail).Error
    if err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }
    c.Status(http.StatusOK)
}
