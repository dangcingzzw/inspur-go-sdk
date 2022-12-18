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
 * This sample demonstrates how to download an cold object
 * from OSS using the OSS SDK for Go.
 */
package examples

import (
	"OSS"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

type RestoreObjectSample struct {
	bucketName string
	objectKey  string
	location   string
	OSSClient  *OSS.OSSClient
}

func newRestoreObjectSample(ak, sk, endpoint, bucketName, objectKey, location string) *RestoreObjectSample {
	OSSClient, err := OSS.New(ak, sk, endpoint)
	if err != nil {
		panic(err)
	}
	return &RestoreObjectSample{OSSClient: OSSClient, bucketName: bucketName, objectKey: objectKey, location: location}
}

func (sample RestoreObjectSample) CreateColdBucket() {
	input := &OSS.CreateBucketInput{}
	input.Bucket = sample.bucketName
	input.Location = sample.location
	input.StorageClass = OSS.StorageClassCold
	_, err := sample.OSSClient.CreateBucket(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Create cold bucket:%s successfully!\n", sample.bucketName)
	fmt.Println()
}

func (sample RestoreObjectSample) CreateObject() {
	input := &OSS.PutObjectInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	input.Body = strings.NewReader("Hello OSS")

	_, err := sample.OSSClient.PutObject(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Create object:%s successfully!\n", sample.objectKey)
	fmt.Println()
}

func (sample RestoreObjectSample) RestoreObject() {
	input := &OSS.RestoreObjectInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	input.Days = 1
	input.Tier = OSS.RestoreTierExpedited

	_, err := sample.OSSClient.RestoreObject(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Create object:%s successfully!\n", sample.objectKey)
	fmt.Println()
}

func (sample RestoreObjectSample) GetObject() {
	input := &OSS.GetObjectInput{}
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

func (sample RestoreObjectSample) DeleteObject() {
	input := &OSS.DeleteObjectInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	_, err := sample.OSSClient.DeleteObject(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Delete object:%s successfully!\n", input.Key)
	fmt.Println()
}

func RunRestoreObjectSample() {
	const (
		endpoint   = "https://your-endpoint"
		ak         = "*** Provide your Access Key ***"
		sk         = "*** Provide your Secret Key ***"
		bucketName = "bucket-test-cold"
		objectKey  = "object-test"
		location   = "yourbucketlocation"
	)

	sample := newRestoreObjectSample(ak, sk, endpoint, bucketName, objectKey, location)

	fmt.Println("Create a new cold bucket for demo")
	sample.CreateColdBucket()

	sample.CreateObject()

	sample.RestoreObject()

	// Wait 6 minutes to get the object
	time.Sleep(time.Duration(6*60) * time.Second)

	sample.GetObject()

	sample.DeleteObject()
}
