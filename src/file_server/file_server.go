package file_server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func ManageFile(w http.ResponseWriter, req *http.Request) {
	url := req.URL.String()[1:]
	list := false
	delete := false
	url_len := len(url)
	del_len := len("/delete")
	add_len := len("/add")
	list_len := len("/list")

	if url[url_len-add_len:] == "/add" {
		// Open a new file for writing only
		url = url[:url_len-add_len]
		file, err := os.OpenFile(
			url,
			os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
			0666,
		)
		if err != nil {
			//log.Fatal(err)
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		defer file.Close()

		// Write bytes to file
		byteSlice, ioerr := ioutil.ReadAll(req.Body)
		if ioerr != nil {
			//log.Fatal(err)
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		bytesWritten, err := file.Write(byteSlice)
		if err != nil {
			//log.Fatal(err)
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		log.Printf("Wrote %d bytes.\n", bytesWritten)
	} else if url[url_len-del_len:] == "/delete" {
		delete = true
		url = url[:url_len-del_len]
		err := os.Remove(url)
		if err != nil {
			//log.Fatal(err)
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			fmt.Fprintf(w, `{"status": "success"}`)
		}
	} else if url[url_len-list_len:] == "/list" {
		// walk all files in directory
		url = url[:url_len-list_len]
		log.Printf(url)
		list = true
		filepath.Walk(url, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				fmt.Fprintln(w, info.Name())
			}
			return nil
		})
	}

	if !delete && !list {
		// always return file if above succeds or is irelavent
		b, err := ioutil.ReadFile(url) // just pass the file name
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		fmt.Fprintf(w, string(b))
	}
}
