package handler

import (
	rPool "FILE-SERVER/cache/redis"
	dblayer "FILE-SERVER/db"
	"FILE-SERVER/util"
	"fmt"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
)

// 分开的信息
type MultiPartUploadInfo struct {
	FileHash   string
	FileSize   int
	UploadID   string
	ChunkSize  int
	ChunkCount int
}

// 初始化分开上传
func InitialMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析用户请求信息
	r.ParseForm()
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize, err := strconv.Atoi(r.Form.Get("filesize"))
	if err != nil {
		w.Write(util.NewRespMsg(-1, "params invalid", nil).JSONBytes())
		return
	}

	// 2. 获取redis的一个连接
	rCoon := rPool.RedisPool().Get()
	defer rCoon.Close()

	// 3. 生成分块上传的初始化信息
	upinfo := MultiPartUploadInfo{
		FileHash:   filehash,
		FileSize:   filesize,
		UploadID:   username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:  5 * 1024 * 1024,
		ChunkCount: int(math.Ceil(float64(filesize) / (5 * 1024 * 1024))),
	}

	// 4. 将初始化信息写入到redis缓存
	rCoon.Do("HSET", "MP_"+upinfo.UploadID, "chunkcount", upinfo.ChunkCount)
	rCoon.Do("HSET", "MP_"+upinfo.UploadID, "filehash", upinfo.FileHash)
	rCoon.Do("HSET", "MP_"+upinfo.UploadID, "filesize", upinfo.FileSize)

	// 5. 将初始化的信息返回给客户端
	w.Write(util.NewRespMsg(0, "OK", upinfo).JSONBytes())
}

func UploadPartHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析用户请求参数
	r.ParseForm()
	// username := r.Form.Get("usename")
	uploadID := r.Form.Get("uploadid")
	chunkIndex := r.Form.Get("index")
	// 2. 获取redis连接池的一个连接
	rCoon := rPool.RedisPool().Get()
	defer rCoon.Close()

	// 3. 获取文件句柄，用于存储分块内容
	fpath := "/data/" + uploadID + "/" + chunkIndex
	os.MkdirAll(path.Dir(fpath), 0744)
	fd, err := os.Create(fpath)
	if err != nil {
		w.Write(util.NewRespMsg(-1, "Upload part failed", nil).JSONBytes())
		return
	}
	defer fd.Close()
	buf := make([]byte, 1024*1024)
	for {
		n, err := r.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil {
			break
		}
	}

	// 4. 更新redis缓存状态
	rCoon.Do("HSET", "MP_"+uploadID, "chkidx_"+chunkIndex, 1)

	// 5. 返回处理结果给客户端
	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}

func CompleteUploadHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析请求参数
	r.ParseForm()
	upid := r.Form.Get("uploadid")
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize := r.Form.Get("filesize")
	filename := r.Form.Get("filename")

	// 2. 获得redis连接池的一个连接
	rCoon := rPool.RedisPool().Get()
	defer rCoon.Close()

	// 3. 通过uploadid查询redis并判断是否所有分块上传完成
	data, err := redis.Values(rCoon.Do("HGETALL", "MP_"+upid))
	if err != nil {
		w.Write(util.NewRespMsg(-1, "complete upload failed", nil).JSONBytes())
		return
	}
	totalCount := 0
	chunkCount := 0
	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))
		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkCount++
		}
	}
	if totalCount != chunkCount {
		w.Write(util.NewRespMsg(-1, "invalid request", nil).JSONBytes())
		return
	}

	// 4. TODO : 合并分块
	// 5. 更新唯一文件表及用户文件夹
	fsize, _ := strconv.Atoi(filesize)

	dblayer.OnFileUploadFinished(filehash, filename, int64(fsize), "")
	dblayer.OnUserFileUploadFinished(username, filehash, filename, int64(fsize))
	// 6. 响应处理结果
	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}

func CancelUploadParHandler(w http.ResponseWriter, r *http.Request) {
	// 删除已存在的分块文件
	// 删除redis缓存状态
	// 更新mysql文件status
}

func MultiPartUploadStatusHandler(w http.ResponseWriter, r *http.Request) {
	// 检查分块上传状态是否有效

	// 获取分块初始化信息

	// 获取已上传的分块信息

}
