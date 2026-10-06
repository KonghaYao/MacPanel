package capabilities

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

const NotSupportedReason = "NOT_SUPPORTED_ON_DARWIN"

var ErrNotSupportedOnMac = errors.New("this feature is not supported on macOS")

type notSupportedResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	ErrorCode string `json:"errorCode,omitempty"`
	Reason    string `json:"reason,omitempty"`
}

func IsNotSupportedOnMac(err error) bool {
	return errors.Is(err, ErrNotSupportedOnMac)
}

func RespondNotSupportedOnDarwin(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, notSupportedResponse{
		Code:      http.StatusNotImplemented,
		Message:   "此功能在 macOS 上不可用",
		ErrorCode: NotSupportedReason,
		Reason:    NotSupportedReason,
	})
	c.Abort()
}
