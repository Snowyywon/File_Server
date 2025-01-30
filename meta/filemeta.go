package meta

import (
	mydb "FILE-SERVER/db"
	"sort"
)

// FileMeta : 文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// 更新fileMetas 新增/更新/文件元信息
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

// 更新fileMetas 到mysql中
func UpdateFileMetaDB(fmeta FileMeta) bool {
	return mydb.OnFileUploadFinished(fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}

// 获取fileMetas元信息结构
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

// 从mysql获取文件元信息
func GetFileMetaDB(fileSha1 string) (FileMeta, error) {
	tfile, err := mydb.GetFileMetaDB(fileSha1)
	if err != nil {
		return FileMeta{}, err
	}
	fmeta := FileMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize,
		Location: tfile.FileAddr.String,
	}
	return fmeta, nil
}

// 获取批量的文件元信息列表
func GetLastFileMeta(count int) []FileMeta {
	fMetaArray := make([]FileMeta, len(fileMetas))
	for _, v := range fileMetas {
		fMetaArray = append(fMetaArray, v)
	}

	sort.Sort(ByUploadTime(fMetaArray))
	return fMetaArray[0:count]
}

// 从mysql中获取批量的文件元信息列表
func GetLastFileMetaDB(limit int) ([]FileMeta, error) {
	tfiles, err := mydb.GetFileMetasDB(limit)
	if err != nil {
		return make([]FileMeta, 0), err
	}
	fmetas := make([]FileMeta, limit)
	for i := 0; i < len(tfiles); i++ {
		tfile := tfiles[i]
		fmeta := FileMeta{
			FileSha1: tfile.FileHash,
			FileName: tfile.FileName.String,
			FileSize: tfile.FileSize,
			Location: tfile.FileAddr.String,
		}
		fmetas = append(fmetas, fmeta)
	}
	return fmetas, nil
}

func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}
