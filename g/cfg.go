package g

import (
	"encoding/json"
	"fmt"
	"github.com/toolkits/file"
	"log"
	"sync"
)

type HttpConfig struct {
	Debug   bool   `json:"debug"`
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type QueueConfig struct {
	Sms        string `json:"sms"`
	Mail       string `json:"mail"`
	QQ         string `json:"qq"`
	Serverchan string `json:"serverchan"`
}

type RedisConfig struct {
	Addr                string   `json:"addr"`
	MaxIdle             int      `json:"maxIdle"`
	HighQueues          []string `json:"highQueues"`
	LowQueues           []string `json:"lowQueues"`
	UserSmsQueue        string   `json:"userSmsQueue"`
	UserMailQueue       string   `json:"userMailQueue"`
	UserQQQueue         string   `json:"userQQQueue"`
	UserServerchanQueue string   `json:"userServerchanQueue"`
}

type ApiConfig struct {
	Portal string `json:"portal"`
	Uic    string `json:"uic"`
	Links  string `json:"links"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Account  string `json:"account"`
	Password string `json:"password"`
	Options  string `json:"options"`
}

type UicConfig struct {
	Addr  string `json:"addr"`
	Table string `json:"table"`
	Idle  int    `json:"idle"`
	Max   int    `json:"max"`
}

type FalconPortalConfig struct {
	Addr  string `json:"addr"`
	Table string `json:"table"`
	Idle  int    `json:"idle"`
	Max   int    `json:"max"`
}

type ShortcutConfig struct {
	FalconPortal     string `json:"falconPortal"`
	FalconDashboard  string `json:"falconDashboard"`
	GrafanaDashboard string `json:"grafanaDashboard"`
	FalconAlarm      string `json:"falconAlarm"`
	FalconUIC        string `json:"falconUIC"`
}
type GlobalConfig struct {
	Debug        bool                `json:"debug"`
	UicToken     string              `json:"uicToken"`
	Http         *HttpConfig         `json:"http"`
	Queue        *QueueConfig        `json:"queue"`
	Redis        *RedisConfig        `json:"redis"`
	Api          *ApiConfig          `json:"api"`
	Shortcut     *ShortcutConfig     `json:"shortcut"`
	Db           *DatabaseConfig     `json:"db"`
	Uic          *UicConfig          `json:"uic"`
	FalconPortal *FalconPortalConfig `json:"falcon_portal"`
	RedirectUrl  string              `json:"redirectUrl"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}
	db := c.Db
	uicdb := c.Uic
	fpdb := c.FalconPortal
	uicdb.Addr = fmt.Sprintf("%s:%s@%s(%s:%d)/%s?%s", db.Account, db.Password, db.Protocol, db.Host, db.Port, uicdb.Table, db.Options)
	fpdb.Addr = fmt.Sprintf("%s:%s@%s(%s:%d)/%s?%s", db.Account, db.Password, db.Protocol, db.Host, db.Port, fpdb.Table, db.Options)
	c.Uic = uicdb
	c.FalconPortal = fpdb
	configLock.Lock()
	defer configLock.Unlock()
	config = &c
	log.Println("read config file:", cfg, "successfully")
}
