package main

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"github.com/NYTimes/gziphandler"
	"github.com/vmihailenco/msgpack/v5"
	"io/ioutil"
	"net/http"
	"os"
)

const filename = `generated.json`

func main(){
	withoutGz := http.HandlerFunc(respond)
	withGz := gziphandler.GzipHandler(withoutGz)

	withMsgPack := http.HandlerFunc(respondWithMsgPack)
	withMsgPackgz := gziphandler.GzipHandler(withMsgPack)

	withZlib := http.HandlerFunc(respondWithZlib)
	withZlibGzip := gziphandler.GzipHandler(withZlib)

	http.Handle("/hello", withoutGz)
	http.Handle("/hellogz", withGz)

	http.Handle("/msg", withMsgPack)
	http.Handle("/msggz", withMsgPackgz)

	http.Handle("/zlib", withZlib)
	http.Handle("/zlibgz", withZlibGzip)


	http.ListenAndServe(":8090", nil)
}


func logAndClose(file *os.File) {
	err := file.Close()
	if err != nil {
		panic(err)
	}
}

func respond(w http.ResponseWriter, req *http.Request) {
	fc := fileContents()

	fmt.Fprintf(w, string(fc))
}

type Item struct {
	Contents string
}

func respondWithMsgPack(w http.ResponseWriter, req *http.Request) {
	fc := fileContents()

	var arbitraryUnmarshaledJson interface{}

	err := json.Unmarshal(fc, &arbitraryUnmarshaledJson)
	if err != nil {
		panic(err)
	}

	b, err := msgpack.Marshal(arbitraryUnmarshaledJson)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition", "inline")

	fmt.Fprintf(w, "%s", b)
}

func respondWithZlib(w http.ResponseWriter, req *http.Request) {
	fc := fileContents()

	var b bytes.Buffer
	zw := zlib.NewWriter(&b)
	zw.Write(fc)
	zw.Close()

	fmt.Fprintf(w, b.String())
}



func fileContents() []byte {
	jsonFile, err := os.OpenFile(filename, os.O_RDONLY, 0600)
	defer logAndClose(jsonFile)

	if err != nil {
		panic(err)
	}

	fileContents, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}
	return fileContents
}
