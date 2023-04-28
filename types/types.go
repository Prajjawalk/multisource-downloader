package types

type UrlInfo struct {
	Url      string
	FileSize int64
}

type DownloadErrorRecord struct {
	UrlIndex   int
	StartChunk int64
	EndChunk   int64
	ErrorMsg   string
}
