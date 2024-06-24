package entrypoint

import (
	. "session/model"

	. "github.com/gogufo/gufo-api-gateway/gufodao"
	"github.com/spf13/viper"
)

func CheckDBStructure() {
	//Check DB and table config
	db, err := ConnectDBv2()
	if err != nil {
		SetErrorLog("dbstructure.go:81: " + err.Error())
		//return "error with db"
	}

	dbtype := viper.GetString("database.type")

	if !db.Conn.Migrator().HasTable(&RefreshTokens{}) {
		if dbtype == "mysql" {
			db.Conn.Set("gorm:table_options", "ENGINE=InnoDB;").Migrator().CreateTable(&RefreshTokens{})
		} else {
			db.Conn.Migrator().CreateTable(&RefreshTokens{})
		}
	}

	if !db.Conn.Migrator().HasTable(&APITokens{}) {
		if dbtype == "mysql" {
			db.Conn.Set("gorm:table_options", "ENGINE=InnoDB;").Migrator().CreateTable(&APITokens{})
		} else {
			db.Conn.Migrator().CreateTable(&APITokens{})
		}
	}

	if !db.Conn.Migrator().HasTable(&ImpersonateTokens{}) {
		if dbtype == "mysql" {
			db.Conn.Set("gorm:table_options", "ENGINE=InnoDB;").Migrator().CreateTable(&ImpersonateTokens{})
		} else {
			db.Conn.Migrator().CreateTable(&ImpersonateTokens{})
		}
	}

}
