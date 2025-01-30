package main

import (
	"FILE-SERVER/handle"
	"fmt"
	"net/http"
)

func main() {
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/file/upload", handle.UploadHandler)
	http.HandleFunc("/file/upload/suc", handle.UploadSucHandler)
	http.HandleFunc("/file/meta", handle.GetFileMetaHandler)
	http.HandleFunc("/file/query", handle.FileQueryHandler)
	http.HandleFunc("/file/download", handle.DownloadHandler)
	http.HandleFunc("/file/update", handle.FileMetaUpdateHandler)
	http.HandleFunc("/file/delete", handle.FileDeleteHandler)
	// 秒传接口
	http.HandleFunc("/file/fastupload", handle.HTTPInterceptor(
		handle.TryFastUploadHandler))

	http.HandleFunc("/file/downloadurl", handle.HTTPInterceptor(
		handle.DownloadURLHandler))

	// 分块上传接口
	http.HandleFunc("/file/mpupload/init",
		handle.HTTPInterceptor(handle.InitialMultipartUploadHandler))
	http.HandleFunc("/file/mpupload/uppart",
		handle.HTTPInterceptor(handle.UploadPartHandler))
	http.HandleFunc("/file/mpupload/complete",
		handle.HTTPInterceptor(handle.CompleteUploadHandler))

	http.HandleFunc("/user/signup", handle.SignupHandler)
	http.HandleFunc("/user/signin", handle.SigninHandler)
	http.HandleFunc("/user/info", handle.HTTPInterceptor(handle.UserInfoHandler))
	err := http.ListenAndServe(":8082", nil)

	if err != nil {
		fmt.Printf("Failed to start server, err: %s", err)
	}
}
