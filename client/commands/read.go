package commands

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/url"

	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/responses"
)

func readInner(blocks responses.PathBlocks, res chan<- bytes.Buffer) {
	for _, block := range blocks {
		any := false
		for _, slave := range block.Slaves {
			if !httputil.IsSlaveAvailable(slave) {
				continue
			}
			resp, err := httpClient.PostForm(slave+"/Block/Read", url.Values{"BlockId": {block.Id}})
			if err != nil {
				log.Println(err)
				continue
			}
			if err := httputil.GetErrorNoClose(resp); err != nil {
				log.Println(err)
				continue
			}
			var b bytes.Buffer
			_, err = io.Copy(&b, resp.Body)
			if err != nil {
				log.Println(err)
				continue
			}
			res <- b
			any = true
			break
		}
		if !any {
			log.Panic("No available slaves for block " + block.Id)
		}
	}
	close(res)
}

func Read(path string) {
	if !existsInner(path) {
		log.Panic("File path " + path + " doesn't exist")
	}
	_, txHandler := startReadTransaction([]string{path})
	defer txHandler.Close()
	resp, err := httpClient.PostForm(mrConfig.GetHost()+"/File/Read", url.Values{"Path": {path}})
	if err != nil {
		log.Panic(err)
	}
	var blocks responses.PathBlocks
	if err := httputil.GetJson(resp, &blocks); err != nil {
		log.Panic(err)
	}
	res := make(chan bytes.Buffer, 2)
	go readInner(blocks, res)
	for buf := range res {
		fmt.Print(buf.String())
	}
}
