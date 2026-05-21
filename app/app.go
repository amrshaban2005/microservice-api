package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/amrshaban2005/microservice-api/domain"
	"github.com/amrshaban2005/microservice-api/logger"
	"github.com/amrshaban2005/microservice-api/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Start() {
	router := gin.Default()

	pool := getDBClient()

	customerRepositoryDB := domain.NewCustomerRepositoryDB(pool)
	accountRepositoryDB := domain.NewAccountRepositoryDB(pool)

	ch := &CustomerHandler{service.NewCustomerService(customerRepositoryDB)}
	ah := &AccountHandler{service.NewAccountService(accountRepositoryDB)}

	router.GET("/customers", ch.GetAllCustomers)
	router.GET("/customers/:id", ch.GetCustomer)
	router.POST("/customers/:id/account", ah.NewAccount)
	router.POST("/customers/:id/account/:account_id", ah.MakeTransaction)

	// starting the server
	address := os.Getenv("SERVER_ADDRESS")
	port := os.Getenv("SERVER_PORT")

	logger.Info(fmt.Sprintf("Starting server on:%s:%s", address, port))

	log.Fatal(router.Run(fmt.Sprintf("%s:%s", address, port)))
}

func getDBClient() *pgxpool.Pool {
	dbUser := os.Getenv("DB_USER")
	dbPasswd := os.Getenv("DB_PASSWD")
	dbAddr := os.Getenv("DB_ADDR")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPasswd, dbAddr, dbPort, dbName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		panic(err)
	}

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		panic(err)
	}
	return pool

}
