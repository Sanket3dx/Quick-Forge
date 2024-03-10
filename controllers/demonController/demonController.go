package demonController

import (
	"net/http"
	mongooDbModel "quick_forge/database/mongooDb/models"
	mysql_models "quick_forge/database/mysql/models"
	"quick_forge/utils"

	"github.com/gin-gonic/gin"
)

func GetAllhandler(ctx *gin.Context) {
	routeInfo, exists := ctx.Get("routeInfo")
	if !exists {
		utils.HandleError(ctx, http.StatusInternalServerError, "Route info not found in context")
		return
	}
	route, ok := routeInfo.(utils.Route)
	if !ok {
		utils.HandleError(ctx, http.StatusInternalServerError, "Invalid routeInfo type")
		return
	}
	dbConfigInfo, exists := ctx.Get("dbConfig")
	if !exists {
		utils.HandleError(ctx, http.StatusInternalServerError, "dbConfig not found in context")
		return
	}
	dbConfig, ok := dbConfigInfo.(utils.Database)
	if !ok {
		utils.HandleError(ctx, http.StatusInternalServerError, "Invalid dbConfig type")
		return
	}

	var allData []map[string]interface{}
	var err error
	queryParms := utils.ConvertURLValuesToMapSlice(ctx.Request.URL.Query())
	switch dbConfig.DBType {
	case "mysql":
		allData, err = mysql_models.GetAllData(dbConfig, route.DBTableName, queryParms)
		if err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, "error fetching data from DB")
			return
		}
	case "mongo":
		allData, err = mongooDbModel.GetAllData(dbConfig, route.DBTableName, queryParms)
		if err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":   false,
		"data":    allData,
		"records": len(allData),
	})
}

func Gethandler(ctx *gin.Context) {
	routeInfo, exists := ctx.Get("routeInfo")
	if !exists {
		utils.HandleError(ctx, http.StatusInternalServerError, "Route info not found in context")
		return
	}
	route, ok := routeInfo.(utils.Route)

	if !ok {
		utils.HandleError(ctx, http.StatusInternalServerError, "Invalid routeInfo type")
		return
	}

	dbConfigInfo, exists := ctx.Get("dbConfig")
	if !exists {
		utils.HandleError(ctx, http.StatusInternalServerError, "dbConfig not found in context")
		return
	}
	dbConfig, ok := dbConfigInfo.(utils.Database)
	if !ok {
		utils.HandleError(ctx, http.StatusInternalServerError, "Invalid dbConfig type")
		return
	}

	var data map[string]interface{}
	var err error

	switch dbConfig.DBType {
	case "mysql":
		data, err = mysql_models.GetData(dbConfig, route, ctx.Param("arg"))
	case "mongo":
		data, err = mongooDbModel.GetData(dbConfig, route, ctx.Param("arg"))
	default:
		utils.HandleError(ctx, http.StatusInternalServerError, "Unsupported database type")
		return
	}

	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Error fetching data from DB")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  data,
	})
}

func Posthandler(ctx *gin.Context) {
	// Retrieve route information from context
	routeInfo, exists := ctx.Get("routeInfo")
	if !exists {
		utils.HandleError(ctx, http.StatusInternalServerError, "Route info not found in context")
		return
	}

	route, ok := routeInfo.(utils.Route)
	if !ok {
		utils.HandleError(ctx, http.StatusInternalServerError, "Invalid routeInfo type")
		return
	}

	dbConfigInfo, exists := ctx.Get("dbConfig")
	if !exists {
		utils.HandleError(ctx, http.StatusInternalServerError, "dbConfig not found in context")
		return
	}
	dbConfig, ok := dbConfigInfo.(utils.Database)
	if !ok {
		utils.HandleError(ctx, http.StatusInternalServerError, "Invalid dbConfig type")
		return
	}

	validatedRequestBody, exists := ctx.Get("validatedRequestBody")
	if !exists {
		utils.HandleError(ctx, http.StatusInternalServerError, "Validated request body not found in context")
		return
	}

	data, ok := validatedRequestBody.(map[string]interface{})
	if !ok {
		utils.HandleError(ctx, http.StatusInternalServerError, "Invalid validatedRequestBody type")
		return
	}

	// Insert data into the database
	var insertedID interface{}
	var err error

	switch dbConfig.DBType {
	case "mysql":
		insertedID, err = mysql_models.InsertData(dbConfig, route.DBTableName, data)
	case "mongo":
		insertedID, err = mongooDbModel.InsertData(dbConfig, route, data)
	default:
		utils.HandleError(ctx, http.StatusInternalServerError, "Unsupported database type")
		return
	}

	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Record insertion failed")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":     false,
		"record_id": insertedID,
	})
}

func Puthandler(ctx *gin.Context) {
	routeInfo, exists := ctx.Get("routeInfo")
	if !exists {
		utils.HandleError(ctx, http.StatusInternalServerError, "Route info not found in context")
		return
	}

	route, ok := routeInfo.(utils.Route)
	if !ok {
		utils.HandleError(ctx, http.StatusInternalServerError, "Invalid routeInfo type")
		return
	}

	dbConfigInfo, exists := ctx.Get("dbConfig")
	if !exists {
		utils.HandleError(ctx, http.StatusInternalServerError, "dbConfig not found in context")
		return
	}
	dbConfig, ok := dbConfigInfo.(utils.Database)
	if !ok {
		utils.HandleError(ctx, http.StatusInternalServerError, "Invalid dbConfig type")
		return
	}

	validatedRequestBody, exists := ctx.Get("validatedRequestBody")
	if !exists {
		utils.HandleError(ctx, http.StatusInternalServerError, "Validated request body not found in context")
		return
	}

	data, ok := validatedRequestBody.(map[string]interface{})
	if !ok {
		utils.HandleError(ctx, http.StatusInternalServerError, "Invalid validatedRequestBody type")
		return
	}

	var err error

	switch dbConfig.DBType {
	case "mysql":
		err = mysql_models.UpdateData(dbConfig, route, data, ctx.Param("arg"))
	case "mongo":
		err = mongooDbModel.UpdateData(dbConfig, route, data, ctx.Param("arg"))
	default:
		utils.HandleError(ctx, http.StatusInternalServerError, "Unsupported database type")
		return
	}

	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Record update failed")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Record updated successfully",
	})
}

func Deletehandler(ctx *gin.Context) {
	routeInfo, exists := ctx.Get("routeInfo")
	if !exists {
		utils.HandleError(ctx, http.StatusInternalServerError, "Route info not found in context")
		return
	}

	route, ok := routeInfo.(utils.Route)
	if !ok {
		utils.HandleError(ctx, http.StatusInternalServerError, "Invalid routeInfo type")
		return
	}

	dbConfigInfo, exists := ctx.Get("dbConfig")
	if !exists {
		utils.HandleError(ctx, http.StatusInternalServerError, "dbConfig not found in context")
		return
	}
	dbConfig, ok := dbConfigInfo.(utils.Database)
	if !ok {
		utils.HandleError(ctx, http.StatusInternalServerError, "Invalid dbConfig type")
		return
	}

	var err error

	switch dbConfig.DBType {
	case "mysql":
		err = mysql_models.DeleteData(dbConfig, route, ctx.Param("arg"))
	case "mongo":
		err = mongooDbModel.DeleteData(dbConfig, route, ctx.Param("arg"))
	default:
		utils.HandleError(ctx, http.StatusInternalServerError, "Unsupported database type")
		return
	}

	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Record deletion failed")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Record deleted successfully",
	})
}
