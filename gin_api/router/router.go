package router

import (
    "app/controllers"
    "app/middleware"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
    loginControlller := controllers.LoginController{}
    systemControlller := controllers.SystemController{}
    userAccountController := controllers.UserAccountController{}
    productController := controllers.ProductController{}
    brandController := controllers.BrandController{}
    orderHeaderController := controllers.OrderHeaderController{}
    orderDetailController := controllers.OrderDetailController{}
    corsConfig := cors.DefaultConfig()
    corsConfig.AllowAllOrigins = true
    corsConfig.AddAllowHeaders("Authorization")
    router := gin.Default()
    router.Static("/uploads", "./uploads")
    router.Use(cors.New(corsConfig)).Use(middleware.Authenticate()).Use(middleware.Authorize())
    router.POST("/api/login", loginControlller.Login)
    router.GET("/api/logout", loginControlller.Logout)
    router.POST("/api/resetPassword", loginControlller.ResetPassword)
    router.GET("/api/changePassword/:token", loginControlller.GetChangePassword)
    router.POST("/api/changePassword/:token", loginControlller.ChangePassword)
    router.GET("/api/user", loginControlller.GetUser)
    router.GET("/api/profile", systemControlller.Profile)
    router.POST("/api/updateProfile", systemControlller.UpdateProfile)
    router.GET("/api/stack", systemControlller.Stack)
    //Not utilized the route.Group() because the issue https://github.com/gin-gonic/gin/issues/531

    router.GET("/api/userAccounts", userAccountController.Index)
    router.POST("/api/userAccounts", userAccountController.Create)
    router.GET("/api/userAccounts/create", userAccountController.GetCreate)
    router.GET("/api/userAccounts/:id", userAccountController.Get)
    router.GET("/api/userAccounts/:id/edit", userAccountController.Edit)
    router.PUT("/api/userAccounts/:id", userAccountController.Update)
    router.GET("/api/userAccounts/:id/delete", userAccountController.GetDelete)
    router.DELETE("/api/userAccounts/:id", userAccountController.Delete)

    router.GET("/api/products", productController.Index)
    router.POST("/api/products", productController.Create)
    router.GET("/api/products/create", productController.GetCreate)
    router.GET("/api/products/:id", productController.Get)
    router.GET("/api/products/:id/edit", productController.Edit)
    router.PUT("/api/products/:id", productController.Update)
    router.GET("/api/products/:id/delete", productController.GetDelete)
    router.DELETE("/api/products/:id", productController.Delete)

    router.GET("/api/brands", brandController.Index)
    router.POST("/api/brands", brandController.Create)
    router.GET("/api/brands/create", brandController.GetCreate)
    router.GET("/api/brands/:id", brandController.Get)
    router.GET("/api/brands/:id/edit", brandController.Edit)
    router.PUT("/api/brands/:id", brandController.Update)
    router.GET("/api/brands/:id/delete", brandController.GetDelete)
    router.DELETE("/api/brands/:id", brandController.Delete)

    router.GET("/api/orderHeaders", orderHeaderController.Index)
    router.POST("/api/orderHeaders", orderHeaderController.Create)
    router.GET("/api/orderHeaders/create", orderHeaderController.GetCreate)
    router.GET("/api/orderHeaders/:id", orderHeaderController.Get)
    router.GET("/api/orderHeaders/:id/edit", orderHeaderController.Edit)
    router.PUT("/api/orderHeaders/:id", orderHeaderController.Update)
    router.GET("/api/orderHeaders/:id/delete", orderHeaderController.GetDelete)
    router.DELETE("/api/orderHeaders/:id", orderHeaderController.Delete)

    router.POST("/api/orderDetails", orderDetailController.Create)
    router.GET("/api/orderDetails/create", orderDetailController.GetCreate)
    router.GET("/api/orderDetails/:orderId/:no/edit", orderDetailController.Edit)
    router.PUT("/api/orderDetails/:orderId/:no", orderDetailController.Update)
    router.GET("/api/orderDetails/:orderId/:no/delete", orderDetailController.GetDelete)
    router.DELETE("/api/orderDetails/:orderId/:no", orderDetailController.Delete)

    return router
}