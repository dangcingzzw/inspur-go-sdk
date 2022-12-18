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
 * This sample demonstrates how to do object-related operations
 * (such as create/delete/get/copy object, do object ACL)
 * on oss using the oss SDK for Go.
 */
package examples

import (
	"fmt"
	"io/ioutil"
	"oss"
	"strings"
)

type ObjectOperationsSample struct {
	bucketName string
	objectKey  string
	location   string
	OSSClient  *oss.OSSClient
}

func newObjectOperationsSample(ak, sk, endpoint, bucketName, objectKey, location string) *ObjectOperationsSample {
	OSSClient, err := oss.New(ak, sk, endpoint)
	if err != nil {
		panic(err)
	}
	return &ObjectOperationsSample{OSSClient: OSSClient, bucketName: bucketName, objectKey: objectKey, location: location}
}

func (sample ObjectOperationsSample) CreateBucket() {
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

func (sample ObjectOperationsSample) GetObjectMeta() {
	input := &oss.GetObjectMetadataInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	output, err := sample.OSSClient.GetObjectMetadata(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Object content-type:%s\n", output.ContentType)
	fmt.Printf("Object content-length:%d\n", output.ContentLength)
	fmt.Println()
}

func (sample ObjectOperationsSample) CreateObject() {
	input := &oss.PutObjectInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	input.Body = strings.NewReader("Hello oss")

	_, err := sample.OSSClient.PutObject(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Create object:%s successfully!\n", sample.objectKey)
	fmt.Println()
}

func (sample ObjectOperationsSample) GetObject() {
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

func (sample ObjectOperationsSample) CopyObject() {
	input := &oss.CopyObjectInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey + "-back"
	input.CopySourceBucket = sample.bucketName
	input.CopySourceKey = sample.objectKey

	_, err := sample.OSSClient.CopyObject(input)
	if err != nil {
		panic(err)
	}
	fmt.Println("Copy object successfully!")
	fmt.Println()
}

func (sample ObjectOperationsSample) DoObjectAcl() {
	input := &oss.SetObjectAclInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	input.ACL = oss.AclPublicRead

	_, err := sample.OSSClient.SetObjectAcl(input)
	if err != nil {
		panic(err)
	}
	fmt.Println("Set object acl successfully!")
	fmt.Println()

	output, err := sample.OSSClient.GetObjectAcl(&oss.GetObjectAclInput{Bucket: sample.bucketName, Key: sample.objectKey})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Object owner - ownerId:%s, ownerName:%s\n", output.Owner.ID, output.Owner.DisplayName)
	for index, grant := range output.Grants {
		fmt.Printf("Grant[%d]\n", index)
		fmt.Printf("GranteeUri:%s, GranteeId:%s, GranteeName:%s\n", grant.Grantee.URI, grant.Grantee.ID, grant.Grantee.DisplayName)
		fmt.Printf("Permission:%s\n", grant.Permission)
	}
}

func (sample ObjectOperationsSample) DeleteObject() {
	input := &oss.DeleteObjectInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey

	_, err := sample.OSSClient.DeleteObject(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Delete object:%s successfully!\n", input.Key)
	fmt.Println()

	input.Key = sample.objectKey + "-back"

	_, err = sample.OSSClient.DeleteObject(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Delete object:%s successfully!\n", input.Key)
	fmt.Println()
}

func RunObjectOperationsSample() {
	const (
		endpoint   = "https://your-endpoint"
		ak         = "*** Provide your Access Key ***"
		sk         = "*** Provide your Secret Key ***"
		bucketName = "bucket-test"
		objectKey  = "object-test"
		location   = "yourbucketlocation"
	)

	sample := newObjectOperationsSample(ak, sk, endpoint, bucketName, objectKey, location)

	fmt.Println("Create a new bucket for demo")
	sample.CreateBucket()

	sample.CreateObject()

	sample.GetObjectMeta()

	sample.GetObject()

	sample.CopyObject()

	sample.DoObjectAcl()

	sample.DeleteObject()
}
