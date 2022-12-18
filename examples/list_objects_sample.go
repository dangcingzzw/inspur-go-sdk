// Copyright 2019 Inspur Technologies Co.,Ltd.
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use
// this file except in compliance with the License.  You may obtain a copy of the
// License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations under the License.

/**
 * This sample demonstrates how to list objects under specified bucket
 * from oss using the oss SDK for Go.
 */
package examples

import (
	"fmt"
	"oss"
	"strconv"
	"strings"
)

type ListObjectsSample struct {
	bucketName string
	location   string
	OSSClient  *oss.OSSClient
}

func newListObjectsSample(ak, sk, endpoint, bucketName, location string) *ListObjectsSample {
	OSSClient, err := oss.New(ak, sk, endpoint)
	if err != nil {
		panic(err)
	}
	return &ListObjectsSample{OSSClient: OSSClient, bucketName: bucketName, location: location}
}

func (sample ListObjectsSample) CreateBucket() {
	input := &oss.CreateBucketInput{}
	input.Bucket = sample.bucketName
	input.Location = sample.location
	_, err := sample.OSSClient.CreateBucket(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Create bucket:%s successfully!\n", sample.bucketName)
	fmt.Println()
}

func (sample ListObjectsSample) DoInsertObjects() []string {

	keyPrefix := "MyObjectKey"

	input := &oss.PutObjectInput{}
	input.Bucket = sample.bucketName
	input.Body = strings.NewReader("Hello oss")
	keys := make([]string, 0, 100)
	for i := 0; i < 100; i++ {
		input.Key = keyPrefix + strconv.Itoa(i)
		_, err := sample.OSSClient.PutObject(input)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Succeed to put object %s\n", input.Key)
		keys = append(keys, input.Key)
	}
	fmt.Println()
	return keys
}

func (sample ListObjectsSample) ListObjects() {
	input := &oss.ListObjectsInput{}
	input.Bucket = sample.bucketName
	output, err := sample.OSSClient.ListObjects(input)
	if err != nil {
		panic(err)
	}
	for index, val := range output.Contents {
		fmt.Printf("Content[%d]-ETag:%s, Key:%s, Size:%d\n",
			index, val.ETag, val.Key, val.Size)
	}
	fmt.Println()
}

func (sample ListObjectsSample) ListObjectsByMarker() {
	input := &oss.ListObjectsInput{}
	input.Bucket = sample.bucketName
	input.MaxKeys = 10
	output, err := sample.OSSClient.ListObjects(input)
	if err != nil {
		panic(err)
	}
	fmt.Println("List the first 10 objects :")
	for index, val := range output.Contents {
		fmt.Printf("Content[%d]-ETag:%s, Key:%s, Size:%d\n",
			index, val.ETag, val.Key, val.Size)
	}
	fmt.Println()

	input.Marker = output.NextMarker
	output, err = sample.OSSClient.ListObjects(input)
	if err != nil {
		panic(err)
	}
	fmt.Println("List the second 10 objects using marker:")
	for index, val := range output.Contents {
		fmt.Printf("Content[%d]-ETag:%s, Key:%s, Size:%d\n",
			index, val.ETag, val.Key, val.Size)
	}
	fmt.Println()
}

func (sample ListObjectsSample) ListObjectsByPage() {

	pageSize := 10
	pageNum := 1
	input := &oss.ListObjectsInput{}
	input.Bucket = sample.bucketName
	input.MaxKeys = pageSize

	for {
		output, err := sample.OSSClient.ListObjects(input)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Page:%d\n", pageNum)
		for index, val := range output.Contents {
			fmt.Printf("Content[%d]-ETag:%s, Key:%s, Size:%d\n",
				index, val.ETag, val.Key, val.Size)
		}
		if output.IsTruncated {
			input.Marker = output.NextMarker
			pageNum++
		} else {
			break
		}
	}

	fmt.Println()
}

func (sample ListObjectsSample) DeleteObjects(keys []string) {
	input := &oss.DeleteObjectsInput{}
	input.Bucket = sample.bucketName

	objects := make([]oss.ObjectToDelete, 0, len(keys))
	for _, key := range keys {
		objects = append(objects, oss.ObjectToDelete{Key: key})
	}
	input.Objects = objects
	_, err := sample.OSSClient.DeleteObjects(input)
	if err != nil {
		panic(err)
	}
	fmt.Println("Delete objects successfully!")
}

func RunListObjectsSample() {

	const (
		endpoint   = "https://your-endpoint"
		ak         = "*** Provide your Access Key ***"
		sk         = "*** Provide your Secret Key ***"
		bucketName = "bucket-test"
		location   = "yourbucketlocation"
	)

	sample := newListObjectsSample(ak, sk, endpoint, bucketName, location)

	fmt.Println("Create a new bucket for demo")
	sample.CreateBucket()

	// First insert 100 objects for demo
	keys := sample.DoInsertObjects()

	// List objects using default parameters, will return up to 1000 objects
	sample.ListObjects()

	// List the first 10 and second 10 objects
	sample.ListObjectsByMarker()

	// List objects in way of pagination
	sample.ListObjectsByPage()

	// Delete all the objects created
	sample.DeleteObjects(keys)
}
