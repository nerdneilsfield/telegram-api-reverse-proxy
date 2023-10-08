package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var (
	remoteURL = flag.String("url", "https://api.telegram.org", "remote api url")
)

func main() {

	flag.Parse()

	fmt.Println("Remote URL: ", *remoteURL)
	fmt.Println("Starting the server on port 18080...")

	http.HandleFunc("/", handler)
	http.ListenAndServe("0.0.0.0:18080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	targetURL, _ := url.Parse(*remoteURL)

	// 设置允许跨域访问的响应头
	origin := r.Header.Get("Access-Control-Allow-Origin")
	if origin == "" {
		origin = "*"
	}
	w.Header().Set("Access-Control-Allow-Origin", origin)

	// 读取请求体
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// 创建一个新的请求
	newRequest, err := http.NewRequest(r.Method, targetURL.String()+r.RequestURI, io.ReadCloser(io.NopCloser(bytes.NewReader(body))))
	if err != nil {
		http.Error(w, "Error creating new request", http.StatusInternalServerError)
		return
	}

	// 复制原始请求的头
	for key, value := range r.Header {
		newRequest.Header[key] = value
	}

	// 使用 HTTP 客户端发送请求并获取响应
	client := &http.Client{}
	response, err := client.Do(newRequest)
	if err != nil {
		errormsg := fmt.Sprintf("Error sending request to Reverse API %v: %v", *remoteURL, err)
		http.Error(w, errormsg, http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	// 复制响应体到响应 writer
	io.Copy(w, response.Body)
}
