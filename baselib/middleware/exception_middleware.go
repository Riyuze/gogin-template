package middleware

import (
	"fmt"
	"gogin-template/baselib/dto"
	"gogin-template/baselib/exception"
	"gogin-template/baselib/identifier"
	"gogin-template/bootstrap"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ExceptionMiddleware(cfg *bootstrap.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			code := strconv.Itoa(http.StatusInternalServerError)
			message := "Internal Server Error"
			httpStatus := http.StatusInternalServerError
			traceId := identifier.GetTraceId(c)
			logReff := identifier.GetLogReff(c)

			switch e := err.(type) {
			case *exception.ErrorException:
				code = e.ErrorCode
				message = e.ErrorMessage
				httpStatus = e.HttpStatusCode
			default:
			}

			cfg.Logger().Errorf("%s - %s - %s", code, message, err)

			resp := dto.Response[any]{
				LogReff:         logReff,
				TraceId:         fmt.Sprint(traceId),
				ResponseCode:    code,
				ResponseMessage: message,
			}
			c.AbortWithStatusJSON(httpStatus, resp)
		}
	}
}
