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
 * This sample demonstrates how to upload multiparts to oss
 * using the oss SDK for Go.
 */
package examples

import (
	"OSS"
	"fmt"
	"strings"
)

type SimpleMultipartUploadSample struct {
	bucketName string
	objectKey  string
	location   string
	OSSClient  *OSS.OSSClient
}

func newSimpleMultipartUploadSample(ak, sk, endpoint, bucketName, objectKey, location string) *SimpleMultipartUploadSample {
	OSSClient, err := OSS.New(ak, sk, endpoint)
	if err != nil {
		panic(err)
	}
	return &SimpleMultipartUploadSample{OSSClient: OSSClient, bucketName: bucketName, objectKey: objectKey, location: location}
}

func (sample SimpleMultipartUploadSample) CreateBucket() {
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

func (sample SimpleMultipartUploadSample) InitiateMultipartUpload() string {
	input := &OSS.InitiateMultipartUploadInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	output, err := sample.OSSClient.InitiateMultipartUpload(input)
	if err != nil {
		panic(err)
	}
	return output.UploadId
}

func (sample SimpleMultipartUploadSample) UploadPart(uploadId string) (string, int) {
	input := &OSS.UploadPartInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	input.UploadId = uploadId
	input.PartNumber = 1
	input.Body = strings.NewReader("Hello oss")
	output, err := sample.OSSClient.UploadPart(input)
	if err != nil {
		panic(err)
	}
	return output.ETag, output.PartNumber
}

func (sample SimpleMultipartUploadSample) CompleteMultipartUpload(uploadId, etag string, partNumber int) {
	input := &OSS.CompleteMultipartUploadInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	input.UploadId = uploadId
	input.Parts = []OSS.Part{
		OSS.Part{PartNumber: partNumber, ETag: etag},
	}
	_, err := sample.OSSClient.CompleteMultipartUpload(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Upload object %s successfully!\n", sample.objectKey)
}

func RunSimpleMultipartUploadSample() {
	const (
		endpoint   = "https://your-endpoint"
		ak         = "*** Provide your Access Key ***"
		sk         = "*** Provide your Secret Key ***"
		bucketName = "bucket-test"
		objectKey  = "object-test"
		location   = "yourbucketlocation"
	)
	sample := newSimpleMultipartUploadSample(ak, sk, endpoint, bucketName, objectKey, location)

	fmt.Println("Create a new bucket for demo")
	sample.CreateBucket()

	// Step 1: initiate multipart upload
	fmt.Println("Step 1: initiate multipart upload")
	uploadId := sample.InitiateMultipartUpload()

	// Step 2: upload a part
	fmt.Println("Step 2: upload a part")

	etag, partNumber := sample.UploadPart(uploadId)

	// Step 3: complete multipart upload
	fmt.Println("Step 3: complete multipart upload")
	sample.CompleteMultipartUpload(uploadId, etag, partNumber)

}
