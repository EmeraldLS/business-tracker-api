package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/EmeraldLS/phsps-api/code"
	"github.com/EmeraldLS/phsps-api/config"
	"github.com/EmeraldLS/phsps-api/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-module/carbon"
)

func InsertCustomer(c *gin.Context) {
	var customer model.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "bind_error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}

	validate := validator.New()
	if err := validate.Struct(customer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "struct_error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}

	if err := config.ValidateEmail(customer.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "email_error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}

	customer.CustomerCode = code.GetMaxCustomerCode() + 1
	customer.CustomerID = code.GenCustomerID(customer.CustomerCode)
	customer.JoinDate = carbon.Now().ToDateTimeString()
	customer.RenewalDate = carbon.Now().AddYear().ToDateTimeString()

	customer.RenewalMonth = strings.ToLower(carbon.Now().AddYear().ToMonthString())
	customer.JobStatusCode = 0
	customer.JobStatus = "Job received"
	customer.UpdatedAt = carbon.Now().ToDateTimeString()
	customer.JoinYear = carbon.Now().Year()
	if err := config.InsertCustomer(customer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "insertion_error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, customer)
}

func GetAllCustomter(c *gin.Context) {
	customers, err := config.GetAllCustomter()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "retrieving_error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, customers)
}

func SearchCustomerByBusinessName(c *gin.Context) {
	s := c.Query("s")
	customers, err := config.SearchCustomerByBusinessName(s)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "retrieving_error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, customers)
}

func GetCustomerByID(c *gin.Context) {
	id := c.Param("customer_code")
	customerCode, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	customer, err := config.GetCustomerByID(customerCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "retrieving_error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, customer)
}

func UpdateJobStatus(c *gin.Context) {
	statusCodeString := c.Param("job_status_code")
	statusCode, err := strconv.Atoi(statusCodeString)
	if err != nil {
		fmt.Println(err)
		return
	}
	customerCodeString := c.Param("customer_code")
	customerCode, err := strconv.Atoi(customerCodeString)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = config.UpdateJobStatus(statusCode, customerCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "update_error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": "success",
		"message":  "Job Status Update successful.",
	})
}

func DeleteCustomer(c *gin.Context) {
	customerCodeString := c.Param("customer_code")
	customerCode, err := strconv.Atoi(customerCodeString)
	if err != nil {
		fmt.Println(err)
		return
	}
	count, err := config.DeleteCustomer(customerCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "deletion_error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": "success",
		"message":  "Customer deleted successfully.",
	})
	fmt.Println(count)
}

func GetAllCustomerInASpecificMonth(c *gin.Context) {
	month := c.Param("month")
	customers, err := config.GetAllCustomerInASpecificMonth(strings.ToLower(month))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "retrieving_error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, customers)
}

func GetTotalAnnualSubFeeOfAYear(c *gin.Context) {
	yearString := c.Param("year")
	year, _ := strconv.Atoi(yearString)
	message, err := config.GetTotalAnnualSubFeeForAYear(year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "retrieving_error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": "success",
		"message":  message,
	})
}

func GetTotalAnnualSubFeeOfAMonth(c *gin.Context) {
	month := c.Param("month")
	message, err := config.GetTotalAnnualSubFeeOfAMonth(strings.ToLower(month))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "retrieving_error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": "success",
		"message":  message,
	})
}
