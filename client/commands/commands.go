package commands

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/responses"
)

type MrConfig struct {
	responses.MrConfig
	Hosts []string
}

func (cfg MrConfig) GetHost() string {
	for _, host := range cfg.Hosts {
		resp, err := http.Get(host + "/IsAlive")
		if err != nil {
			continue
		}
		var status responses.MasterStatus
		err = httputil.GetJson(resp, &status)
		if err != nil {
			continue
		}
		if status.State == "active" {
			return host
		}
	}
	log.Panic("No master available")
	return ""
}

var mrConfig MrConfig

func Init() {
	mrConfig.Hosts = strings.Split(os.Getenv("MR_HOSTS"), ",")
	resp, err := http.Get(mrConfig.GetHost() + "/GetMrConfig")
	if err != nil {
		log.Panic("Couldn't get MR config from master, error: " + err.Error())
	}
	if err := httputil.GetJson(resp, &mrConfig); err != nil {
		log.Panic("Couldn't get MR config from master, error: " + err.Error())
	}
}
