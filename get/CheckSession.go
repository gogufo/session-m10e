package get

import (
	"fmt"
	fn "session/functions"

	. "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
)

func CheckSession(t *pb.Request) (response *pb.Response) {
	SetErrorLog("CheckSession Function")

	ans := make(map[string]interface{})

	token := *t.Token

	SetErrorLog(fmt.Sprintf("Token: %s", token))

	ans = fn.UpdateSession(token)

	response = Interfacetoresponse(t, ans)
	return response

}
