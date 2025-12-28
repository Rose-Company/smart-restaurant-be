package cmd

import (
	"app-noti/common"
	"app-noti/server"
	logger2 "app-noti/services/logger"
	postgres3 "app-noti/services/postgres"
	"app-noti/services/rest_api_service"
	"context"
	"log"

	"github.com/spf13/cobra"
)

var restApiServiceCmd = &cobra.Command{
	Use:   "server",
	Short: "Start Your Server",
	Long:  "Let start a server with your opinion",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		common.FetchMasterErrData()

		loggerPkg := logger2.NewLogger("Logger")
		if err := loggerPkg.Run(); err != nil {
			log.Panic(err)
		}

		logger := loggerPkg.Get()

		start, _ := cmd.Flags().GetBool("start")

		if start {
			svr := server.NewServer("SupplierLoyaltyService", 8081)
			restHdl := rest_api_service.RestHandler(svr)
			err, postgres := postgres3.NewMainPostgres(common.PREFIX_MAIN_POSTGRES)
			if err != nil {
				logger.Error().Println("NewMainPostgres", err)
				return
			}

			svr.AddLogger(logger)
			svr.InitContext(ctx)
			svr.InitService(postgres)
			svr.AddHandler(restHdl)
			if err := svr.Run(); err != nil {
				logger.Error().Printf("Server is stopped by %v", err.Error())
			}
		}
	},
}
