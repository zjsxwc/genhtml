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

func ListDir(dirPth string) (files []string, names []string,  err error) {
	suffixList := []string{".mp4", ".avi", ".wmv",".mkv",".mov",".flv"}

	files = make([]string, 0, 10)
	names = make([]string, 0, 10)

	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, nil, err
	}

	PthSep := string(os.PathSeparator)
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		for _, suffix := range suffixList{
			suffix = strings.ToUpper(suffix)
			if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
				files = append(files, dirPth+PthSep+fi.Name())
				names = append(names, fi.Name())
			}
		}

	}

	return files,  names,nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getCurrentPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	splitstring := strings.Split(path, "\\")
	size := len(splitstring)
	splitstring = strings.Split(path, splitstring[size-1])
	ret := strings.Replace(splitstring[0], "\\", "/", size-1)
	return ret
}


func main()  {
	fmt.Println("genhtml")
	path := getCurrentPath()
	_, names, err := ListDir(path)

	jsondata, err := json.Marshal(names)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(jsondata))
	fmt.Println(path)

	templateHtml,err := ioutil.ReadFile("template.html")
	html := strings.Replace(string(templateHtml), "[[NAMES]]", string(jsondata), -1)
	f, err := os.Create("index.html")
	f.WriteString(html)
	f.Close()

	isServer := flag.Bool("server", false, "Run as http server")
	port := flag.String("port", "8083", "Server listen port")
	flag.Parse()
	if *isServer {
		fmt.Println("start server")
		h := http.FileServer(http.Dir(path))
		http.Handle("/", TraceHandler{h})
		println("Listening on port ", *port,"...")
		log.Fatal("ListenAndServe: ", http.ListenAndServe(":"+ *port, nil))
	}

}
func (r TraceHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	println("get",req.URL.Path," from ",req.RemoteAddr)

	r.h.ServeHTTP(w, req)
}
type TraceHandler struct {
	h http.Handler
}
