package entrypoint

import (
  "fmt"
	. "session/version"
	. "session/global"

	"github.com/spf13/viper"
)

func Init() {
	EntryPoint()
	updateversion()
}

func updateversion() {
setingskey := fmt.Sprintf("%s.entrypointversion", MicroServiceName)
viper.Set(setingskey, VERSIONPLUGIN)
}
