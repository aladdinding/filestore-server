package handler

import (
	rpool "filestore-server/cache/redis"
	"filestore-server/util"
	"fmt"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

// MultipartUploadInfo:初始化信息
type MultipartUploadInfo struct {
	FileHash   string
	FileSize   int
	UploadID   string
	ChunkSize  int
	ChunkCount int
}

func InitialMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析用户请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize, err := strconv.Atoi(r.Form.Get("filesize"))
	if err != nil {
		w.Write(util.NewRespMsg(-1, "param invalid", nil).JSONBytes())
		return
	}

	// 2. 获得redis的一个连接
	rConn := rpool.RedisPool().Get()
	defer rConn.Close()

	// 3. 生成分块上传的初始化信息
	upInfo := MultipartUploadInfo{
		FileHash:   filehash,
		FileSize:   filesize,
		UploadID:   username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:  5 * 1024 * 1024,
		ChunkCount: int(math.Ceil(float64(filesize) / (5 * 1024 * 1024))),
	}

	//4. 将初始化信息
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "chunkcount", upInfo.ChunkCount)
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "filehash", upInfo.FileHash)
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "filesize", upInfo.FileSize)

	w.Write(util.NewRespMsg(0, "OK", upInfo).JSONBytes())

}

func UploadPartHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析用户请求参数
	r.ParseForm()
	uploadID := r.Form.Get("upload")
	chunkIndex := r.Form.Get("index")

	// 2. 获得redis连接池重的一个连接
	rConn := rpool.RedisPool().Get()
	defer rConn.Close()

	// 3. 获取文件句柄，用于存储分块的内容
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
	rConn.Do("HSET", "MP_"+uploadID, "chkidx_"+chunkIndex, 1)

	//5. 返回处理结果到客户端
	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}
