package main

import (
	"fmt"

	. "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"

	ad "session/admin"
	gt "session/get"
	. "session/global"
	pt "session/post"
	. "session/version"
)

func Init(t *pb.Request) (response *pb.Response) {
	SetErrorLog("Im Init of Amy Session")

	method := *t.Method
	param := *t.Param

	if *t.Module != MicroServiceName {
		method = *t.IR.Method
		param = *t.IR.Param
	}

	switch param {
	case "admin":
		return admincheck(t)
	}

	SetErrorLog(fmt.Sprintf("Method: %s", method))
	SetErrorLog(fmt.Sprintf("Param: %s", param))

	switch method {
	case "GET":
		switch param {
		case "info":
			response = info(t)
		case "health":
			response = health(t)
		default:
			response = gt.Init(t)
		}
	case "POST":
		response = pt.Init(t)
	default:
		response = ErrorReturn(t, 404, "000012", "Missing argument")

	}

	return response

}

func info(t *pb.Request) (response *pb.Response) {
	ans := make(map[string]interface{})
	ans["pluginname"] = "Session microservice"
	ans["version"] = VERSIONPLUGIN
	ans["description"] = ""
	response = Interfacetoresponse(t, ans)
	return response
}

func health(t *pb.Request) (response *pb.Response) {
	ans := make(map[string]interface{})
	ans["health"] = "OK"
	response = Interfacetoresponse(t, ans)
	return response
}

func admincheck(t *pb.Request) (response *pb.Response) {

	if *t.IsAdmin != 1 {
		response = ErrorReturn(t, 401, "000012", "You have no admin rights")
	}

	return ad.Init(t)

}
