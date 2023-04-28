package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/Prajjawalk/multisource-downloader/types"
	"github.com/Prajjawalk/multisource-downloader/utils"
	"github.com/Prajjawalk/multisource-downloader/worker"
)

func main() {
	var (
		urlList string
		output  string
	)

	//usage example
	//go run ./main.go --output=output.tar.gz --urls https://filebin.net/vun7wywrbksi7sw8/eth2-beaconchain-explorer-1.19.3.tar.gz https://filebin.net/vun7wywrbksi7sw8/eth2-beaconchain-explorer-1.19.3__1_.tar.gz
	flag.StringVar(&urlList, "urls", "", "list of urls to download file from (--urls url1 url2 url3...)")
	flag.StringVar(&output, "output", "", "output file name")

	flag.Parse()

	urls := strings.Split(urlList, " ")
	filename := output
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()
	chunkSize := int64(0)
	urlInfoArray := []types.UrlInfo{}
	for _, url := range urls {
		filesize, err := utils.ConfirmSupportAndFileSize(url)
		if err != nil {
			fmt.Println("Error getting file size:", err, "from url", url)
			continue
		}

		if filesize <= 0 {
			fmt.Println("File size is unknown from url: ", url)
			continue
		}

		urlInfoArray = append(urlInfoArray, types.UrlInfo{
			Url:      url,
			FileSize: filesize,
		})
	}

	sort.Slice(urlInfoArray, func(i, j int) bool { return urlInfoArray[i].FileSize < urlInfoArray[j].FileSize })

	chunkSize = urlInfoArray[len(urlInfoArray)-1].FileSize / int64(len(urlInfoArray))
	errorRecord := make([]types.DownloadErrorRecord, 0)

	var wg sync.WaitGroup

	for i := 0; i < len(urlInfoArray); i++ {
		start := int64(i) * chunkSize
		end := start + chunkSize - 1
		if i == len(urlInfoArray)-1 {
			end = urlInfoArray[i].FileSize - 1 // last chunk may be larger
		}
		wg.Add(1)
		go worker.DownloadWorker(start, end, i, &errorRecord, urlInfoArray, &wg, file)
	}

	if len(errorRecord) != 0 {
		for _, rec := range errorRecord {
			index := rec.UrlIndex
			retryIndex := index + 1
			if index == len(errorRecord)-1 {
				retryIndex = index - 1
			}
			wg.Add(1)
			fmt.Printf("Retrying download from url %v...", urlInfoArray[rec.UrlIndex].Url)
			go worker.RetryWorker(rec.StartChunk, rec.EndChunk, retryIndex, urlInfoArray, &wg, file)
		}
	}

	wg.Wait()

	fmt.Println("Download complete")
}
