package docker

import (
	"strings"

	"github.com/1Panel-dev/1Panel/agent/i18n"
)

func PullErrorPrefix(errMsg string) string {
	if errMsg == "" {
		return ""
	}
	if strings.Contains(errMsg, "no matching manifest") {
		return i18n.GetMsgByKey("PullImagePlatformMismatch") + ":"
	}
	if strings.Contains(errMsg, "no such host") {
		return i18n.GetMsgByKey("ErrNoSuchHost") + ":"
	}
	if strings.Contains(errMsg, "Error response from daemon") {
		return i18n.GetMsgByKey("PullImageTimeout") + ":"
	}
	return ""
}
