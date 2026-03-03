package main

import (
	"cloud/store/ceph"
	"fmt"

	"gopkg.in/amz.v1/s3"
)

func main() {
	bucket := ceph.GetCephBucket("testbucket1")

	// 创建一个新的bucket
	err := bucket.PutBucket(s3.PublicRead)
	fmt.Printf("create bucket %s", bucket)

	// 查询这个bucket下面指定的object
	res, err := bucket.List("", "", "", 100)
	fmt.Printf("list bucket %s", res)

	// 新上传一个对象
	err = bucket.Put("/testupload/a.txt", []byte("just for test"), "octet-stream", s3.PublicRead)
	fmt.Printf("upload err: %s", err.Error())

	// 查询这个bucket下面指定条件的object keys
	res, err = bucket.List("", "", "", 100)
	fmt.Printf("list bucket %s", res)
}
