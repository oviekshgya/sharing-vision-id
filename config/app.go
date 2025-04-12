package config

import (
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log"
	"sharing-vision-id/pkg"
)

var Config *Configuration

type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfiguration
	Oauth    OauthConfiguration
}
type ServerConfiguration struct {
	AppName string
	Port    string
	Secret  string
	Mode    string
	Env     string
}

type DatabaseConfiguration struct {
	Host                 string
	Port                 string
	Username             string
	Password             string
	Database             string
	MAX_IDLE_CONNECTIONS int
	MAX_OPEN_CONNECTIONS int
	MaxLifeTimeSeconds   int
	DRIVER               string
}

type OauthConfiguration struct {
	OAUTH_CLIENT_ID     string
	OAUTH_CLIENT_SECRET string
	OAUTH_REDIRECT_URL  string
}

func AppConfig() {
	configuration := &Configuration{}

	//viper.SetConfigFile("./.env")
	//viper.SetConfigType("env")
	//viper.AutomaticEnv()
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error saat membaca file .env: %v", err)
	}

	configuration.Server.AppName = viper.GetString("SERVICE_NAME")
	configuration.Server.Port = viper.GetString("SERVICE_PORT")
	configuration.Server.Secret = viper.GetString("SERVICE_SECRET")
	configuration.Server.Mode = viper.GetString("SERVICE_MODE")
	configuration.Server.Env = viper.GetString("SERVICE_ENV")

	configuration.Database.Host = viper.GetString("DB_HOST")
	configuration.Database.Port = viper.GetString("DB_PORT")
	configuration.Database.Username = viper.GetString("DB_USER")
	configuration.Database.Password = viper.GetString("DB_PASSWORD")
	configuration.Database.Database = viper.GetString("DB_NAME")
	configuration.Database.MAX_IDLE_CONNECTIONS = viper.GetInt("DB_MAX_IDLE_CONNECTIONS")
	configuration.Database.MAX_OPEN_CONNECTIONS = viper.GetInt("DB_MAX_OPEN_CONNECTIONS")
	configuration.Database.MaxLifeTimeSeconds = viper.GetInt("DB_MAX_LIFE_TIME")
	configuration.Database.DRIVER = viper.GetString("DB_DRIVER")

	configuration.Oauth.OAUTH_CLIENT_ID = viper.GetString("OAUTH_CLIENT_ID")
	configuration.Oauth.OAUTH_CLIENT_SECRET = viper.GetString("OAUTH_CLIENT_SECRET")
	configuration.Oauth.OAUTH_REDIRECT_URL = viper.GetString("OAUTH_REDIRECT_URL")
	Config = configuration
	pkg.GoogleOAuthConfig = &oauth2.Config{
		ClientID:     Config.Oauth.OAUTH_CLIENT_ID,
		ClientSecret: configuration.Oauth.OAUTH_CLIENT_SECRET,
		RedirectURL:  Config.Oauth.OAUTH_REDIRECT_URL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

}

func GetConfig() *Configuration {
	return Config
}
