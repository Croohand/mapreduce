package commands

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/Croohand/mapreduce/common/blockutil"
	"github.com/Croohand/mapreduce/common/httputil"
)

func Exists(path string) bool {
	if !blockutil.ValidateFilePath(path) {
		log.Fatal("invalid file path " + path)
	}
	resp, err := http.PostForm(mrConfig.Host+"/File/IsExists", url.Values{"Path": []string{path}})
	if err != nil {
		log.Fatal(err)
	}
	var exists struct{ Exists bool }
	if err := httputil.GetJson(resp, &exists); err != nil {
		log.Fatal(err)
	}
	if exists.Exists {
		fmt.Printf("File %s exists\n", path)
	} else {
		fmt.Printf("File %s doesn't exist\n", path)
	}
	return exists.Exists
}
