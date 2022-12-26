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

package OSS

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// ListObjects lists objects in a bucket.
//
// You can use this API to list objects in a bucket. By default, a maximum of 1000 objects are listed.
func (OSSClient OSSClient) ListObjects(input *ListObjectsInput, extensions ...extensionOptions) (output *ListObjectsOutput, err error) {
	if input == nil {
		return nil, errors.New("ListObjectsInput is nil")
	}
	output = &ListObjectsOutput{}
	err = OSSClient.doActionWithBucket("ListObjects", HTTP_GET, input.Bucket, input, output, extensions)
	if err != nil {
		output = nil
	} else {
		if location, ok := output.ResponseHeaders[HEADER_BUCKET_REGION]; ok {
			output.Location = location[0]
		}
		if output.EncodingType == "url" {
			err = decodeListObjectsOutput(output)
			if err != nil {
				doLog(LEVEL_ERROR, "Failed to get ListObjectsOutput with error: %v.", err)
				output = nil
			}
		}
	}
	return
}

// ListVersions lists versioning objects in a bucket.
//
// You can use this API to list versioning objects in a bucket. By default, a maximum of 1000 versioning objects are listed.
func (OSSClient OSSClient) ListVersions(input *ListVersionsInput, extensions ...extensionOptions) (output *ListVersionsOutput, err error) {
	if input == nil {
		return nil, errors.New("ListVersionsInput is nil")
	}
	output = &ListVersionsOutput{}
	err = OSSClient.doActionWithBucket("ListVersions", HTTP_GET, input.Bucket, input, output, extensions)
	if err != nil {
		output = nil
	} else {
		if location, ok := output.ResponseHeaders[HEADER_BUCKET_REGION]; ok {
			output.Location = location[0]
		}
		if output.EncodingType == "url" {
			err = decodeListVersionsOutput(output)
			if err != nil {
				doLog(LEVEL_ERROR, "Failed to get ListVersionsOutput with error: %v.", err)
				output = nil
			}
		}
	}
	return
}

