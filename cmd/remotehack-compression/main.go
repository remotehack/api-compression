package main

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/remotehack/api-compression/v2/pkg/noise"
	"github.com/vmihailenco/msgpack/v5"
)

func main() {

	respondWithNoise()
	return
	respondWithSound()

	withoutGz := http.HandlerFunc(respond)
	withGz := gziphandler.GzipHandler(withoutGz)

	withMsgPack := http.HandlerFunc(respondWithMsgPack)
	withMsgPackgz := gziphandler.GzipHandler(withMsgPack)

	withZlib := http.HandlerFunc(respondWithZlib)
	withZlibGzip := gziphandler.GzipHandler(withZlib)

	//withSound := http.HandlerFunc(respondWithSound)

	http.Handle("/hello", withoutGz)
	http.Handle("/hellogz", withGz)

	http.Handle("/msg", withMsgPack)
	http.Handle("/msggz", withMsgPackgz)

	http.Handle("/zlib", withZlib)
	http.Handle("/zlibgz", withZlibGzip)

	//http.Handle("/sound", withSound)

	http.ListenAndServe(":8090", nil)
}

func logAndClose(file *os.File) {
	err := file.Close()
	if err != nil {
		panic(err)
	}
}

func respond(w http.ResponseWriter, req *http.Request) {
	fc := fileContents("generated.json")

	fmt.Fprintf(w, string(fc))
}

type Item struct {
	Contents string
}

func respondWithMsgPack(w http.ResponseWriter, req *http.Request) {
	fc := fileContents("generated.json")

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
	fc := fileContents("generated.json")

	var b bytes.Buffer
	zw := zlib.NewWriter(&b)
	zw.Write(fc)
	zw.Close()

	fmt.Fprintf(w, b.String())
}

func respondWithSound() {
	f, err := os.Open("Honk.mp3")
	if err != nil {
		log.Fatal("failed to open Honk")
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()
	done := make(chan bool)
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
}

func fileContents(filename string) []byte {
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

func respondWithNoise() {
	fc := fileContents("generated.json")
	ns := noise.New(fc)

	fmt.Println("noise filled\n\n\n")

	sr := beep.SampleRate(44100)
	speaker.Init(sr, sr.N(time.Second/10))
	speaker.Play(ns)
	select {}
}
