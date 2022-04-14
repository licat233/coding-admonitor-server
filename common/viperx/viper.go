package viperx

import (
	"admonitor-server/common/utils"
	"fmt"
	"path"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	configPath := path.Join(utils.CurrentProjectRootPath(), "deploy/server/")
	viper.AddConfigPath(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	viper.WatchConfig() //程序在运行中实时读取配置文件
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
	})
}
