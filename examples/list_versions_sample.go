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
 * This sample demonstrates how to list versions under specified bucket
 * from oss using the oss SDK for Go.
 */
package examples

import (
	"OSS"
	"fmt"
	"strconv"
	"strings"
)

type ListVersionsSample struct {
	bucketName string
	location   string
	OSSClient  *OSS.OSSClient
}

func newListVersionsSample(ak, sk, endpoint, bucketName, location string) *ListVersionsSample {
	OSSClient, err := OSS.New(ak, sk, endpoint)
	if err != nil {
		panic(err)
	}
	return &ListVersionsSample{OSSClient: OSSClient, bucketName: bucketName, location: location}
}

func (sample ListVersionsSample) CreateBucket() {
	input := &OSS.CreateBucketInput{}
	input.Bucket = sample.bucketName
	input.Location = sample.location
	_, err := sample.OSSClient.CreateBucket(input)
	if err != nil {
		panic(err)
	}

	setBucketVersioningInput := &OSS.SetBucketVersioningInput{}
	setBucketVersioningInput.Bucket = sample.bucketName
	setBucketVersioningInput.Status = OSS.VersioningStatusEnabled
	_, err = sample.OSSClient.SetBucketVersioning(setBucketVersioningInput)

	if err != nil {
		panic(err)
	}
	fmt.Printf("Create bucket:%s successfully!\n", sample.bucketName)
	fmt.Println()
}

func (sample ListVersionsSample) preparePutObject(input *OSS.PutObjectInput) {
	_, err := sample.OSSClient.PutObject(input)
	if err != nil {
		panic(err)
	}
}

func (sample ListVersionsSample) PrepareFoldersAndObjects() {

	keyPrefix := "MyObjectKeyVersions"
	folderPrefix := "src"
	subFolderPrefix := "test"

	input := &OSS.PutObjectInput{}
	input.Bucket = sample.bucketName

	// First prepare folders and sub folders
	for i := 0; i < 5; i++ {
		key := folderPrefix + strconv.Itoa(i) + "/"
		input.Key = key
		sample.preparePutObject(input)
		for j := 0; j < 3; j++ {
			subKey := key + subFolderPrefix + strconv.Itoa(j) + "/"
			input.Key = subKey
			sample.preparePutObject(input)
		}
	}

	// Insert 2 objects in each folder
	input.Body = strings.NewReader("Hello oss")
	listVersionsInput := &OSS.ListVersionsInput{}
	listVersionsInput.Bucket = sample.bucketName
	output, err := sample.OSSClient.ListVersions(listVersionsInput)
	if err != nil {
		panic(err)
	}
	for _, version := range output.Versions {
		for i := 0; i < 2; i++ {
			objectKey := version.Key + keyPrefix + strconv.Itoa(i)
			input.Key = objectKey
			sample.preparePutObject(input)
		}
	}

	// Insert 2 objects in root path
	input.Key = keyPrefix + strconv.Itoa(0)
	sample.preparePutObject(input)
	input.Key = keyPrefix + strconv.Itoa(1)
	sample.preparePutObject(input)

	fmt.Println("Prepare folders and objects finished")
	fmt.Println()
}

func (sample ListVersionsSample) ListVersionsInFolders() {
	fmt.Println("List versions in folder src0/")
	input := &OSS.ListVersionsInput{}
	input.Bucket = sample.bucketName
	input.Prefix = "src0/"
	output, err := sample.OSSClient.ListVersions(input)
	if err != nil {
		panic(err)
	}
	for index, val := range output.Versions {
		fmt.Printf("Version[%d]-ETag:%s, Key:%s, Size:%d, VersionId:%s\n",
			index, val.ETag, val.Key, val.Size, val.VersionId)
	}

	fmt.Println()

	fmt.Println("List versions in sub folder src0/test0/")

	input.Prefix = "src0/test0/"
	output, err = sample.OSSClient.ListVersions(input)
	if err != nil {
		panic(err)
	}
	for index, val := range output.Versions {
		fmt.Printf("Version[%d]-ETag:%s, Key:%s, Size:%d, VersionId:%s\n",
			index, val.ETag, val.Key, val.Size, val.VersionId)
	}

	fmt.Println()
}

