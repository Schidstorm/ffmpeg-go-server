package ffmpegLib

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"os/user"
)

type Settings struct {
	ConversionDestinationDirectory string
	UploadDestinationDirectory     string
	HttpListenEndpoint             string
	Dsn                            string
}

var settings *Settings = initSettings()

func initSettings() *Settings {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}

	loadConfigFile(currentUser.HomeDir)

	if viper.GetString("PostgresPasswordFile") != "" {
		buffer, err := ioutil.ReadFile(viper.GetString("PostgresPasswordFile"))
		if err != nil {
			panic(err)
		}
		viper.Set("PostgresPassword", string(buffer))
	}

	return &Settings{
		ConversionDestinationDirectory: viper.GetString("ConversionDestinationDirectory"),
		UploadDestinationDirectory:     viper.GetString("UploadDestinationDirectory"),
		HttpListenEndpoint:             viper.GetString("HttpListenEndpoint"),
		Dsn: fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Europe/Berlin",
			viper.GetString("PostgresHostname"),
			viper.GetString("PostgresUsername"),
			viper.GetString("PostgresPassword"),
			viper.GetString("PostgresDbName"),
			viper.GetInt32("PostgresPort"),
		),
	}
}

func loadConfigFile(homePath string) {
	viper.SetDefault("ConversionDestinationDirectory", homePath)
	viper.SetDefault("UploadDestinationDirectory", homePath)
	viper.SetDefault("HttpListenEndpoint", "0.0.0.0:8080")
	viper.SetDefault("PostgresUsername", "ffmpeg")
	viper.SetDefault("PostgresHostname", "localhost")
	viper.SetDefault("PostgresPassword", "ffmpeg")
	viper.SetDefault("PostgresPasswordFile", "")
	viper.SetDefault("PostgresDbName", "ffmpeg")
	viper.SetDefault("PostgresPort", "9920")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("FFMPEG_GO_SERVER_")
	viper.AddConfigPath("/etc/ffmpeg-server/")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return
		}

		logrus.Errorln(err)
	}
}

func GetSettings() *Settings {
	return settings
}
