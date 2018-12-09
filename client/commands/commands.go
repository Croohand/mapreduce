package commands

import (
	"log"
	"net/http"
	"os"

	"github.com/Croohand/mapreduce/common/httputil"
)

var mrConfig struct {
	BlockSize            int
	ReplicationFactor    int
	MinReplicationFactor int
	Host                 string
}

func Init() {
	mrConfig.Host = os.Getenv("MR_HOST")
	resp, err := http.Get(mrConfig.Host + "/GetMRConfig")
	if err != nil {
		log.Fatal("couldn't get MR config from master, error: " + err.Error())
	}
	if err := httputil.GetJson(resp, &mrConfig); err != nil {
		log.Fatal("couldn't get MR config from master, error: " + err.Error())
	}
}
