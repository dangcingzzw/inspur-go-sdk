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
 * This sample demonstrates how to download an object
 * from oss in different ways using the oss SDK for Go.
 */
package examples

import (
	"fmt"
	"io/ioutil"
	"os"
	"oss"
	"path/filepath"
	"strings"
)

type DownloadSample struct {
	bucketName string
	objectKey  string
	location   string
	OSSClient  *oss.OSSClient
}

func newDownloadSample(ak, sk, endpoint, bucketName, objectKey, location string) *DownloadSample {
	OSSClient, err := oss.New(ak, sk, endpoint)
	if err != nil {
		panic(err)
	}
	return &DownloadSample{OSSClient: OSSClient, bucketName: bucketName, objectKey: objectKey, location: location}
}

func (sample DownloadSample) CreateBucket() {
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

func (sample DownloadSample) PutObject() {
	input := &oss.PutObjectInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	input.Body = strings.NewReader("Hello oss")

	_, err := sample.OSSClient.PutObject(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Put object:%s successfully!\n", sample.objectKey)
	fmt.Println()
}

func (sample DownloadSample) GetObject() {
	input := &oss.GetObjectInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey

	output, err := sample.OSSClient.GetObject(input)
	if err != nil {
		panic(err)
	}
	defer func() {
		errMsg := output.Body.Close()
		if errMsg != nil {
			panic(errMsg)
		}
	}()
	fmt.Println("Object content:")
	body, err := ioutil.ReadAll(output.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
	fmt.Println()
}

func (sample DownloadSample) PutFile(sampleFilePath string) {
	input := &oss.PutFileInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	input.SourceFile = sampleFilePath

	_, err := sample.OSSClient.PutFile(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Put object:%s with file:%s successfully!\n", sample.objectKey, sampleFilePath)
	fmt.Println()
}

func (sample DownloadSample) DeleteObject() {
	input := &oss.DeleteObjectInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey

	_, err := sample.OSSClient.DeleteObject(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Delete object:%s successfully!\n", sample.objectKey)
	fmt.Println()
}

func (DownloadSample) createSampleFile(sampleFilePath string) {
	if err := os.MkdirAll(filepath.Dir(sampleFilePath), os.ModePerm); err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(sampleFilePath, []byte("Hello oss from file"), os.ModePerm); err != nil {
		panic(err)
	}
}

func RunDownloadSample() {
	const (
		endpoint   = "https://your-endpoint"
		ak         = "*** Provide your Access Key ***"
		sk         = "*** Provide your Secret Key ***"
		bucketName = "bucket-test"
		objectKey  = "object-test"
		location   = "yourbucketlocation"
	)
	sample := newDownloadSample(ak, sk, endpoint, bucketName, objectKey, location)

	fmt.Println("Create a new bucket for demo")
	sample.CreateBucket()

	fmt.Println("Uploading a new object to oss from string")
	sample.PutObject()

	fmt.Println("Download object to string")
	sample.GetObject()

	fmt.Println("Uploading a new object to oss from file")
	sampleFilePath := "/temp/text.txt"
	sample.createSampleFile(sampleFilePath)
	defer func() {
		errMsg := os.Remove(sampleFilePath)
		if errMsg != nil {
			panic(errMsg)
		}
	}()
	sample.PutFile(sampleFilePath)

	fmt.Println("Download file to string")
	sample.GetObject()

	sample.DeleteObject()
}
