package router

import (
	"net/http"
	"os"

	"github.com/EmeraldLS/phsps-api/controller"
	"github.com/EmeraldLS/phsps-api/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Run() {

	port := os.Getenv("PORT")
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type"},
		AllowMethods:     []string{"POST", "PUT", "GET", "DELETE"},
	}))
	api := r.Group("apiv1")
	{
		api.POST("/register", controller.Register)
		api.POST("/login", controller.Login)
		secured := api.Group("/secured")
		{
			secured.Use(middleware.Auth)
			secured.GET("/hi", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"response": "success",
					"message":  "Hello world",
				})
			})
			secured.POST("/customer", controller.InsertCustomer)
			secured.GET("/customer/:customer_code", controller.GetCustomerByID)
			secured.GET("/customer/search_customer", controller.SearchCustomerByBusinessName)
			secured.GET("/customer", controller.GetAllCustomter)
			secured.PUT("/customer/:customer_code/:job_status_code", controller.UpdateJobStatus)
			secured.DELETE("/customer/:customer_code", controller.DeleteCustomer)
			secured.GET("/month/:month", controller.GetAllCustomerInASpecificMonth)
			secured.GET("/sum/:year", controller.GetTotalAnnualSubFeeOfAYear)
			secured.GET("/sum/month/:month", controller.GetTotalAnnualSubFeeOfAMonth)
			secured.PUT("/logout", controller.Logout)
		}
	}

	r.Run("0.0.0.0:" + port)
}
