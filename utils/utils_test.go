package utils

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfirmSupportAndFileSize(t *testing.T) {
	testcases := []struct {
		name               string
		expectedRangeStart int
		expectedRangeStop  int
		expectedFileSize   int
		errorMsg           error
		addAcceptRanges    bool
		addContentLength   bool
	}{
		{
			name:               "test ok",
			expectedRangeStart: 0,
			expectedRangeStop:  10,
			expectedFileSize:   10,
			errorMsg:           nil,
			addAcceptRanges:    true,
			addContentLength:   true,
		},
		{
			name:               "accept range header not present",
			expectedRangeStart: 0,
			expectedRangeStop:  10,
			expectedFileSize:   0,
			errorMsg:           errors.New("server error: Accept-Ranges Header does not exist in HTTP Response"),
			addAcceptRanges:    false,
			addContentLength:   true,
		},
		{
			name:               "content-length header not present",
			expectedRangeStart: 0,
			expectedRangeStop:  10,
			expectedFileSize:   0,
			errorMsg:           nil,
			addAcceptRanges:    true,
			addContentLength:   false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			expectedRangeStart := tc.expectedRangeStart
			expectedRangeStop := tc.expectedRangeStop
			expectedSize := int64(tc.expectedFileSize)

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tc.addAcceptRanges {
					w.Header().Add("Accept-Ranges", fmt.Sprintf("bytes=%d-%d", expectedRangeStart, expectedRangeStop))
				}

				if tc.addContentLength {
					w.Header().Add("Content-Length", "10")
				}
			}))
			defer ts.Close()

			filesize, err := ConfirmSupportAndFileSize(ts.URL)

			assert.Equal(t, expectedSize, filesize)
			assert.Equal(t, tc.errorMsg, err)
		})
	}
}
