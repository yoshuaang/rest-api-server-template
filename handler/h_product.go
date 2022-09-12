package handler

import (
	"fmt"
	"golang-api-server/config"
	"golang-api-server/model"
	"golang-api-server/util"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

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
