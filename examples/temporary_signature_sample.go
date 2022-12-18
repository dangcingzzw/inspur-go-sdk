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
 * This sample demonstrates how to do common operations in temporary signature way
 * on OSS using the OSS SDK for Go.
 */
package examples

import (
	"OSS"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type TemporarySignatureSample struct {
	bucketName string
	objectKey  string
	location   string
	OSSClient  *OSS.OSSClient
}

func newTemporarySignatureSample(ak, sk, endpoint, bucketName, objectKey, location string) *TemporarySignatureSample {
	OSSClient, err := OSS.New(ak, sk, endpoint)
	if err != nil {
		panic(err)
	}
	return &TemporarySignatureSample{OSSClient: OSSClient, bucketName: bucketName, objectKey: objectKey, location: location}
}

func (sample TemporarySignatureSample) CreateBucket() {
	input := &OSS.CreateSignedUrlInput{}
	input.Bucket = sample.bucketName
	input.Method = OSS.HttpMethodPut
	input.Expires = 3600
	output, err := sample.OSSClient.CreateSignedUrl(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s using temporary signature url:\n", "CreateBucket")
	fmt.Println(output.SignedUrl)

	data := strings.NewReader(fmt.Sprintf("<CreateBucketConfiguration><LocationConstraint>%s</LocationConstraint></CreateBucketConfiguration>", sample.location))

	_, err = sample.OSSClient.CreateBucketWithSignedUrl(output.SignedUrl, output.ActualSignedRequestHeaders, data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Create bucket:%s successfully!\n", sample.bucketName)
	fmt.Println()
}

func (sample TemporarySignatureSample) ListBuckets() {
	input := &OSS.CreateSignedUrlInput{}
	input.Method = OSS.HttpMethodGet
	input.Expires = 3600
	output, err := sample.OSSClient.CreateSignedUrl(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s using temporary signature url:\n", "ListBuckets")
	fmt.Println(output.SignedUrl)

	listBucketsOutput, err := sample.OSSClient.ListBucketsWithSignedUrl(output.SignedUrl, output.ActualSignedRequestHeaders)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Owner.DisplayName:%s, Owner.ID:%s\n", listBucketsOutput.Owner.DisplayName, listBucketsOutput.Owner.ID)
	for index, val := range listBucketsOutput.Buckets {
		fmt.Printf("Bucket[%d]-Name:%s,CreationDate:%s\n", index, val.Name, val.CreationDate)
	}
	fmt.Println()
}

func (sample TemporarySignatureSample) DoBucketCors() {

	rawData := "<CORSConfiguration>" +
		"<CORSRule>" +
		"<AllowedOrigin>http://www.a.com</AllowedOrigin>" +
		"<AllowedMethod>PUT</AllowedMethod>" +
		"<AllowedMethod>POST</AllowedMethod>" +
		"<AllowedMethod>DELETE</AllowedMethod>" +
		"<AllowedHeader>*</AllowedHeader>" +
		"</CORSRule>" +
		"<CORSRule>" +
		"<AllowedOrigin>http://www.b.com</AllowedOrigin>" +
		"<AllowedMethod>GET</AllowedMethod>" +
		"</CORSRule>" +
		"</CORSConfiguration>"

	input := &OSS.CreateSignedUrlInput{}
	input.Method = OSS.HttpMethodPut
	input.Bucket = sample.bucketName
	input.SubResource = OSS.SubResourceCors
	input.Expires = 3600
	input.Headers = map[string]string{OSS.HEADER_MD5_CAMEL: OSS.Base64Md5([]byte(rawData))}
	output, err := sample.OSSClient.CreateSignedUrl(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s using temporary signature url:\n", "SetBucketCors")
	fmt.Println(output.SignedUrl)

	data := strings.NewReader(rawData)
	_, err = sample.OSSClient.SetBucketCorsWithSignedUrl(output.SignedUrl, output.ActualSignedRequestHeaders, data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Set bucket cors:%s successfully!\n", sample.bucketName)
	fmt.Println()

	input.Method = OSS.HttpMethodGet
	output, err = sample.OSSClient.CreateSignedUrl(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s using temporary signature url:\n", "GetBucketCors")
	fmt.Println(output.SignedUrl)

	getBucketCorsOutput, err := sample.OSSClient.GetBucketCorsWithSignedUrl(output.SignedUrl, output.ActualSignedRequestHeaders)
	if err != nil {
		panic(err)
	}
	for index, corsRule := range getBucketCorsOutput.CorsRules {
		fmt.Printf("CorsRule[%d]\n", index)
		fmt.Printf("ID:%s, AllowedOrigin:%s, AllowedMethod:%s, AllowedHeader:%s, MaxAgeSeconds:%d, ExposeHeader:%s\n",
			corsRule.ID, strings.Join(corsRule.AllowedOrigin, "|"), strings.Join(corsRule.AllowedMethod, "|"),
			strings.Join(corsRule.AllowedHeader, "|"), corsRule.MaxAgeSeconds, strings.Join(corsRule.ExposeHeader, "|"))
	}
	fmt.Println()

	input.Method = OSS.HttpMethodDelete
	output, err = sample.OSSClient.CreateSignedUrl(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s using temporary signature url:\n", "DeleteBucketCors")
	fmt.Println(output.SignedUrl)

	_, err = sample.OSSClient.DeleteBucketCorsWithSignedUrl(output.SignedUrl, output.ActualSignedRequestHeaders)
	if err != nil {
		panic(err)
	}
	fmt.Println("Delete bucket cors successfully!")
	fmt.Println()
}

func (sample TemporarySignatureSample) PutObject() {
	input := &OSS.CreateSignedUrlInput{}
	input.Method = OSS.HttpMethodPut
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	input.Expires = 3600
	output, err := sample.OSSClient.CreateSignedUrl(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s using temporary signature url:\n", "PutObject")
	fmt.Println(output.SignedUrl)

	data := strings.NewReader("Hello OSS")
	_, err = sample.OSSClient.PutObjectWithSignedUrl(output.SignedUrl, output.ActualSignedRequestHeaders, data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Put object:%s successfully!\n", sample.objectKey)
	fmt.Println()
}

func (TemporarySignatureSample) createSampleFile(sampleFilePath string) {
	if err := os.MkdirAll(filepath.Dir(sampleFilePath), os.ModePerm); err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(sampleFilePath, []byte("Hello OSS from file"), os.ModePerm); err != nil {
		panic(err)
	}
}

func (sample TemporarySignatureSample) PutFile(sampleFilePath string) {
	input := &OSS.CreateSignedUrlInput{}
	input.Method = OSS.HttpMethodPut
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	input.Expires = 3600
	output, err := sample.OSSClient.CreateSignedUrl(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s using temporary signature url:\n", "PutFile")
	fmt.Println(output.SignedUrl)

	_, err = sample.OSSClient.PutFileWithSignedUrl(output.SignedUrl, output.ActualSignedRequestHeaders, sampleFilePath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Put file:%s successfully!\n", sample.objectKey)
	fmt.Println()
}

func (sample TemporarySignatureSample) GetObject() {
	input := &OSS.CreateSignedUrlInput{}
	input.Method = OSS.HttpMethodGet
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	input.Expires = 3600
	output, err := sample.OSSClient.CreateSignedUrl(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s using temporary signature url:\n", "GetObject")
	fmt.Println(output.SignedUrl)

	getObjectOutput, err := sample.OSSClient.GetObjectWithSignedUrl(output.SignedUrl, output.ActualSignedRequestHeaders)
	if err != nil {
		panic(err)
	}
	defer func() {
		errMsg := getObjectOutput.Body.Close()
		if errMsg != nil {
			panic(errMsg)
		}
	}()
	fmt.Println("Object content:")
	body, err := ioutil.ReadAll(getObjectOutput.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
	fmt.Println()
}

func (sample TemporarySignatureSample) DoObjectAcl() {
	input := &OSS.CreateSignedUrlInput{}
	input.Method = OSS.HttpMethodPut
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	input.SubResource = OSS.SubResourceAcl
	input.Expires = 3600
	input.Headers = map[string]string{OSS.HEADER_ACL_AMZ: string(OSS.AclPublicRead)}
	output, err := sample.OSSClient.CreateSignedUrl(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s using temporary signature url:\n", "SetObjectAcl")
	fmt.Println(output.SignedUrl)

	_, err = sample.OSSClient.SetObjectAclWithSignedUrl(output.SignedUrl, output.ActualSignedRequestHeaders, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Set object acl:%s successfully!\n", sample.objectKey)
	fmt.Println()

	input.Method = OSS.HttpMethodGet
	output, err = sample.OSSClient.CreateSignedUrl(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s using temporary signature url:\n", "GetObjectAcl")
	fmt.Println(output.SignedUrl)

	getObjectAclOutput, err := sample.OSSClient.GetObjectAclWithSignedUrl(output.SignedUrl, output.ActualSignedRequestHeaders)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Object owner - ownerId:%s, ownerName:%s\n", getObjectAclOutput.Owner.ID, getObjectAclOutput.Owner.DisplayName)
	for index, grant := range getObjectAclOutput.Grants {
		fmt.Printf("Grant[%d]\n", index)
		fmt.Printf("GranteeUri:%s, GranteeId:%s, GranteeName:%s\n", grant.Grantee.URI, grant.Grantee.ID, grant.Grantee.DisplayName)
		fmt.Printf("Permission:%s\n", grant.Permission)
	}
	fmt.Println()
}

func (sample TemporarySignatureSample) DeleteObject() {
	input := &OSS.CreateSignedUrlInput{}
	input.Method = OSS.HttpMethodDelete
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	input.Expires = 3600
	output, err := sample.OSSClient.CreateSignedUrl(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s using temporary signature url:\n", "DeleteObject")
	fmt.Println(output.SignedUrl)

	_, err = sample.OSSClient.DeleteObjectWithSignedUrl(output.SignedUrl, output.ActualSignedRequestHeaders)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Delete object:%s successfully!\n", sample.objectKey)
	fmt.Println()
}

func (sample TemporarySignatureSample) DeleteBucket() {
	input := &OSS.CreateSignedUrlInput{}
	input.Method = OSS.HttpMethodDelete
	input.Bucket = sample.bucketName
	input.Expires = 3600
	output, err := sample.OSSClient.CreateSignedUrl(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s using temporary signature url:\n", "DeleteBucket")
	fmt.Println(output.SignedUrl)

	_, err = sample.OSSClient.DeleteBucketWithSignedUrl(output.SignedUrl, output.ActualSignedRequestHeaders)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Delete bucket:%s successfully!\n", sample.bucketName)
	fmt.Println()
}

func RunTemporarySignatureSample() {
	const (
		endpoint   = "https://your-endpoint"
		ak         = "*** Provide your Access Key ***"
		sk         = "*** Provide your Secret Key ***"
		bucketName = "bucket-test"
		objectKey  = "object-test"
		location   = "yourbucketlocation"
	)

	sample := newTemporarySignatureSample(ak, sk, endpoint, bucketName, objectKey, location)

	// Create bucket
	sample.CreateBucket()

	// List buckets
	sample.ListBuckets()

	// Set/Get/Delete bucket cors
	sample.DoBucketCors()

	// Put object
	sample.PutObject()

	// Get object
	sample.GetObject()

	// Put file
	sampleFilePath := "/temp/sampleText.txt"
	sample.createSampleFile(sampleFilePath)

	sample.PutFile(sampleFilePath)
	// Get object
	sample.GetObject()

	// Set/Get object acl
	sample.DoObjectAcl()

	// Delete object
	sample.DeleteObject()

	// Delete bucket
	sample.DeleteBucket()
}
