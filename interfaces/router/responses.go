package router

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/xxarupkaxx/anke-two/repository/middleware"
	"github.com/xxarupkaxx/anke-two/usecase"
	"github.com/xxarupkaxx/anke-two/usecase/input"
	"net/http"
	"strconv"
)

type ResponseAPI interface {
	PostResponse(c echo.Context) error
	GetResponse(c echo.Context) error
	EditResponse(c echo.Context) error
	DeleteResponse(c echo.Context) error
}

type Response struct {
	usecase.ResponseUsecase
	middleware.IMiddleware
}

func NewResponseAPI(responseUsecase usecase.ResponseUsecase, IMiddleware middleware.IMiddleware) *Response {
	return &Response{ResponseUsecase: responseUsecase, IMiddleware: IMiddleware}
}

func (r *Response) PostResponse(c echo.Context) error {
	userID, err := r.GetUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
	}

	in := input.Responses{}
	in.UserID = userID

	if err := c.Bind(&in); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	statusCode, err := usecase.ValidateRequest(c, in)
	if err != nil {
		switch statusCode {
		case http.StatusBadRequest:
			c.Logger().Info(err)
			return echo.NewHTTPError(statusCode)
		case http.StatusInternalServerError:
			c.Logger().Error(err)
			return echo.NewHTTPError(statusCode)
		}
	}

	out, err := r.ResponseUsecase.PostResponse(c.Request().Context(), in)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, out)
}

func (r *Response) GetResponse(c echo.Context) error {
	responseID, err := strconv.Atoi(c.Param("responseID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	in := input.GetResponse{ResponseID: responseID}

	out, err := r.ResponseUsecase.GetResponse(c.Request().Context(), in)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, out)
}

func (r *Response) EditResponse(c echo.Context) error {
	responseID, err := strconv.Atoi(c.Param("responseID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	in := input.EditResponse{}
	in.ResponseID = responseID

	if err := c.Bind(&in); err != nil {
		c.Logger().Info(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	statusCode, err := usecase.ValidateRequest(c, in)
	if err != nil {
		switch statusCode {
		case http.StatusBadRequest:
			c.Logger().Info(err)
			return echo.NewHTTPError(statusCode)
		case http.StatusInternalServerError:
			c.Logger().Error(err)
			return echo.NewHTTPError(statusCode)
		}
	}

	err = r.ResponseUsecase.EditResponse(c.Request().Context(), in)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}

func (r *Response) DeleteResponse(c echo.Context) error {
	responseID, err := strconv.Atoi(c.Param("responseID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	in := input.DeleteResponse{ResponseID: responseID}

	err = r.ResponseUsecase.DeleteResponse(c.Request().Context(), in)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}
