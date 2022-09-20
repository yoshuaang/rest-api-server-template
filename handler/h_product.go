package handler

import (
	"encoding/json"
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

// Product handler for:
// - Receiving product to be saved, updated, or deleted
type ProductHandler struct {
	l *log.Logger
	v *util.Validation
}

// Return a new product handler with the given logger & validation
func NewProductHandler(l *log.Logger, v *util.Validation) *ProductHandler {
	return &ProductHandler{l, v}
}

func (p *ProductHandler) GetProduct(writer http.ResponseWriter, request *http.Request) {
	DB, err_db := gorm.Open(mysql.Open(config.DSN_TEST), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err_db != nil {
		p.l.Println("[ERROR] Failed to connect to default database", err_db.Error())
	}

	var product model.Products
	productResult := DB.Find(&product).
		Where("active = ?", 2)
	if productResult.Error != nil {
		p.l.Println("[ERROR] Querying products", productResult.Error.Error())
	}

	util.ToJSON(product, writer)
}

func (p *ProductHandler) CreateProduct(writer http.ResponseWriter, request *http.Request) {
	target := "CreateProduct"
	p.l.Printf("[DEBUG] %v: receive products to be saved", target)

	products := model.Products{}
	b, _ := ioutil.ReadAll(request.Body)
	err := json.Unmarshal(b, &products)
	if err != nil {
		p.l.Printf("[ERROR] %v: deserializing products %v", target, err)

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
		p.l.Println("[ERROR] Failed to connect to default database", err_db.Error())
	}

	trx := DB.CreateInBatches(products, 50)
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

func (p *ProductHandler) UpdateProduct(rw http.ResponseWriter, rq *http.Request) {
	target := "UpdateProduct"
	p.l.Printf("[DEBUG] %v: receive products to be updated", target)

	products := model.Products{}
	b, _ := ioutil.ReadAll(rq.Body)
	err := json.Unmarshal(b, &products)
	if err != nil {
		p.l.Printf("[ERROR] %v: deserializing products %v", target, err)

		rw.WriteHeader(http.StatusBadRequest)
		util.ToJSON(&util.ErrorResponse{
			Error: util.ErrorBody{
				Code:    util.ErrValidation.Error(),
				Message: err.Error(),
				Target:  target,
			},
		}, rw)
		return
	}

	DB, err_db := gorm.Open(mysql.Open(config.DSN_TEST), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err_db != nil {
		p.l.Println("[ERROR] Failed to connect to default database", err_db.Error())
	}

	txDB := DB.Begin()
	res := txDB.Save(&products)
	if res.Error != nil {
		p.l.Printf("[ERROR] %v: update product %v", target, res.Error.Error())
		txDB.Rollback()

		rw.WriteHeader(http.StatusInternalServerError)
		util.ToJSON(&util.ErrorResponse{
			Error: util.ErrorBody{
				Code:    util.ErrInternalServer.Error(),
				Message: res.Error.Error(),
				Target:  target,
			},
		}, rw)
		return
	}

	res = txDB.Commit()
	if res.Error != nil {
		p.l.Printf("[ERROR] %v: commit update product %v", target, res.Error.Error())
		txDB.Rollback()

		rw.WriteHeader(http.StatusInternalServerError)
		util.ToJSON(&util.ErrorResponse{
			Error: util.ErrorBody{
				Code:    util.ErrInternalServer.Error(),
				Message: res.Error.Error(),
				Target:  target,
			},
		}, rw)
		return
	}

	// Returns success message
	rw.WriteHeader(http.StatusOK)

	util.ToJSON(&util.SuccessResponse{
		Message: "Successfully updated the data",
	}, rw)
}

func (p *ProductHandler) DeleteProduct(rw http.ResponseWriter, rq *http.Request) {
	target := "DeleteProduct"
	p.l.Printf("[DEBUG] %v: receive products to be deleted", target)

	products := model.Products{}
	b, _ := ioutil.ReadAll(rq.Body)
	err := json.Unmarshal(b, &products)
	if err != nil {
		p.l.Printf("[ERROR] %v: deserializing products %v", target, err)

		rw.WriteHeader(http.StatusBadRequest)
		util.ToJSON(&util.ErrorResponse{
			Error: util.ErrorBody{
				Code:    util.ErrValidation.Error(),
				Message: err.Error(),
				Target:  target,
			},
		}, rw)
		return
	}

	DB, err_db := gorm.Open(mysql.Open(config.DSN_TEST), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err_db != nil {
		p.l.Println("[ERROR] Failed to connect to default database", err_db.Error())
	}

	res := DB.Delete(&products)
	if res.Error != nil {
		p.l.Printf("[ERROR] %v: commit update product %v", target, res.Error.Error())
		res.Rollback()

		rw.WriteHeader(http.StatusInternalServerError)
		util.ToJSON(&util.ErrorResponse{
			Error: util.ErrorBody{
				Code:    util.ErrInternalServer.Error(),
				Message: res.Error.Error(),
				Target:  target,
			},
		}, rw)
		return
	}

	// Returns success message
	rw.WriteHeader(http.StatusOK)

	util.ToJSON(&util.SuccessResponse{
		Message: "Successfully deleted the data",
	}, rw)
}
