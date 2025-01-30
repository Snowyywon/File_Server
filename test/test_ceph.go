package main

import (
	"FILE-SERVER/store/ceph"
	"fmt"
	"os"
)

func main() {
	bucket := ceph.GetCephBucket("testbucket1")
	d, _ := bucket.Get("/ceph/1b186aa919832cb3477e8a0d66acb699fc84215a")
	// fmt.Println(d)
	tmpFile, _ := os.Create("/tmp/test_file")
	tmpFile.Write(d)

	res, _ := bucket.List("", "", "", 100)
	fmt.Printf("object keys : %+v \n", res)
	return

	// 列出桶中的对象
	// objects, err := bucket.List("", "", "", 100)
	// if err != nil {
	// 	fmt.Printf("Error listing objects in bucket: %v\n", err)
	// 	return
	// }

	// // 打印所有对象的名称
	// fmt.Println("Objects in bucket 'userfile':")
	// for _, object := range objects {
	// 	fmt.Println(object) // 这里假设 List 返回的是对象名，具体视实现而定
	// }

	// return

	// 创建一个新的bucket
	// err := bucket.PutBucket(s3.PublicRead)
	// fmt.Println("create bucket err : ", err)

	// 查询这个bucket下面指定条件的object keys

	// res, _ := bucket.List("", "", "", 100)
	// fmt.Printf("object keys : %+v \n", res)

	//新上传一个对象
	// err = bucket.Put("/test/upload/a.txt", []byte("just for test"), "octet-stream",
	// 	s3.PublicRead)
	// fmt.Printf("upload err : %+v\n", err)

	// 查询这个bucket下面指定条件的object keys
	// res, err = bucket.List("", "", "", 100)
	// fmt.Printf("object keys : %+v \n", res)

}
