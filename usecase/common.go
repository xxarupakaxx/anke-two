package usecase

import (
	"github.com/labstack/echo/v4"
	mymiddleware "github.com/xxarupkaxx/anke-two/interfaces/middleware"
	"net/http"
)

func ValidateRequest(c echo.Context, request interface{}) (int, error){
	validate, err := mymiddleware.GetValidator(c)
	if err != nil {
	return http.StatusInternalServerError, err
	}

	err = validate.StructCtx(c.Request().Context(), request)
	if err != nil {
	return http.StatusBadRequest, err
	}

	return http.StatusOK, nil

}