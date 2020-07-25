package main

import (
	"fmt"
	"os"
)

const filename = `generated.json`

func main(){
	jsonFile, err := os.OpenFile(filename, os.O_RDONLY, 0600)
	defer logAndClose(jsonFile)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", jsonFile)
}


func logAndClose(file *os.File) {
	err := file.Close()
	if err != nil {
		panic(err)
	}
}
