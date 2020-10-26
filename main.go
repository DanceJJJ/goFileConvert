package main

import (
	"flag"
	"io/ioutil"
	"julian/goFileConvert/convert"
	config2 "julian/goFileConvert/tools/config"
	"julian/goFileConvert/utils"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var (
	config string
	extensionToContentType = map[string]string {
		".html": "text/html; charset=utf-8",
		".css": "text/css; charset=utf-8",
		".js": "application/javascript",
		".xml": "text/xml; charset=utf-8",
		".jpg":  "image/jpeg",
	}
)

func main() {

	//获取命令行参数
	flag.StringVar(&config, "c", "config/settings.yml", "Config file")
	flag.Parse()

	//获取配置信息
	config2.ConfigSetup(config)

	//文档转化接口
	http.HandleFunc("/convert", convert.Handle)

	//读取静态文件
	http.HandleFunc("/cache/", fileHandler)

	//服务启动
	_ = http.ListenAndServe(":" + config2.ApplicationConfig.Port, nil)
}

func fileHandler(w http.ResponseWriter, r *http.Request) {

	path := "." + r.URL.Path

	f, err := os.Open(path)
	if err != nil {
		utils.Error(w, utils.ToHTTPError(err))
		return
	}
	defer f.Close()

	d, err := f.Stat()
	if err != nil {
		utils.Error(w, utils.ToHTTPError(err))
		return
	}

	if d.IsDir() {
		utils.DirList(w, r, f)
		return
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		utils.Error(w, utils.ToHTTPError(err))
		return
	}

	ext := filepath.Ext(path)
	if contentType := extensionToContentType[ext]; contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}

	w.Header().Set("Content-Length", strconv.FormatInt(d.Size(), 10))
	w.Write(data)
}
