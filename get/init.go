package get

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
	case "checksession":
		response = CheckSession(t)
	default:
		response = ErrorReturn(t, 404, "000012", "Missing argument")
	}

	return response

}
