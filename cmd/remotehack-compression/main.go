package main

import (
	"fmt"
	"github.com/NYTimes/gziphandler"
	"io/ioutil"
	"net/http"
	"os"
)

const filename = `generated.json`

func main(){
	withoutGz := http.HandlerFunc(respond)
	withGz := gziphandler.GzipHandler(withoutGz)

	http.Handle("/hello", withoutGz)
	http.Handle("/hellogz", withGz)

	http.ListenAndServe(":8090", nil)
}


func logAndClose(file *os.File) {
	err := file.Close()
	if err != nil {
		panic(err)
	}
}

func respond(w http.ResponseWriter, req *http.Request) {
	jsonFile, err := os.OpenFile(filename, os.O_RDONLY, 0600)
	defer logAndClose(jsonFile)

	if err != nil {
		panic(err)
	}

	fileContents, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, string(fileContents))
}
