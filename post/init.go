package post

import (
	. "session/global"

	. "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
)

func Init(t *pb.Request) (response *pb.Response) {

	param := *t.Param

	if *t.Module != MicroServiceName {

		param = *t.IR.Param
	}

	switch param {
	case "setsession":
		response = SetSessionapi(t)
	case "refresh_token":
		response = RefreshSesion(t)
	default:
		response = ErrorReturn(t, 404, "000012", "Missing argument")
	}

	return response

}