func (sample ListVersionsSample) ListVersionsByPage() {

	pageSize := 10
	pageNum := 1
	input := &OSS.ListVersionsInput{}
	input.Bucket = sample.bucketName
	input.MaxKeys = pageSize

	for {
		output, err := sample.OSSClient.ListVersions(input)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Page:%d\n", pageNum)
		for index, val := range output.Versions {
			fmt.Printf("Version[%d]-ETag:%s, Key:%s, Size:%d, VersionId:%s\n",
				index, val.ETag, val.Key, val.Size, val.VersionId)
		}
		if output.IsTruncated {
			input.KeyMarker = output.NextKeyMarker
			input.VersionIdMarker = output.NextVersionIdMarker
			pageNum++
		} else {
			break
		}
	}

	fmt.Println()
}

func (sample ListVersionsSample) listVersionsByPrefixes(commonPrefixes []string) {
	input := &OSS.ListVersionsInput{}
	input.Bucket = sample.bucketName
	input.Delimiter = "/"
	for _, prefix := range commonPrefixes {
		input.Prefix = prefix
		output, err := sample.OSSClient.ListVersions(input)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Folder %s:\n", prefix)
		for index, val := range output.Versions {
			fmt.Printf("Version[%d]-ETag:%s, Key:%s, Size:%d, VersionId:%s\n",
				index, val.ETag, val.Key, val.Size, val.VersionId)
		}
		fmt.Println()
		sample.listVersionsByPrefixes(output.CommonPrefixes)
	}
}

func (sample ListVersionsSample) ListVersionsGroupByFolder() {
	fmt.Println("List versions group by folder")
	input := &OSS.ListVersionsInput{}
	input.Bucket = sample.bucketName
	input.Delimiter = "/"
	output, err := sample.OSSClient.ListVersions(input)
	if err != nil {
		panic(err)
	}
	fmt.Println("Root path:")
	for index, val := range output.Versions {
		fmt.Printf("Version[%d]-ETag:%s, Key:%s, Size:%d, VersionId:%s\n",
			index, val.ETag, val.Key, val.Size, val.VersionId)
	}
	fmt.Println()
	sample.listVersionsByPrefixes(output.CommonPrefixes)
}

func (sample ListVersionsSample) BatchDeleteVersions() {
	input := &OSS.ListVersionsInput{}
	input.Bucket = sample.bucketName
	output, err := sample.OSSClient.ListVersions(input)
	if err != nil {
		panic(err)
	}
	objects := make([]OSS.ObjectToDelete, 0, len(output.Versions))
	for _, val := range output.Versions {
		objects = append(objects, OSS.ObjectToDelete{Key: val.Key, VersionId: val.VersionId})
	}
	deleteObjectsInput := &OSS.DeleteObjectsInput{}
	deleteObjectsInput.Bucket = sample.bucketName
	deleteObjectsInput.Objects = objects[:]
	_, err = sample.OSSClient.DeleteObjects(deleteObjectsInput)
	if err != nil {
		panic(err)
	}
	fmt.Println("Delete versions successfully!")
	fmt.Println()
}

func RunListVersionsSample() {

	const (
		endpoint   = "https://your-endpoint"
		ak         = "*** Provide your Access Key ***"
		sk         = "*** Provide your Secret Key ***"
		bucketName = "bucket-test"
		location   = "yourbucketlocation"
	)

	sample := newListVersionsSample(ak, sk, endpoint, bucketName, location)

	fmt.Println("Create a new bucket for demo")
	sample.CreateBucket()

	// First prepare folders and objects
	sample.PrepareFoldersAndObjects()

	// List versions in folders
	sample.ListVersionsInFolders()

	// List versions in way of pagination
	sample.ListVersionsByPage()

	// List versions group by folder
	sample.ListVersionsGroupByFolder()

	sample.BatchDeleteVersions()
}