// HeadObject checks whether an object exists.
//
// You can use this API to check whether an object exists.
func (OSSClient OSSClient) HeadObject(input *HeadObjectInput, extensions ...extensionOptions) (output *BaseModel, err error) {
	if input == nil {
		return nil, errors.New("HeadObjectInput is nil")
	}
	output = &BaseModel{}
	err = OSSClient.doActionWithBucketAndKey("HeadObject", HTTP_HEAD, input.Bucket, input.Key, input, output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// SetObjectMetadata sets object metadata.
func (OSSClient OSSClient) SetObjectMetadata(input *SetObjectMetadataInput, extensions ...extensionOptions) (output *SetObjectMetadataOutput, err error) {
	output = &SetObjectMetadataOutput{}
	err = OSSClient.doActionWithBucketAndKey("SetObjectMetadata", HTTP_PUT, input.Bucket, input.Key, input, output, extensions)
	if err != nil {
		output = nil
	} else {
		ParseSetObjectMetadataOutput(output)
	}
	return
}

// DeleteObject deletes an object.
//
// You can use this API to delete an object from a specified bucket.
func (OSSClient OSSClient) DeleteObject(input *DeleteObjectInput, extensions ...extensionOptions) (output *DeleteObjectOutput, err error) {
	if input == nil {
		return nil, errors.New("DeleteObjectInput is nil")
	}
	output = &DeleteObjectOutput{}
	err = OSSClient.doActionWithBucketAndKey("DeleteObject", HTTP_DELETE, input.Bucket, input.Key, input, output, extensions)
	if err != nil {
		output = nil
	} else {
		ParseDeleteObjectOutput(output)
	}
	return
}

// DeleteObjects deletes objects in a batch.
//
// You can use this API to batch delete objects from a specified bucket.
func (OSSClient OSSClient) DeleteObjects(input *DeleteObjectsInput, extensions ...extensionOptions) (output *DeleteObjectsOutput, err error) {
	if input == nil {
		return nil, errors.New("DeleteObjectsInput is nil")
	}
	output = &DeleteObjectsOutput{}
	err = OSSClient.doActionWithBucket("DeleteObjects", HTTP_POST, input.Bucket, input, output, extensions)
	if err != nil {
		output = nil
	} else if output.EncodingType == "url" {
		err = decodeDeleteObjectsOutput(output)
		if err != nil {
			doLog(LEVEL_ERROR, "Failed to get DeleteObjectsOutput with error: %v.", err)
			output = nil
		}
	}
	return
}

// SetObjectAcl sets ACL for an object.
//
// You can use this API to set the ACL for an object in a specified bucket.
func (OSSClient OSSClient) SetObjectAcl(input *SetObjectAclInput, extensions ...extensionOptions) (output *BaseModel, err error) {
	if input == nil {
		return nil, errors.New("SetObjectAclInput is nil")
	}
	output = &BaseModel{}
	err = OSSClient.doActionWithBucketAndKey("SetObjectAcl", HTTP_PUT, input.Bucket, input.Key, input, output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// GetObjectAcl gets the ACL of an object.
//
// You can use this API to obtain the ACL of an object in a specified bucket.
func (OSSClient OSSClient) GetObjectAcl(input *GetObjectAclInput, extensions ...extensionOptions) (output *GetObjectAclOutput, err error) {
	if input == nil {
		return nil, errors.New("GetObjectAclInput is nil")
	}
	output = &GetObjectAclOutput{}
	err = OSSClient.doActionWithBucketAndKey("GetObjectAcl", HTTP_GET, input.Bucket, input.Key, input, output, extensions)
	if err != nil {
		output = nil
	} else {
		if versionID, ok := output.ResponseHeaders[HEADER_VERSION_ID]; ok {
			output.VersionId = versionID[0]
		}
	}
	return
}

// RestoreObject restores an object.
func (OSSClient OSSClient) RestoreObject(input *RestoreObjectInput, extensions ...extensionOptions) (output *BaseModel, err error) {
	if input == nil {
		return nil, errors.New("RestoreObjectInput is nil")
	}
	output = &BaseModel{}
	err = OSSClient.doActionWithBucketAndKey("RestoreObject", HTTP_POST, input.Bucket, input.Key, input, output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// GetObjectMetadata gets object metadata.
//
// You can use this API to send a HEAD request to the object of a specified bucket to obtain its metadata.
func (OSSClient OSSClient) GetObjectMetadata(input *GetObjectMetadataInput, extensions ...extensionOptions) (output *GetObjectMetadataOutput, err error) {
	if input == nil {
		return nil, errors.New("GetObjectMetadataInput is nil")
	}
	output = &GetObjectMetadataOutput{}
	err = OSSClient.doActionWithBucketAndKey("GetObjectMetadata", HTTP_HEAD, input.Bucket, input.Key, input, output, extensions)
	if err != nil {
		output = nil
	} else {
		ParseGetObjectMetadataOutput(output)
	}
	return
}

func (OSSClient OSSClient) GetAttribute(input *GetAttributeInput, extensions ...extensionOptions) (output *GetAttributeOutput, err error) {
	if input == nil {
		return nil, errors.New("GetAttributeInput is nil")
	}
	output = &GetAttributeOutput{}
	err = OSSClient.doActionWithBucketAndKey("GetAttribute", HTTP_HEAD, input.Bucket, input.Key, input, output, extensions)
	if err != nil {
		output = nil
	} else {
		ParseGetAttributeOutput(output)
	}
	return
}

// GetObject downloads object.
//
// You can use this API to download an object in a specified bucket.
func (OSSClient OSSClient) GetObject(input *GetObjectInput, extensions ...extensionOptions) (output *GetObjectOutput, err error) {
	if input == nil {
		return nil, errors.New("GetObjectInput is nil")
	}
	output = &GetObjectOutput{}
	err = OSSClient.doActionWithBucketAndKey("GetObject", HTTP_GET, input.Bucket, input.Key, input, output, extensions)
	if err != nil {
		output = nil
	} else {
		ParseGetObjectOutput(output)
	}
	return
}
func (OSSClient OSSClient) DoesObjectExist(input *DoesObjectExistInput, extensions ...extensionOptions) (output *DoesObjectExistOutput, err error) {
	if input == nil {
		return nil, errors.New("DoesObjectExistInput is nil")
	}
	output = &DoesObjectExistOutput{}
	err = OSSClient.doActionWithBucketAndKey("DoesObjectExist", HTTP_GET, input.Bucket, input.Key, input, output, extensions)
	if err != nil {
		output = nil
	} else {
		ParseDoesObjectExistOutput(output)
	}
	return
}

// PutObject uploads an object to the specified bucket.
func (OSSClient OSSClient) PutObject(input *PutObjectInput, extensions ...extensionOptions) (output *PutObjectOutput, err error) {
	if input == nil {
		return nil, errors.New("PutObjectInput is nil")
	}

	if input.ContentType == "" && input.Key != "" {
		if contentType, ok := mimeTypes[strings.ToLower(input.Key[strings.LastIndex(input.Key, ".")+1:])]; ok {
			input.ContentType = contentType
		}
	}
	output = &PutObjectOutput{}
	var repeatable bool
	if input.Body != nil {
		if _, ok := input.Body.(*strings.Reader); !ok {
			repeatable = false
		}
		if input.ContentLength > 0 {
			input.Body = &readerWrapper{reader: input.Body, totalCount: input.ContentLength}
		}
	}
	if repeatable {
		err = OSSClient.doActionWithBucketAndKey("PutObject", HTTP_PUT, input.Bucket, input.Key, input, output, extensions)
	} else {
		err = OSSClient.doActionWithBucketAndKeyUnRepeatable("PutObject", HTTP_PUT, input.Bucket, input.Key, input, output, extensions)
	}
	if err != nil {
		output = nil
	} else {
		ParsePutObjectOutput(output)
	}
	return
}

func (OSSClient OSSClient) getContentType(input *PutObjectInput, sourceFile string) (contentType string) {
	if contentType, ok := mimeTypes[strings.ToLower(input.Key[strings.LastIndex(input.Key, ".")+1:])]; ok {
		return contentType
	}
	if contentType, ok := mimeTypes[strings.ToLower(sourceFile[strings.LastIndex(sourceFile, ".")+1:])]; ok {
		return contentType
	}
	return
}

func (OSSClient OSSClient) isGetContentType(input *PutObjectInput) bool {
	if input.ContentType == "" && input.Key != "" {
		return true
	}
	return false
}

func (OSSClient OSSClient) NewFolder(input *NewFolderInput, extensions ...extensionOptions) (output *NewFolderOutput, err error) {
	if input == nil {
		return nil, errors.New("NewFolderInput is nil")
	}

	if !strings.HasSuffix(input.Key, "/") {
		input.Key += "/"
	}

	output = &NewFolderOutput{}
	err = OSSClient.doActionWithBucketAndKey("NewFolder", HTTP_PUT, input.Bucket, input.Key, input, output, extensions)
	if err != nil {
		output = nil
	} else {
		ParseNewFolderOutput(output)
		output.ObjectUrl = fmt.Sprintf("%s/%s/%s", OSSClient.conf.endpoint, input.Bucket, input.Key)
	}
	return
}

// PutFile uploads a file to the specified bucket.
func (OSSClient OSSClient) PutFile(input *PutFileInput, extensions ...extensionOptions) (output *PutObjectOutput, err error) {
	if input == nil {
		return nil, errors.New("PutFileInput is nil")
	}

	var body io.Reader
	sourceFile := strings.TrimSpace(input.SourceFile)
	if sourceFile != "" {
		fd, _err := os.Open(sourceFile)
		if _err != nil {
			err = _err
			return nil, err
		}
		defer func() {
			errMsg := fd.Close()
			if errMsg != nil {
				doLog(LEVEL_WARN, "Failed to close file with reason: %v", errMsg)
			}
		}()

		stat, _err := fd.Stat()
		if _err != nil {
			err = _err
			return nil, err
		}
		fileReaderWrapper := &fileReaderWrapper{filePath: sourceFile}
		fileReaderWrapper.reader = fd
		if input.ContentLength > 0 {
			if input.ContentLength > stat.Size() {
				input.ContentLength = stat.Size()
			}
			fileReaderWrapper.totalCount = input.ContentLength
		} else {
			fileReaderWrapper.totalCount = stat.Size()
		}
		body = fileReaderWrapper
	}

	_input := &PutObjectInput{}
	_input.PutObjectBasicInput = input.PutObjectBasicInput
	_input.Body = body

	if OSSClient.isGetContentType(_input) {
		_input.ContentType = OSSClient.getContentType(_input, sourceFile)
	}

	output = &PutObjectOutput{}
	err = OSSClient.doActionWithBucketAndKey("PutFile", HTTP_PUT, _input.Bucket, _input.Key, _input, output, extensions)
	if err != nil {
		output = nil
	} else {
		ParsePutObjectOutput(output)
	}
	return
}

// CopyObject creates a copy for an existing object.
//
// You can use this API to create a copy for an object in a specified bucket.
func (OSSClient OSSClient) CopyObject(input *CopyObjectInput, extensions ...extensionOptions) (output *CopyObjectOutput, err error) {
	if input == nil {
		return nil, errors.New("CopyObjectInput is nil")
	}

	if strings.TrimSpace(input.CopySourceBucket) == "" {
		return nil, errors.New("Source bucket is empty")
	}
	if strings.TrimSpace(input.CopySourceKey) == "" {
		return nil, errors.New("Source key is empty")
	}

	output = &CopyObjectOutput{}
	err = OSSClient.doActionWithBucketAndKey("CopyObject", HTTP_PUT, input.Bucket, input.Key, input, output, extensions)
	if err != nil {
		output = nil
	} else {
		ParseCopyObjectOutput(output)
	}
	return
}

func (OSSClient OSSClient) AppendObject(input *AppendObjectInput, extensions ...extensionOptions) (output *AppendObjectOutput, err error) {
	if input == nil {
		return nil, errors.New("AppendObjectInput is nil")
	}
	if input.ContentType == "" && input.Key != "" {
		if contentType, ok := mimeTypes[strings.ToLower(input.Key[strings.LastIndex(input.Key, ".")+1:])]; ok {
			input.ContentType = contentType
		}
	}
	var repeatable bool
	if input.Body != nil {
		if _, ok := input.Body.(*strings.Reader); !ok {
			repeatable = false
		}
		if input.ContentLength > 0 {
			input.Body = &readerWrapper{reader: input.Body, totalCount: input.ContentLength}
		}
	}
	fmt.Printf("re:%s", repeatable)
	if repeatable {
		err = OSSClient.doActionWithBucketAndKey("AppendObject", HTTP_PUT, input.Bucket, input.Key, input, output, extensions)
	} else {
		err = OSSClient.doActionWithBucketAndKey("AppendObject", HTTP_PUT, input.Bucket, input.Key, input, output, extensions)
	}
	if err != nil {
		output = nil
	} else {
		if err = ParseAppendObjectOutput(output); err != nil {
			output = nil
		}
	}
	return
}

func (OSSClient OSSClient) ModifyObject(input *ModifyObjectInput, extensions ...extensionOptions) (output *ModifyObjectOutput, err error) {
	if input == nil {
		return nil, errors.New("ModifyObjectInput is nil")
	}

	output = &ModifyObjectOutput{}
	var repeatable bool
	if input.Body != nil {
		if _, ok := input.Body.(*strings.Reader); !ok {
			repeatable = false
		}
		if input.ContentLength > 0 {
			input.Body = &readerWrapper{reader: input.Body, totalCount: input.ContentLength}
		}
	}
	if repeatable {
		err = OSSClient.doActionWithBucketAndKey("ModifyObject", HTTP_PUT, input.Bucket, input.Key, input, output, extensions)
	} else {
		err = OSSClient.doActionWithBucketAndKeyUnRepeatable("ModifyObject", HTTP_PUT, input.Bucket, input.Key, input, output, extensions)
	}
	if err != nil {
		output = nil
	} else {
		ParseModifyObjectOutput(output)
	}
	return
}

func (OSSClient OSSClient) RenameFile(input *RenameFileInput, extensions ...extensionOptions) (output *RenameFileOutput, err error) {
	if input == nil {
		return nil, errors.New("RenameFileInput is nil")
	}

	output = &RenameFileOutput{}
	err = OSSClient.doActionWithBucketAndKey("RenameFile", HTTP_POST, input.Bucket, input.Key, input, output, extensions)
	if err != nil {
		output = nil
	}
	return
}

func (OSSClient OSSClient) RenameFolder(input *RenameFolderInput, extensions ...extensionOptions) (output *RenameFolderOutput, err error) {
	if input == nil {
		return nil, errors.New("RenameFolderInput is nil")
	}

	if !strings.HasSuffix(input.Key, "/") {
		input.Key += "/"
	}
	if !strings.HasSuffix(input.NewObjectKey, "/") {
		input.NewObjectKey += "/"
	}
	output = &RenameFolderOutput{}
	err = OSSClient.doActionWithBucketAndKey("RenameFolder", HTTP_POST, input.Bucket, input.Key, input, output, extensions)
	if err != nil {
		output = nil
	}
	return
}
