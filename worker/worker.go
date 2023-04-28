package worker

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/Prajjawalk/multisource-downloader/types"
)

func DownloadWorker(start, end int64, idx int, record *[]types.DownloadErrorRecord, urlInfoArray []types.UrlInfo, wg *sync.WaitGroup, file *os.File) {
	defer wg.Done()
	req, err := http.NewRequest("GET", urlInfoArray[idx].Url, nil)
	if err != nil {
		fmt.Printf("Error creating request for range %d-%d: %s\n", start, end, err)
		*record = append(*record, types.DownloadErrorRecord{
			UrlIndex:   idx,
			StartChunk: start,
			EndChunk:   end,
			ErrorMsg:   fmt.Sprintf("Error creating request for range %d-%d: %s\n", start, end, err),
		})
		return
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error downloading range %d-%d: %s\n", start, end, err)
		*record = append(*record, types.DownloadErrorRecord{
			UrlIndex:   idx,
			StartChunk: start,
			EndChunk:   end,
			ErrorMsg:   fmt.Sprintf("Error downloading range %d-%d: %s\n", start, end, err),
		})
		return
	}
	defer resp.Body.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Printf("Error writing range %d-%d to file: %s\n", start, end, err)
		*record = append(*record, types.DownloadErrorRecord{
			UrlIndex:   idx,
			StartChunk: start,
			EndChunk:   end,
			ErrorMsg:   fmt.Sprintf("Error writing range %d-%d to file: %s\n", start, end, err),
		})
		return
	}
}

func RetryWorker(start, end int64, idx int, urlInfoArray []types.UrlInfo, wg *sync.WaitGroup, file *os.File) {
	defer wg.Done()
	req, err := http.NewRequest("GET", urlInfoArray[idx].Url, nil)
	if err != nil {
		fmt.Printf("Retry failed: Error creating request for range %d-%d: %s\n", start, end, err)
		return
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Retry failed: Error downloading range %d-%d: %s\n", start, end, err)
		return
	}
	defer resp.Body.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Printf("Retry failed: Error writing range %d-%d to file: %s\n", start, end, err)
		return
	}
}
