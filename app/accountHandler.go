package app

import (
	"net/http"

	"github.com/amrshaban2005/microservice-api/dto"
	"github.com/amrshaban2005/banking-lib/logger"
	"github.com/amrshaban2005/microservice-api/service"
	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	service service.AccountService
}

func (ah *AccountHandler) NewAccount(c *gin.Context) {

	var req dto.NewAccountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Error while read account params " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}
	req.CustomerId = c.Param("id")
	newAccount, err := ah.service.NewAccount(c.Request.Context(), req)
	if err != nil {
		c.JSON(err.Code, gin.H{"error": err.AsMessage()})
		return
	}
	c.JSON(http.StatusCreated, newAccount)
}

func (ah *AccountHandler) MakeTransaction(c *gin.Context) {
	var req dto.NewTransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Error while read account params " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	req.AccountId = c.Param("account_id")

	transactionResponse, err := ah.service.MakeTransaction(c.Request.Context(), req)
	if err != nil {
		c.JSON(err.Code, gin.H{"error": err.AsMessage()})
		return
	}
	c.JSON(http.StatusCreated, transactionResponse)
}
