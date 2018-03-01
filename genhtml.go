package main

import (
	"fmt"
	"os"
	"strings"
	"io/ioutil"
	"os/exec"
	"encoding/json"
	"log"
	"net/http"
	"flag"
	"path/filepath"
)

func listDir(dirPth string) (files []string, names []string, err error) {
	suffixList := []string{".mp4", ".avi", ".wmv", ".mkv", ".mov", ".flv"}

	files = make([]string, 0, 10)
	names = make([]string, 0, 10)

	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, nil, err
	}

	pthSep := string(os.PathSeparator)
	for _, fi := range dir {
		if fi.IsDir() {
			// 忽略目录
			continue
		}
		for _, suffix := range suffixList {
			suffix = strings.ToUpper(suffix)
			if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
				//匹配文件
				files = append(files, dirPth + pthSep + fi.Name())
				names = append(names, fi.Name())
			}
		}

	}

	return files, names, nil
}

func getCurrentPath() string {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		fmt.Printf("exec.LookPath(%s), err: %s\n", os.Args[0], err)
		return ""
	}
	path, err := filepath.Abs(file)
	if err != nil {
		fmt.Printf("filepath.Abs(%s), err: %s\n", file, err)
		return ""
	}

	dir := filepath.Dir(path)
	return dir
}

type TraceHandler struct {
	h http.Handler
}

func genHtmlFile() {
	path := getCurrentPath()
	_, names, err := listDir(path)

	jsonData, err := json.Marshal(names)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(jsonData))
	fmt.Println(path)

	templateHtml, err := ioutil.ReadFile("template.html")
	html := strings.Replace(string(templateHtml), "[[NAMES]]", string(jsonData), -1)
	f, err := os.Create("index.html")
	f.WriteString(html)
	f.Close()
}

const RET_CODE_OK = 100;

func handCmdRefresh(w http.ResponseWriter, r *http.Request) {
	fmt.Println("start refresh")
	genHtmlFile()
	w.Write([]byte{RET_CODE_OK})
}

func (r TraceHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	println("get", req.URL.Path, " from ", req.RemoteAddr)

	r.h.ServeHTTP(w, req)
}

func main() {
	fmt.Println("genhtml")
	path := getCurrentPath()

	genHtmlFile()

	isServer := flag.Bool("server", false, "Run as http server")
	port := flag.String("port", "8083", "Server listen port")
	flag.Parse()
	if *isServer {
		fmt.Println("start server")
		h := http.FileServer(http.Dir(path))
		http.Handle("/", TraceHandler{h})
		http.HandleFunc("/cmd-refresh", handCmdRefresh)
		fmt.Println("Listening on port ", *port, "...")
		log.Fatal("ListenAndServe: ", http.ListenAndServe(":" + *port, nil))
	}

}
