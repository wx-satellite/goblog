package config

import (
	"goblog/pkg/logger"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

// Viper Viper 库实例
var Viper *viper.Viper

// StrMap 简写 —— map[string]interface{}
type StrMap map[string]interface{}

// init() 函数在 import 的时候立刻被加载
func init() {
	// 1. 初始化 Viper 库
	Viper = viper.New()
	// 2. 设置文件名称
	Viper.SetConfigName(".env")

	// 如果配置文件是：1.env，那么只设置Viper.SetConfigName("1")即可不需要再设置Viper.SetConfigType
	// 因为，viper 可以根据 1.env 的后缀名 .env 自行解析
	// 本案例中配置文件是 .env 相当于没有指定后缀，因此需要设置 SetConfigType

	// 3. 配置类型，支持 "json", "toml", "yaml", "yml", "properties",
	//             "props", "prop", "env", "dotenv"
	Viper.SetConfigType("env")

	//viper.SetConfigName("config") // 配置文件名称(无扩展名)
	//viper.SetConfigType("yaml") // 如果配置文件的名称中没有扩展名，则需要配置此项

	// 4. 环境变量配置文件查找的路径，相对于 main.go
	Viper.AddConfigPath(".")

	// 5. 开始读根目录下的 .env 文件，读不到会报错
	err := Viper.ReadInConfig()
	logger.Error(err)

	// 6. 设置环境变量前缀，用以区分 Go 的系统环境变量
	// 例如原本环境变量是：APP_NAME，设置了前缀之后：APPENV_APP_NAME
	// 注意在测试环境变量的时候，需要在命令行中启动程序：go run main.go
	//Viper.SetEnvPrefix("appenv")

	// 7. Viper.Get() 时，优先读取环境变量
	Viper.AutomaticEnv()
}

// Env 读取环境变量，支持默认值
func Env(envName string, defaultValue ...interface{}) interface{} {
	if len(defaultValue) > 0 {
		return Get(envName, defaultValue[0])
	}
	return Get(envName)
}

// Add 新增配置项
// map 结构支持使用 . 的方式访问键
func Add(name string, configuration map[string]interface{}) {
	Viper.Set(name, configuration)
}

// Get 获取配置项，允许使用点式获取，如：app.name
func Get(path string, defaultValue ...interface{}) interface{} {
	// 不存在的情况
	if !Viper.IsSet(path) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	return Viper.Get(path)
}

// GetString 获取 String 类型的配置信息
func GetString(path string, defaultValue ...interface{}) string {
	return cast.ToString(Get(path, defaultValue...))
}

// GetInt 获取 Int 类型的配置信息
func GetInt(path string, defaultValue ...interface{}) int {
	return cast.ToInt(Get(path, defaultValue...))
}

// GetInt64 获取 Int64 类型的配置信息
func GetInt64(path string, defaultValue ...interface{}) int64 {
	return cast.ToInt64(Get(path, defaultValue...))
}

// GetUint 获取 Uint 类型的配置信息
func GetUint(path string, defaultValue ...interface{}) uint {
	return cast.ToUint(Get(path, defaultValue...))
}

// GetBool 获取 Bool 类型的配置信息
func GetBool(path string, defaultValue ...interface{}) bool {
	return cast.ToBool(Get(path, defaultValue...))
}
