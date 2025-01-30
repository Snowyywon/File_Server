package oss

import (
	cfg "FILE-SERVER/config"
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var ossCli *oss.Client

// 返回oss.client对象
func Client() *oss.Client {
	if ossCli != nil {
		return ossCli
	}
	ossCli, err := oss.New(cfg.OSSEndpoint, cfg.OSSAccesskeyID, cfg.OSSAccessKeySecret)

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return ossCli
}

// 用于获取bucket存储空间
func Bucket() *oss.Bucket {
	cli := Client()
	if cli != nil {
		bucket, err := cli.Bucket(cfg.OSSBucket)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		return bucket
	}
	return nil
}

// 临时授权下载
func DownloadURL(objName string) string {
	signedURL, err := Bucket().SignURL(objName, oss.HTTPGet, 3600)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	return signedURL
}

func BuildLifecycleRule(bucketName string) {
	ruleTest1 := oss.BuildLifecycleRuleByDays("rule1", "test/", true, 30)
	rules := []oss.LifecycleRule{ruleTest1}

	Client().SetBucketLifecycle(bucketName, rules)
}
