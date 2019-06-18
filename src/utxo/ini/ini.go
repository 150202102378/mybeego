package ini

import (
	"net/http"
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

var (
	config     TomlConfig
	runMode    string
	log        = logrus.New()
	psqlDB     *gorm.DB
	psqlDBLock sync.Mutex
	httpClient *http.Client
)

func init() {
	initConfig()
	initRunMode()
	initLog()
}
