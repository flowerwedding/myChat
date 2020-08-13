/**
 * @Title  CONFIG
 * @description  #
 * @Author  沈来
 * @Update  2020/8/13 15:27
 **/
package global

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	SensitiveWords []string
)

func initConfig(){
	viper.SetConfigName("chatroom")
	viper.AddConfigPath(RootDir + "/config")

	if  err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	SensitiveWords = viper.GetStringSlice("sensitive")

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event){
		_ = viper.ReadInConfig()

		SensitiveWords = viper.GetStringSlice("sensitive")
	})
}