package main

import (
	"github.com/schidstorm/ffmpeg-go-server/application/lib/ffmpegLib"
	"github.com/schidstorm/ffmpeg-go-server/application/service"
	restframework "github.com/schidstorm/rest-framework"
	"gorm.io/driver/postgres"
)

func main() {
	app := restframework.LoadApplication(postgres.Open(GetSettings().Dsn))

	app.RegisterModel(service.Todo{})
	app.RegisterModel(service.FfmpegTask{})

	consumer := ffmpegLib.NewConsumer(app.Db)
	consumer.Run()

	app.Start(GetSettings().HttpListenEndpoint)

}
