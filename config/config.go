package config

import (
	"github.com/devfeel/dotweb-start/const"
	"github.com/devfeel/dotweb-start/core"
	"github.com/devfeel/dotweb-start/util/file"
	"github.com/devfeel/dotweb-start/util/json"
	"encoding/xml"
	"github.com/devfeel/dotlog"
	"io/ioutil"
	"strconv"
)

var (
	CurrentConfig  *AppConfig
	CurrentBaseDir string
	innerLogger    dotlog.Logger
	appSetMap      *core.CMap
	allowIPMap     *core.CMap
	redisMap       *core.CMap
)

func SetBaseDir(baseDir string) {
	CurrentBaseDir = baseDir
}

//初始化配置文件
func InitConfig(configFile string) *AppConfig {
	innerLogger = dotlog.GetLogger(_const.LoggerName_Inner)
	CurrentBaseDir = file.GetCurrentDirectory()
	innerLogger.Info("AppConfig::InitConfig 配置文件[" + configFile + "]开始...")
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		innerLogger.Warn("AppConfig::InitConfig 配置文件[" + configFile + "]无法解析 - " + err.Error())
		panic(err)
	}

	var result AppConfig
	err = xml.Unmarshal(content, &result)
	if err != nil {
		innerLogger.Warn("AppConfig::InitConfig 配置文件[" + configFile + "]解析失败 - " + err.Error())
		panic(err)
	}

	//init config base
	CurrentConfig = &result

	//init AppConfig
	innerLogger.Info("AppConfig::InitConfig Load AppSet Start")
	tmpAppSetMap := core.NewCMap()
	for _, v := range result.AppSets {
		tmpAppSetMap.Set(v.Key, v.Value)
		innerLogger.Info("AppConfig::InitConfig Load AppSet => " + jsonutil.GetJsonString(&v))
	}
	appSetMap = tmpAppSetMap
	innerLogger.Info("AppConfig::InitConfig Load AppSet Finished [" + strconv.Itoa(appSetMap.Len()) + "]")

	//init redisConfig
	innerLogger.Info("AppConfig::InitConfig Start Load RedisInfo")
	tmpRedisMap := core.NewCMap()
	for k, v := range result.Redises {
		tmpRedisMap.Set(v.ID, result.Redises[k])
		innerLogger.Info("AppConfig::InitConfig Load RedisInfo => " + jsonutil.GetJsonString(v))
	}
	redisMap = tmpRedisMap
	innerLogger.Info("AppConfig::InitConfig Finish Load RedisInfo")

	innerLogger.Info("AppConfig::InitConfig 配置文件[" + configFile + "]完成")

	CurrentConfig.ConfigPath = GetAppConfig("ConfigPath")

	return CurrentConfig
}

func GetAppConfig(key string) string {
	return appSetMap.GetString(key)
}

func GetAppSetMap() *core.CMap {
	return appSetMap
}

func GetRedisInfo(redisID string) (*RedisInfo, bool) {
	redis, exists := redisMap.Get(redisID)
	if exists {
		return redis.(*RedisInfo), exists
	} else {
		return nil, false
	}
}

//检测IP是否被允许访问
func CheckAllowIP(ip string) bool {
	return allowIPMap.Exists(ip)
}
