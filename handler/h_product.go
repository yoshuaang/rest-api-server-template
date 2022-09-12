package handler

import (
	"encoding/json"
	"fmt"
	"golang-api-server/config"
	"golang-api-server/model"
	"golang-api-server/util"
	"io/ioutil"
	"log"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Is a key used for the Company departments object in the context
type KeyProducts struct{}

func ProductHandler(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusNoContent)
}

func GetProduct(writer http.ResponseWriter, request *http.Request) {
	DB, err_db := gorm.Open(mysql.Open(config.DSN_TEST), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err_db != nil {
		fmt.Println("[ERROR] Failed to connect to default database", err_db.Error())
	}

	var product model.Products
	productResult := DB.Find(&product).
		Where("active = ?", 2)
	if productResult.Error != nil {
		fmt.Println("[ERROR] Querying vendors", productResult.Error.Error())
	}

	util.ToJSON(product, writer)
}

func CreateProduct(writer http.ResponseWriter, request *http.Request) {
	target := "CreateProduct"
	log.Printf("[DEBUG] %v: receive product to be saved", target)

	products := model.Products{}
	b, _ := ioutil.ReadAll(request.Body)
	err := json.Unmarshal(b, &products)
	if err != nil {
		fmt.Printf("[ERROR] %v: deserializing company departments %v", target, err)

		writer.WriteHeader(http.StatusBadRequest)
		util.ToJSON(&util.ErrorResponse{
			Error: util.ErrorBody{
				Code:    util.ErrValidation.Error(),
				Message: err.Error(),
				Target:  target,
			},
		}, writer)
		return
	}

	DB, err_db := gorm.Open(mysql.Open(config.DSN_TEST), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err_db != nil {
		fmt.Println("[ERROR] Failed to connect to default database", err_db.Error())
	}

	trx := DB.CreateInBatches(products, 250)
	if trx.Error != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.ToJSON(&util.ErrorResponse{
			Error: util.ErrorBody{
				Code:    util.ErrInternalServer.Error(),
				Message: trx.Error.Error(),
				Target:  target,
			},
		}, writer)
		return
	}
	writer.WriteHeader(http.StatusOK)

	util.ToJSON(&util.SuccessResponse{
		Message: "Successfully saved the data",
	}, writer)
}
