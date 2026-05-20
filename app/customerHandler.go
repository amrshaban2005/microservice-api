package app

import (
	"net/http"

	"github.com/amrshaban2005/microservice-api/service"
	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	service service.CustomerService
}

func (ch *CustomerHandler) GetAllCustomers(c *gin.Context){
	status := c.Query("status")
	customers,err :=ch.service.GetAllCustomers(c.Request.Context(),status)
	if err != nil {
		c.JSON(err.Code, gin.H{"error": err.AsMessage()})
		return
	}
	c.JSON(http.StatusOK, customers)
}

func (ch *CustomerHandler) GetCustomer(c *gin.Context){
	id := c.Param("id")

	customer,err :=ch.service.GetCustomer(c.Request.Context(),id)
	if err != nil {
		c.JSON(err.Code, gin.H{"error": err.AsMessage()})
		return
	}
	c.JSON(http.StatusOK, customer)
}
