package setting

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-ini/ini"
)

type App struct {
	JwtSecret       string
	JwtExpireHour   int
	PageSize        int
	RuntimeRootPath string
	ImagePrefixUrl  string
	ImageSavePath   string
	ImageMaxSize    int
	ImageAllowExts  []string
	LogSavePath     string
	LogSaveName     string
	LogFileExt      string
	TimeFormat      string
}

var AppSetting = &App{}

type Server struct {
	RunMode      string
	DebugMode    string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

type Database struct {
	Type     string
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

var DatabaseSetting = &Database{}

type Redis struct {
	Host        string
	Paswword    string
	MaxIdle     int
	MaxActive   int
	IdelTimeout time.Duration
}

var RedisSetting = &Redis{}

type Email struct {
	Email    string
	Password string
}

var EmailSetting = &Email{}

func Setup() {
	cfg, err := ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.int: %v'", err)
	}
	err = cfg.Section("app").MapTo(AppSetting)
	if err != nil {
		log.Fatalf("cfg.MapTo AppSetting err: %v", err)
	}
	// MB
	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024
	err = cfg.Section("server").MapTo(ServerSetting)
	if err != nil {
		log.Fatalf("cfg.MapTo ServerSetting error: %v", err)
	}
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.ReadTimeout * time.Second
	err = cfg.Section("database").MapTo(DatabaseSetting)
	if err != nil {
		log.Fatalf("cfg.MapTo DatabaseSetting error: %v", err)
	}
	err = cfg.Section("redis").MapTo(RedisSetting)
	if err != nil {
		log.Fatalf("cfg.MapTo RedisSetting err: %v", err)
	}
	err = cfg.Section("email").MapTo(EmailSetting)
	if err != nil {
		log.Fatalf("cfg.MapTo EmailSetting err: %v", err)
	}
}

func IsLocalTest() bool {
	host, err := os.Hostname()
	if err != nil {
		fmt.Errorf("isLocalTest Error: %v", err)
	} else {
		if host == "TWA01119484" || host == "linjianxinde-MacBook-Pro.local" {
			return true
		}
	}
	return false
}

func (serverSetting *Server) IsDebug() bool {
	return serverSetting.RunMode == "debug"
}
