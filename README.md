# multisource-downloader

This is the multisource downloader written in golang. It takes multiple urls as input and downloads file in multiple parts from those urls. The file is divided into multiple chunks depending on number of urls and every url is alloted a chunk of file (starting bytes and ending bytes) which it needs to download. A new download worker is initiated and multiple download workers download file in parallel which results in faster download speeds.

### Requirements

* [Go 1.18+](https://golang.org/dl/)

### Tests

```
$ go test -v ./...
```

### Build

```
$ go build .
```

### Usage

```
$ multisource-downloader --output="output file name" --urls url1 url2 url3...
```
