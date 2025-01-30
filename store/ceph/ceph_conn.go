package ceph

import (
	"gopkg.in/amz.v1/aws"

	"gopkg.in/amz.v1/s3"
)

var cephCoon *s3.S3

// 获取s3连接
func GetCephCoonetction() *s3.S3 {
	// 1. 初始化ceph的一些信息
	if cephCoon != nil {
		return cephCoon
	}
	auth := aws.Auth{
		AccessKey: "JVCVH5B21XW34WHMCUOV",
		SecretKey: "NeolBWjIOqKlZ3DeKS4ZHdr7mTZ1UBY9CU93JVOZ",
	}

	curRegion := aws.Region{
		Name:                 "default",
		EC2Endpoint:          "http://127.0.0.1:7480",
		S3Endpoint:           "http://127.0.0.1:7480",
		S3BucketEndpoint:     "",
		S3LocationConstraint: false,
		S3LowercaseBucket:    false,
		Sign:                 aws.SignV2,
	}
	// 2. 创建S3类型的连接
	return s3.New(auth, curRegion)
}

// 获取s3的bucket
func GetCephBucket(bucket string) *s3.Bucket {
	conn := GetCephCoonetction()
	return conn.Bucket(bucket)
}

// PutObject : 上传文件到ceph集群
func PutObject(bucket string, path string, data []byte) error {
	return GetCephBucket(bucket).Put(path, data, "octet-stream", s3.PublicRead)
}
