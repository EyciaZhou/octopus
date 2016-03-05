package MOctopus

import (
	"database/sql"
	"github.com/EyciaZhou/configparser"
	log "github.com/Sirupsen/logrus"
	"fmt"
)

type Config struct {
	DBAddress  string `default:"127.0.0.1"`
	DBPort     string `default:"3306"`
	DBName     string `default:"version"`
	DBUsername string `default:"root"`
	DBPassword string `default:"fmttm233"`
}

var config Config
var (
	db *sql.DB
)

func init() {
	configparser.AutoLoadConfig("M.msghub", &config)

	var err error
	log.Info("M.msghub Start Connect mysql")
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.DBUsername, config.DBPassword, config.DBAddress, config.DBPort, config.DBName)
	db, err = sql.Open("mysql", url)
	if err != nil {
		log.Panic("M.msghub Can't Connect DB REASON : " + err.Error())
		return
	}
	err = db.Ping()
	if err != nil {
		log.Panic("M.msghub Can't Connect DB REASON : " + err.Error())
		return
	}
	log.Info("M.msghub connected")
}

