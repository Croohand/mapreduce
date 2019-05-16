package commands

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/responses"
)

func List(prefix string) {
	resp, err := http.PostForm(mrConfig.Host+"/File/List", url.Values{"Prefix": {prefix}})
	if err != nil {
		log.Panic(err)
	}
	var files responses.ListedFiles
	if err := httputil.GetJson(resp, &files); err != nil {
		log.Panic(err)
	}
	fmt.Println(strings.Join(files, "\n"))
}
