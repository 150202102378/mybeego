package ini

import (
	"os"
	u "utxo/utils"

	"github.com/BurntSushi/toml"
)

//TomlConfig is struct of config.toml
type TomlConfig struct {
	RunMode string
	BTCCore map[string]BTCCoreConf
	Psql    map[string]PsqlConf
	Email   EmailConf
}

func initConfig() {
	configFilePath := u.GetProjectPath("utxo") + "ini/resources/config.toml"
	_, err := toml.DecodeFile(configFilePath, &config)
	if err != nil {
		log.Fatal("Init Config File Error")
		os.Exit(3)
	}
}

func initRunMode() {
	if len(os.Args) >= 2 {
		runMode = os.Args[1]
	} else {
		runMode = config.RunMode
	}
}
