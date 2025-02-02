package db

import (
	mydb "FILE-SERVER/db/mysql"
	"fmt"
	"time"
)

// 用户文件表结构体
type UserFile struct {
	UseNname    string
	FileHash    string
	FileName    string
	FileSize    int64
	UploadAt    string
	LastUpdated string
}

// 更新文件表
func OnUserFileUploadFinished(username, filehash, filename string,
	filesize int64) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_user_file(`user_name`, `file_sha1`, `file_name`, " +
			"`file_size`, `upload_at`) values(?, ?, ?, ?, ?)")

	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, filehash, filename, filesize, time.Now())
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

// 批量获取用户上传的文件信息列表
func QueryUserFileMetas(username string, limit int) ([]UserFile, error) {
	stmt, err := mydb.DBConn().Prepare(
		"select file_sha1, file_name, file_size, upload_at, last_update from " +
			"tbl_user_file where user_name = ? limit ?")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, limit)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	var userFiles []UserFile
	for rows.Next() {
		ufile := UserFile{}
		err = rows.Scan(&ufile.FileHash, &ufile.FileName, &ufile.FileSize, &ufile.UploadAt, &ufile.LastUpdated)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		userFiles = append(userFiles, ufile)
	}

	return userFiles, nil
}
