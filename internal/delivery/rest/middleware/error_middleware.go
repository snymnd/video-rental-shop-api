package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"vrs-api/internal/customerrors"
	"vrs-api/internal/dto"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) == 0 {
			return
		}
		err := ctx.Errors[0]

		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			fieldErrors := make([]string, 0)

			for _, fe := range ve {
				log.Println(fe.Error())

				fieldErrors = append(fieldErrors, fmt.Sprintf("invalid input on field %s", fe.Field()))
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
			log.Println(ce.ErrorLog)

			ctx.JSON(ce.GetHTTPErrorCode(), dto.ResponseError(dto.ErrorResponse{
				Message: ce.ErrorMessage,
			}))
			return
		}

		log.Println(err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ResponseError(dto.ErrorResponse{Message: "something went wrong"}))
	}
}
