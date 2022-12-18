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
 * This sample demonstrates how to delete objects under specified bucket
 * from oss using the oss SDK for Go.
 */
package examples

import (
	"OSS"
	"fmt"
	"strconv"
	"strings"
)

const (
	MyObjectKey string = "MyObjectKey"
)

type DeleteObjectsSample struct {
	bucketName string
	location   string
	OSSClient  *OSS.OSSClient
}

func newDeleteObjectsSample(ak, sk, endpoint, bucketName, location string) *DeleteObjectsSample {
	OSSClient, err := OSS.New(ak, sk, endpoint)
	if err != nil {
		panic(err)
	}
	return &DeleteObjectsSample{OSSClient: OSSClient, bucketName: bucketName, location: location}
}

func (sample DeleteObjectsSample) CreateBucket() {
	input := &OSS.CreateBucketInput{}
	input.Bucket = sample.bucketName
	input.Location = sample.location
	_, err := sample.OSSClient.CreateBucket(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Create bucket:%s successfully!\n", sample.bucketName)
	fmt.Println()
}

func (sample DeleteObjectsSample) BatchPutObjects() {
	content := "Thank you for using Object Storage Service"
	keyPrefix := MyObjectKey

	input := &OSS.PutObjectInput{}
	input.Bucket = sample.bucketName
	input.Body = strings.NewReader(content)
	for i := 0; i < 100; i++ {
		input.Key = keyPrefix + strconv.Itoa(i)
		_, err := sample.OSSClient.PutObject(input)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Succeed to put object %s\n", input.Key)
	}
}

func (sample DeleteObjectsSample) BatchDeleteObjects() {
	input := &OSS.ListObjectsInput{}
	input.Bucket = sample.bucketName
	output, err := sample.OSSClient.ListObjects(input)
	if err != nil {
		panic(err)
	}
	objects := make([]OSS.ObjectToDelete, 0, len(output.Contents))
	for _, content := range output.Contents {
		objects = append(objects, OSS.ObjectToDelete{Key: content.Key})
	}
	deleteObjectsInput := &OSS.DeleteObjectsInput{}
	deleteObjectsInput.Bucket = sample.bucketName
	deleteObjectsInput.Objects = objects[:]
	deleteObjectsOutput, err := sample.OSSClient.DeleteObjects(deleteObjectsInput)
	if err != nil {
		panic(err)
	}
	for _, deleted := range deleteObjectsOutput.Deleteds {
		fmt.Printf("Delete %s successfully\n", deleted.Key)
	}
	fmt.Println()
	for _, deleteError := range deleteObjectsOutput.Errors {
		fmt.Printf("Delete %s failed, code:%s, message:%s\n", deleteError.Key, deleteError.Code, deleteError.Message)
	}
	fmt.Println()
}

func RunDeleteObjectsSample() {
	const (
		endpoint   = "https://your-endpoint"
		ak         = "*** Provide your Access Key ***"
		sk         = "*** Provide your Secret Key ***"
		bucketName = "bucket-test"
		location   = "yourbucketlocation"
	)
	sample := newDeleteObjectsSample(ak, sk, endpoint, bucketName, location)

	fmt.Println("Create a new bucket for demo")
	sample.CreateBucket()

	// Batch put objects into the bucket
	sample.BatchPutObjects()

	// Delete all objects uploaded recently under the bucket
	sample.BatchDeleteObjects()
}
