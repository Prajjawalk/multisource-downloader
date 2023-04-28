package worker

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/Prajjawalk/multisource-downloader/types"
	"github.com/stretchr/testify/assert"
)

func TestDownloadWorker(t *testing.T) {
	testcases := []struct {
		name   string
		start  int64
		end    int64
		idx    int
		record *[]types.DownloadErrorRecord
	}{
		{
			name:   "test ok",
			start:  0,
			end:    10,
			idx:    0,
			record: &[]types.DownloadErrorRecord{},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Accept-Ranges", fmt.Sprintf("bytes=%d-%d", tc.start, tc.end))
				w.Header().Add("Content-Length", "10")
				w.Write([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
			}))
			defer ts.Close()
			var wg sync.WaitGroup
			wg.Add(1)
			urlInfoArray := []types.UrlInfo{
				{
					Url:      ts.URL,
					FileSize: 10,
				},
			}
			file, err := os.Create("output.txt")
			defer os.Remove("output.txt")
			if err != nil {
				t.Error("Error creating file:", err)
			}
			DownloadWorker(tc.start, tc.end, tc.idx, tc.record, urlInfoArray, &wg, file)
			stats, _ := file.Stat()
			assert.Equal(t, stats.Size(), int64(10))
			assert.Equal(t, len(*tc.record), 0)
		})
	}
}

func TestRetryWorker(t *testing.T) {
	testcases := []struct {
		name  string
		start int64
		end   int64
		idx   int
	}{
		{
			name:  "test ok",
			start: 0,
			end:   10,
			idx:   0,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Accept-Ranges", fmt.Sprintf("bytes=%d-%d", tc.start, tc.end))
				w.Header().Add("Content-Length", "10")
				w.Write([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
			}))
			defer ts.Close()
			var wg sync.WaitGroup
			wg.Add(1)
			urlInfoArray := []types.UrlInfo{
				{
					Url:      ts.URL,
					FileSize: 10,
				},
			}
			file, err := os.Create("output.txt")
			defer os.Remove("output.txt")
			if err != nil {
				t.Error("Error creating file:", err)
			}
			RetryWorker(tc.start, tc.end, tc.idx, urlInfoArray, &wg, file)
			stats, _ := file.Stat()
			assert.Equal(t, stats.Size(), int64(10))
		})
	}
}
