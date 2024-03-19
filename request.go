package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const totalReqCount = Threads * PerThread

func getBody() []byte {
	body, err := json.Marshal(Body)
	if err != nil {
		log.Fatalln("Error: ", err)
	}

	return body
}

func getResponse(res *http.Response) ([]byte, error) {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func setHeaders(req *http.Request) {
	for key, value := range Headers {
		req.Header.Set(key, value)
	}
}

func makeRequest(client http.Client, body *[]byte, index int, successCount *int64) {
	req, err := http.NewRequest(Method, Url, bytes.NewBuffer(*body))
	setHeaders(req)

	if err != nil {
		log.Fatalln("ReqID err: ", index, " Error: ", err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println("ReqID send req err: ", index, " Error: ", err)
		return
	}

	reqRes, err := getResponse(res)
	if err != nil {
		log.Println("ReqID response err: ", index, "Error: ", err)
		return
	}
	reqString := string(reqRes)

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		atomic.AddInt64(successCount, 1)
	}

	var data map[string]any
	json.Unmarshal(reqRes, &data)
	log.Println("ReqID ", index, "Data: ", reqString)

	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			log.Fatalln("ReqID close: ", index, " Error: ", err)
			return
		}
	}(res.Body)
}

func run() {
	httpClient := http.Client{
		Timeout: 60 * time.Second,
	}

	var successCount int64 = 0
	var wg sync.WaitGroup
	var timeNow = time.Now()
	var body = getBody()

	for index := 0; index < Threads; index++ {
		wg.Add(1)
		go func(i int) {
			for x := 0; x < PerThread; x++ {
				makeRequest(httpClient, &body, (PerThread*i)+x, &successCount)
			}
			wg.Done()
		}(index)
	}

	wg.Wait()
	fmt.Println(
		"Done; Success: ", successCount,
		"; Failed: ", totalReqCount-successCount,
		"; Time: ", time.Since(timeNow),
	)
}
