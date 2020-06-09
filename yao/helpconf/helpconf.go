package helpconf

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

func ViperReadInConfig(fileName, DirName string) error {
	viper.SetConfigName(fileName)
	viper.AddConfigPath(DirName)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Print("Config file changed:", e.Name)
	})
	return nil
}
