package postgres

import (
	"fmt"
	"gorm.io/gorm/logger"
)

type ConfigureParams struct {
	User      string
	Password  string
	Host      string
	Port      int
	Database  string
	Params    string
	DebugMode logger.LogLevel
}

func GetPostgresUri(params ConfigureParams) string {
	if params.Port > 0 {
		return fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?%v", params.User, params.Password, params.Host, params.Port, params.Database, params.Params)
	}
	return fmt.Sprintf("postgresql://%v:%v@%v/%v?%v", params.User, params.Password, params.Host, params.Database, params.Params)
}
