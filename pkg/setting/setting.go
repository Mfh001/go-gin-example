package setting

import (
	"github.com/jinzhu/configor"
	"log"
	"time"

	"github.com/go-ini/ini"
)

type App struct {
	JwtSecret string
	PageSize  int
	PrefixUrl string

	RuntimeRootPath string

	ImageSavePath  string
	ImageMaxSize   int
	ImageAllowExts []string

	ExportSavePath string
	QrCodeSavePath string
	FontSavePath   string

	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string
	WXAppID     string
	WXSecret    string
}

var AppSetting = &App{}

type Server struct {
	RunMode          string
	HttpPort         int
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	SSLOpen          bool
	CrtPath          string
	KeyPath          string
	WXCodeExpireTime int
}

var ServerSetting = &Server{}

type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

var DatabaseSetting = &Database{}

type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
	LockTimeout int
}

var RedisSetting = &Redis{}

var Platform = struct {
	Games []struct {
		Name   string
		Idx    int
		IsOpen bool
	}
	OrderTypes []struct {
		Name   string
		Idx    int
		IsOpen bool
	}
	ServerZones []struct {
		Name string
		Idx  int
		Num  int
	}
	InsteadTypes []struct {
		Name string
		Idx  int
	}
	Levels []struct {
		Name  string
		Idx   int
		Stars int
		Price int
	}
}{}

type LevelCell struct {
	Idx   int
	Price int
}

var PlatFormLevelAll []LevelCell

var cfg *ini.File

// Setup initialize the configuration instance
func Setup() {
	var err error
	cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
	}

	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)
	mapTo("redis", RedisSetting)

	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second

	configor.Load(&Platform, "conf/platform.yaml")
	levelInit()
	//fmt.Printf("config: %#v", Platform)
}

func levelInit() {
	for i := 0; i < len(Platform.Levels); i++ {
		for j := 0; j <= Platform.Levels[i].Stars; j++ {
			idx := Platform.Levels[i].Idx*1000 + j
			l := LevelCell{
				Idx:   idx,
				Price: Platform.Levels[i].Price,
			}
			PlatFormLevelAll = append(PlatFormLevelAll, l)
		}
	}
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
