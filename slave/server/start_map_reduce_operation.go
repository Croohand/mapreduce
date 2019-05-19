package server

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/responses"
)

type reducerResult struct {
	reducer string
	result  responses.PathBlocks
}

func doTasks(blocks responses.PathBlocks, txId string, mappers int, reducersAddrs []string) (responses.PathBlocks, error) {
	reducers := len(reducersAddrs)

	hasSources := map[string]bool{}
	ensureSources := func(where string) error {
		if !hasSources[where] {
			srcsPath := fsutil.GetSourcesDir(txId)
			err := httputil.WriteSources(srcsPath, txId, where)
			if err != nil {
				return err
			}
			resp, err := http.PostForm(where+"/Source/Build", url.Values{"TransactionId": {txId}})
			if err != nil {
				return err
			}
			err = httputil.GetError(resp)
			if err != nil {
				return err
			}
			hasSources[where] = true
		}
		return nil
	}

	mapTask := func(block fsutil.BlockInfoEx, pool <-chan bool, errs chan<- error, done chan<- bool) {
		var err error
		for _, slave := range block.Slaves {
			err = nil
			var resp *http.Response
			err = ensureSources(slave)
			if err != nil {
				continue
			}

			resp, err = http.PostForm(slave+"/Operation/Map", url.Values{"BlockId": {block.Id}, "TransactionId": {txId}, "Reducers": {strconv.Itoa(reducers)}})
			if err != nil {
				continue
			}
			if err = httputil.GetError(resp); err != nil {
				continue
			}

			where := []string{}
			for i, addr := range reducersAddrs {
				where = append(where, strconv.Itoa(i)+" "+addr)
			}
			resp, err = http.PostForm(slave+"/Operation/SendResults", url.Values{"BlockId": {block.Id}, "TransactionId": {txId}, "Where": where})
			if err != nil {
				continue
			}
			if err = httputil.GetError(resp); err != nil {
				continue
			}
		}

		if err != nil {
			errs <- errors.New("Couldn't finish map task for block " + block.Id + ": " + err.Error())
			return
		}
		<-pool
		done <- true
	}

	pool := make(chan bool, mappers)
	errs := make(chan error)
	done := make(chan bool)
	for _, block := range blocks {
		select {
		case pool <- true:
		case err := <-errs:
			return nil, err
		}
		go mapTask(block, pool, errs, done)
	}

	for _ = range blocks {
		select {
		case <-done:
		case err := <-errs:
			return nil, err
		}
	}

	reduceTask := func(reducer string, errs chan<- error, res chan<- reducerResult) {
		err := ensureSources(reducer)
		if err != nil {
			errs <- err
			return
		}
		resp, err := http.PostForm(reducer+"/Operation/Reduce", url.Values{"TransactionId": {txId}})
		if err != nil {
			errs <- err
			return
		}
		var blocks responses.PathBlocks
		if err := httputil.GetJson(resp, &blocks); err != nil {
			errs <- err
			return
		}
		res <- reducerResult{reducer, blocks}
	}

	res := make(chan reducerResult)
	for _, reducer := range reducersAddrs {
		select {
		case err := <-errs:
			return nil, err
		default:
		}
		go reduceTask(reducer, errs, res)
	}

	gotBlocks := map[string]responses.PathBlocks{}
	for _ = range reducersAddrs {
		select {
		case r := <-res:
			gotBlocks[r.reducer] = r.result
		case err := <-errs:
			return nil, err
		}
	}

	blocks = responses.PathBlocks{}
	for _, reducer := range reducersAddrs {
		blocks = append(blocks, gotBlocks[reducer]...)
	}
	return blocks, nil
}

func startMapReduceOperation(in []string, out, readTxId, txId string, mappers, reducers int) {
	defer func() {
		if e := recover(); e != nil {
			log.Print(e)
		}
	}()
	readTxHandler := httputil.NewTxHandler()
	writeTxHandler := httputil.NewTxHandler()
	defer readTxHandler.Close()
	defer writeTxHandler.Close()
	go httputil.PingTransaction(Config.MasterAddr, readTxId, readTxHandler)
	go httputil.PingTransaction(Config.MasterAddr, txId, writeTxHandler)

	resp, err := http.PostForm(Config.MasterAddr+"/GetAvailableReducers", url.Values{"Number": {strconv.Itoa(reducers)}})
	if err != nil {
		panic(err)
	}
	reducersAddrs := []string{}
	if err := httputil.GetJson(resp, &reducersAddrs); err != nil {
		panic(err)
	}

	var blocks responses.PathBlocks
	for _, path := range in {
		resp, err = http.PostForm(Config.MasterAddr+"/File/Read", url.Values{"Path": {path}})
		if err != nil {
			panic(err)
		}
		var curBlocks responses.PathBlocks
		if err := httputil.GetJson(resp, &curBlocks); err != nil {
			panic(err)
		}
		blocks = append(blocks, curBlocks...)
	}

	outBlocks, err := doTasks(blocks, txId, mappers, reducersAddrs)

	if err != nil {
		panic(err)
	}

	err = httputil.TryWritePath(Config.MasterAddr, txId, out, outBlocks, false)

	if err != nil {
		panic(err)
	}

	err = httputil.TryValidateBlocks(Config.MasterAddr, txId, outBlocks)
	httputil.CleanUp(txId, outBlocks)
	removeTransaction(txId)

	if err != nil {
		panic(err)
	}
}
