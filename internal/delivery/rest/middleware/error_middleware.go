package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"vrs-api/internal/customerrors"
	"vrs-api/internal/dto"
	"vrs-api/internal/util/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		log := logger.GetLogger()

		if len(ctx.Errors) == 0 {
			return
		}
		err := ctx.Errors[0]

		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			fieldErrors := make([]dto.DetailsError, 0)

			for _, fe := range ve {
				log.Error(fe.Error())

				fieldErrors = append(fieldErrors, dto.DetailsError{
					Title:   fe.Field(),
					Message: fmt.Sprintf("invalid input on field %s", fe.Field()),
				})
			}

			errorResponse := dto.ErrorResponse{
				Message: "input validation errors",
				Details: fieldErrors,
			}

			ctx.JSON(http.StatusBadRequest, dto.ResponseError(errorResponse))
			return
		}

		var ce *customerrors.CustomError
		if errors.As(err, &ce) {
			log.Error(ce.ErrorLog)

			ctx.JSON(ce.GetHTTPErrorCode(), dto.ResponseError(dto.ErrorResponse{
				Message: ce.ErrorMessage,
				Details: ce.Details,
			}))
			return
		}

		log.Error(err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ResponseError(dto.ErrorResponse{Message: "something went wrong"}))
	}
}
