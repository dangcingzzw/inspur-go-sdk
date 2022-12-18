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

package oss

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
)

// ListBucketsWithSignedUrl lists buckets with the specified signed url and signed request headers
func (OSSClient OSSClient) ListBucketsWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *ListBucketsOutput, err error) {
	output = &ListBucketsOutput{}
	err = OSSClient.doHTTPWithSignedURL("ListBuckets", HTTP_GET, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	}
	return
}

// CreateBucketWithSignedUrl creates bucket with the specified signed url and signed request headers and data
func (OSSClient OSSClient) CreateBucketWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header, data io.Reader) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = OSSClient.doHTTPWithSignedURL("CreateBucket", HTTP_PUT, signedUrl, actualSignedRequestHeaders, data, output, true)
	if err != nil {
		output = nil
	}
	return
}

// DeleteBucketWithSignedUrl deletes bucket with the specified signed url and signed request headers
func (OSSClient OSSClient) DeleteBucketWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = OSSClient.doHTTPWithSignedURL("DeleteBucket", HTTP_DELETE, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	}
	return
}

// SetBucketStoragePolicyWithSignedUrl sets bucket storage class with the specified signed url and signed request headers and data
func (OSSClient OSSClient) SetBucketStoragePolicyWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header, data io.Reader) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = OSSClient.doHTTPWithSignedURL("SetBucketStoragePolicy", HTTP_PUT, signedUrl, actualSignedRequestHeaders, data, output, true)
	if err != nil {
		output = nil
	}
	return
}

// GetBucketStoragePolicyWithSignedUrl gets bucket storage class with the specified signed url and signed request headers
func (OSSClient OSSClient) GetBucketStoragePolicyWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *GetBucketStoragePolicyOutput, err error) {
	output = &GetBucketStoragePolicyOutput{}
	err = OSSClient.doHTTPWithSignedURL("GetBucketStoragePolicy", HTTP_GET, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	}
	return
}

// ListObjectsWithSignedUrl lists objects in a bucket with the specified signed url and signed request headers
func (OSSClient OSSClient) ListObjectsWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *ListObjectsOutput, err error) {
	output = &ListObjectsOutput{}
	err = OSSClient.doHTTPWithSignedURL("ListObjects", HTTP_GET, signedUrl, actualSignedRequestHeaders, nil, output, true)
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

// ListVersionsWithSignedUrl lists versioning objects in a bucket with the specified signed url and signed request headers
func (OSSClient OSSClient) ListVersionsWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *ListVersionsOutput, err error) {
	output = &ListVersionsOutput{}
	err = OSSClient.doHTTPWithSignedURL("ListVersions", HTTP_GET, signedUrl, actualSignedRequestHeaders, nil, output, true)
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

// ListMultipartUploadsWithSignedUrl lists the multipart uploads that are initialized but not combined or aborted in a
// specified bucket with the specified signed url and signed request headers
func (OSSClient OSSClient) ListMultipartUploadsWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *ListMultipartUploadsOutput, err error) {
	output = &ListMultipartUploadsOutput{}
	err = OSSClient.doHTTPWithSignedURL("ListMultipartUploads", HTTP_GET, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	} else if output.EncodingType == "url" {
		err = decodeListMultipartUploadsOutput(output)
		if err != nil {
			doLog(LEVEL_ERROR, "Failed to get ListMultipartUploadsOutput with error: %v.", err)
			output = nil
		}
	}
	return
}

// SetBucketQuotaWithSignedUrl sets the bucket quota with the specified signed url and signed request headers and data
func (OSSClient OSSClient) SetBucketQuotaWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header, data io.Reader) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = OSSClient.doHTTPWithSignedURL("SetBucketQuota", HTTP_PUT, signedUrl, actualSignedRequestHeaders, data, output, true)
	if err != nil {
		output = nil
	}
	return
}

// GetBucketQuotaWithSignedUrl gets the bucket quota with the specified signed url and signed request headers
func (OSSClient OSSClient) GetBucketQuotaWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *GetBucketQuotaOutput, err error) {
	output = &GetBucketQuotaOutput{}
	err = OSSClient.doHTTPWithSignedURL("GetBucketQuota", HTTP_GET, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	}
	return
}

// HeadBucketWithSignedUrl checks whether a bucket exists with the specified signed url and signed request headers
func (OSSClient OSSClient) HeadBucketWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = OSSClient.doHTTPWithSignedURL("HeadBucket", HTTP_HEAD, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	}
	return
}

// HeadObjectWithSignedUrl checks whether an object exists with the specified signed url and signed request headers
func (OSSClient OSSClient) HeadObjectWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = OSSClient.doHTTPWithSignedURL("HeadObject", HTTP_HEAD, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	}
	return
}

// GetBucketMetadataWithSignedUrl gets the metadata of a bucket with the specified signed url and signed request headers
func (OSSClient OSSClient) GetBucketMetadataWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *GetBucketMetadataOutput, err error) {
	output = &GetBucketMetadataOutput{}
	err = OSSClient.doHTTPWithSignedURL("GetBucketMetadata", HTTP_HEAD, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	} else {
		ParseGetBucketMetadataOutput(output)
	}
	return
}

// GetBucketStorageInfoWithSignedUrl gets storage information about a bucket with the specified signed url and signed request headers
func (OSSClient OSSClient) GetBucketStorageInfoWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *GetBucketStorageInfoOutput, err error) {
	output = &GetBucketStorageInfoOutput{}
	err = OSSClient.doHTTPWithSignedURL("GetBucketStorageInfo", HTTP_GET, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	}
	return
}

// GetBucketLocationWithSignedUrl gets the location of a bucket with the specified signed url and signed request headers
func (OSSClient OSSClient) GetBucketLocationWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *GetBucketLocationOutput, err error) {
	output = &GetBucketLocationOutput{}
	err = OSSClient.doHTTPWithSignedURL("GetBucketLocation", HTTP_GET, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	}
	return
}

// SetBucketAclWithSignedUrl sets the bucket ACL with the specified signed url and signed request headers and data
func (OSSClient OSSClient) SetBucketAclWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header, data io.Reader) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = OSSClient.doHTTPWithSignedURL("SetBucketAcl", HTTP_PUT, signedUrl, actualSignedRequestHeaders, data, output, true)
	if err != nil {
		output = nil
	}
	return
}

// GetBucketAclWithSignedUrl gets the bucket ACL with the specified signed url and signed request headers
func (OSSClient OSSClient) GetBucketAclWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *GetBucketAclOutput, err error) {
	output = &GetBucketAclOutput{}
	err = OSSClient.doHTTPWithSignedURL("GetBucketAcl", HTTP_GET, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	}
	return
}

// SetBucketPolicyWithSignedUrl sets the bucket policy with the specified signed url and signed request headers and data
func (OSSClient OSSClient) SetBucketPolicyWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header, data io.Reader) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = OSSClient.doHTTPWithSignedURL("SetBucketPolicy", HTTP_PUT, signedUrl, actualSignedRequestHeaders, data, output, true)
	if err != nil {
		output = nil
	}
	return
}

// GetBucketPolicyWithSignedUrl gets the bucket policy with the specified signed url and signed request headers
func (OSSClient OSSClient) GetBucketPolicyWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *GetBucketPolicyOutput, err error) {
	output = &GetBucketPolicyOutput{}
	err = OSSClient.doHTTPWithSignedURL("GetBucketPolicy", HTTP_GET, signedUrl, actualSignedRequestHeaders, nil, output, false)
	if err != nil {
		output = nil
	}
	return
}

// DeleteBucketPolicyWithSignedUrl deletes the bucket policy with the specified signed url and signed request headers
func (OSSClient OSSClient) DeleteBucketPolicyWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = OSSClient.doHTTPWithSignedURL("DeleteBucketPolicy", HTTP_DELETE, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	}
	return
}

// SetBucketCorsWithSignedUrl sets CORS rules for a bucket with the specified signed url and signed request headers and data
func (OSSClient OSSClient) SetBucketCorsWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header, data io.Reader) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = OSSClient.doHTTPWithSignedURL("SetBucketCors", HTTP_PUT, signedUrl, actualSignedRequestHeaders, data, output, true)
	if err != nil {
		output = nil
	}
	return
}

// GetBucketCorsWithSignedUrl gets CORS rules of a bucket with the specified signed url and signed request headers
func (OSSClient OSSClient) GetBucketCorsWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *GetBucketCorsOutput, err error) {
	output = &GetBucketCorsOutput{}
	err = OSSClient.doHTTPWithSignedURL("GetBucketCors", HTTP_GET, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	}
	return
}

// DeleteBucketCorsWithSignedUrl deletes CORS rules of a bucket with the specified signed url and signed request headers
func (OSSClient OSSClient) DeleteBucketCorsWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = OSSClient.doHTTPWithSignedURL("DeleteBucketCors", HTTP_DELETE, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	}
	return
}

// SetBucketVersioningWithSignedUrl sets the versioning status for a bucket with the specified signed url and signed request headers and data
func (OSSClient OSSClient) SetBucketVersioningWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header, data io.Reader) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = OSSClient.doHTTPWithSignedURL("SetBucketVersioning", HTTP_PUT, signedUrl, actualSignedRequestHeaders, data, output, true)
	if err != nil {
		output = nil
	}
	return
}

// GetBucketVersioningWithSignedUrl gets the versioning status of a bucket with the specified signed url and signed request headers
func (OSSClient OSSClient) GetBucketVersioningWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *GetBucketVersioningOutput, err error) {
	output = &GetBucketVersioningOutput{}
	err = OSSClient.doHTTPWithSignedURL("GetBucketVersioning", HTTP_GET, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	}
	return
}

// SetBucketWebsiteConfigurationWithSignedUrl sets website hosting for a bucket with the specified signed url and signed request headers and data
func (OSSClient OSSClient) SetBucketWebsiteConfigurationWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header, data io.Reader) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = OSSClient.doHTTPWithSignedURL("SetBucketWebsiteConfiguration", HTTP_PUT, signedUrl, actualSignedRequestHeaders, data, output, true)
	if err != nil {
		output = nil
	}
	return
}

// GetBucketWebsiteConfigurationWithSignedUrl gets the website hosting settings of a bucket with the specified signed url and signed request headers
func (OSSClient OSSClient) GetBucketWebsiteConfigurationWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *GetBucketWebsiteConfigurationOutput, err error) {
	output = &GetBucketWebsiteConfigurationOutput{}
	err = OSSClient.doHTTPWithSignedURL("GetBucketWebsiteConfiguration", HTTP_GET, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	}
	return
}

// DeleteBucketWebsiteConfigurationWithSignedUrl deletes the website hosting settings of a bucket with the specified signed url and signed request headers
func (OSSClient OSSClient) DeleteBucketWebsiteConfigurationWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = OSSClient.doHTTPWithSignedURL("DeleteBucketWebsiteConfiguration", HTTP_DELETE, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	}
	return
}

// SetBucketLoggingConfigurationWithSignedUrl sets the bucket logging with the specified signed url and signed request headers and data
func (OSSClient OSSClient) SetBucketLoggingConfigurationWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header, data io.Reader) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = OSSClient.doHTTPWithSignedURL("SetBucketLoggingConfiguration", HTTP_PUT, signedUrl, actualSignedRequestHeaders, data, output, true)
	if err != nil {
		output = nil
	}
	return
}

// GetBucketLoggingConfigurationWithSignedUrl gets the logging settings of a bucket with the specified signed url and signed request headers
func (OSSClient OSSClient) GetBucketLoggingConfigurationWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *GetBucketLoggingConfigurationOutput, err error) {
	output = &GetBucketLoggingConfigurationOutput{}
	err = OSSClient.doHTTPWithSignedURL("GetBucketLoggingConfiguration", HTTP_GET, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	}
	return
}

// SetBucketLifecycleConfigurationWithSignedUrl sets lifecycle rules for a bucket with the specified signed url and signed request headers and data
func (OSSClient OSSClient) SetBucketLifecycleConfigurationWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header, data io.Reader) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = OSSClient.doHTTPWithSignedURL("SetBucketLifecycleConfiguration", HTTP_PUT, signedUrl, actualSignedRequestHeaders, data, output, true)
	if err != nil {
		output = nil
	}
	return
}

// GetBucketLifecycleConfigurationWithSignedUrl gets lifecycle rules of a bucket with the specified signed url and signed request headers
func (OSSClient OSSClient) GetBucketLifecycleConfigurationWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *GetBucketLifecycleConfigurationOutput, err error) {
	output = &GetBucketLifecycleConfigurationOutput{}
	err = OSSClient.doHTTPWithSignedURL("GetBucketLifecycleConfiguration", HTTP_GET, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	}
	return
}

// DeleteBucketLifecycleConfigurationWithSignedUrl deletes lifecycle rules of a bucket with the specified signed url and signed request headers
func (OSSClient OSSClient) DeleteBucketLifecycleConfigurationWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = OSSClient.doHTTPWithSignedURL("DeleteBucketLifecycleConfiguration", HTTP_DELETE, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	}
	return
}

// SetBucketTaggingWithSignedUrl sets bucket tags with the specified signed url and signed request headers and data
func (OSSClient OSSClient) SetBucketTaggingWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header, data io.Reader) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = OSSClient.doHTTPWithSignedURL("SetBucketTagging", HTTP_PUT, signedUrl, actualSignedRequestHeaders, data, output, true)
	if err != nil {
		output = nil
	}
	return
}

// GetBucketTaggingWithSignedUrl gets bucket tags with the specified signed url and signed request headers
func (OSSClient OSSClient) GetBucketTaggingWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *GetBucketTaggingOutput, err error) {
	output = &GetBucketTaggingOutput{}
	err = OSSClient.doHTTPWithSignedURL("GetBucketTagging", HTTP_GET, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	}
	return
}

// DeleteBucketTaggingWithSignedUrl deletes bucket tags with the specified signed url and signed request headers
func (OSSClient OSSClient) DeleteBucketTaggingWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = OSSClient.doHTTPWithSignedURL("DeleteBucketTagging", HTTP_DELETE, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	}
	return
}

// SetBucketNotificationWithSignedUrl sets event notification for a bucket with the specified signed url and signed request headers and data
func (OSSClient OSSClient) SetBucketNotificationWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header, data io.Reader) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = OSSClient.doHTTPWithSignedURL("SetBucketNotification", HTTP_PUT, signedUrl, actualSignedRequestHeaders, data, output, true)
	if err != nil {
		output = nil
	}
	return
}

// GetBucketNotificationWithSignedUrl gets event notification settings of a bucket with the specified signed url and signed request headers
func (OSSClient OSSClient) GetBucketNotificationWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *GetBucketNotificationOutput, err error) {
	output = &GetBucketNotificationOutput{}
	err = OSSClient.doHTTPWithSignedURL("GetBucketNotification", HTTP_GET, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	}
	return
}

// DeleteObjectWithSignedUrl deletes an object with the specified signed url and signed request headers
func (OSSClient OSSClient) DeleteObjectWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *DeleteObjectOutput, err error) {
	output = &DeleteObjectOutput{}
	err = OSSClient.doHTTPWithSignedURL("DeleteObject", HTTP_DELETE, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	} else {
		ParseDeleteObjectOutput(output)
	}
	return
}

// DeleteObjectsWithSignedUrl deletes objects in a batch with the specified signed url and signed request headers and data
func (OSSClient OSSClient) DeleteObjectsWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header, data io.Reader) (output *DeleteObjectsOutput, err error) {
	output = &DeleteObjectsOutput{}
	err = OSSClient.doHTTPWithSignedURL("DeleteObjects", HTTP_POST, signedUrl, actualSignedRequestHeaders, data, output, true)
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

// SetObjectAclWithSignedUrl sets ACL for an object with the specified signed url and signed request headers and data
func (OSSClient OSSClient) SetObjectAclWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header, data io.Reader) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = OSSClient.doHTTPWithSignedURL("SetObjectAcl", HTTP_PUT, signedUrl, actualSignedRequestHeaders, data, output, true)
	if err != nil {
		output = nil
	}
	return
}

// GetObjectAclWithSignedUrl gets the ACL of an object with the specified signed url and signed request headers
func (OSSClient OSSClient) GetObjectAclWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *GetObjectAclOutput, err error) {
	output = &GetObjectAclOutput{}
	err = OSSClient.doHTTPWithSignedURL("GetObjectAcl", HTTP_GET, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	} else {
		if versionID, ok := output.ResponseHeaders[HEADER_VERSION_ID]; ok {
			output.VersionId = versionID[0]
		}
	}
	return
}

// RestoreObjectWithSignedUrl restores an object with the specified signed url and signed request headers and data
func (OSSClient OSSClient) RestoreObjectWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header, data io.Reader) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = OSSClient.doHTTPWithSignedURL("RestoreObject", HTTP_POST, signedUrl, actualSignedRequestHeaders, data, output, true)
	if err != nil {
		output = nil
	}
	return
}

// GetObjectMetadataWithSignedUrl gets object metadata with the specified signed url and signed request headers
func (OSSClient OSSClient) GetObjectMetadataWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *GetObjectMetadataOutput, err error) {
	output = &GetObjectMetadataOutput{}
	err = OSSClient.doHTTPWithSignedURL("GetObjectMetadata", HTTP_HEAD, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	} else {
		ParseGetObjectMetadataOutput(output)
	}
	return
}

// GetObjectWithSignedUrl downloads object with the specified signed url and signed request headers
func (OSSClient OSSClient) GetObjectWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *GetObjectOutput, err error) {
	output = &GetObjectOutput{}
	err = OSSClient.doHTTPWithSignedURL("GetObject", HTTP_GET, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	} else {
		ParseGetObjectOutput(output)
	}
	return
}

// PutObjectWithSignedUrl uploads an object to the specified bucket with the specified signed url and signed request headers and data
func (OSSClient OSSClient) PutObjectWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header, data io.Reader) (output *PutObjectOutput, err error) {
	output = &PutObjectOutput{}
	err = OSSClient.doHTTPWithSignedURL("PutObject", HTTP_PUT, signedUrl, actualSignedRequestHeaders, data, output, true)
	if err != nil {
		output = nil
	} else {
		ParsePutObjectOutput(output)
	}
	return
}

// PutFileWithSignedUrl uploads a file to the specified bucket with the specified signed url and signed request headers and sourceFile path
func (OSSClient OSSClient) PutFileWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header, sourceFile string) (output *PutObjectOutput, err error) {
	var data io.Reader
	sourceFile = strings.TrimSpace(sourceFile)
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

		var contentLength int64
		if value, ok := actualSignedRequestHeaders[HEADER_CONTENT_LENGTH_CAMEL]; ok {
			contentLength = StringToInt64(value[0], -1)
		} else if value, ok := actualSignedRequestHeaders[HEADER_CONTENT_LENGTH]; ok {
			contentLength = StringToInt64(value[0], -1)
		} else {
			contentLength = stat.Size()
		}
		if contentLength > stat.Size() {
			return nil, errors.New("ContentLength is larger than fileSize")
		}
		fileReaderWrapper.totalCount = contentLength
		data = fileReaderWrapper
	}

	output = &PutObjectOutput{}
	err = OSSClient.doHTTPWithSignedURL("PutObject", HTTP_PUT, signedUrl, actualSignedRequestHeaders, data, output, true)
	if err != nil {
		output = nil
	} else {
		ParsePutObjectOutput(output)
	}
	return
}

// CopyObjectWithSignedUrl creates a copy for an existing object with the specified signed url and signed request headers
func (OSSClient OSSClient) CopyObjectWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *CopyObjectOutput, err error) {
	output = &CopyObjectOutput{}
	err = OSSClient.doHTTPWithSignedURL("CopyObject", HTTP_PUT, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	} else {
		ParseCopyObjectOutput(output)
	}
	return
}

// AbortMultipartUploadWithSignedUrl aborts a multipart upload in a specified bucket by using the multipart upload ID with the specified signed url and signed request headers
func (OSSClient OSSClient) AbortMultipartUploadWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = OSSClient.doHTTPWithSignedURL("AbortMultipartUpload", HTTP_DELETE, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	}
	return
}

// InitiateMultipartUploadWithSignedUrl initializes a multipart upload with the specified signed url and signed request headers
func (OSSClient OSSClient) InitiateMultipartUploadWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *InitiateMultipartUploadOutput, err error) {
	output = &InitiateMultipartUploadOutput{}
	err = OSSClient.doHTTPWithSignedURL("InitiateMultipartUpload", HTTP_POST, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	} else {
		ParseInitiateMultipartUploadOutput(output)
		if output.EncodingType == "url" {
			err = decodeInitiateMultipartUploadOutput(output)
			if err != nil {
				doLog(LEVEL_ERROR, "Failed to get InitiateMultipartUploadOutput with error: %v.", err)
				output = nil
			}
		}
	}
	return
}

// UploadPartWithSignedUrl uploads a part to a specified bucket by using a specified multipart upload ID
// with the specified signed url and signed request headers and data
func (OSSClient OSSClient) UploadPartWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header, data io.Reader) (output *UploadPartOutput, err error) {
	output = &UploadPartOutput{}
	err = OSSClient.doHTTPWithSignedURL("UploadPart", HTTP_PUT, signedUrl, actualSignedRequestHeaders, data, output, true)
	if err != nil {
		output = nil
	} else {
		ParseUploadPartOutput(output)
	}
	return
}

// CompleteMultipartUploadWithSignedUrl combines the uploaded parts in a specified bucket by using the multipart upload ID
// with the specified signed url and signed request headers and data
func (OSSClient OSSClient) CompleteMultipartUploadWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header, data io.Reader) (output *CompleteMultipartUploadOutput, err error) {
	output = &CompleteMultipartUploadOutput{}
	err = OSSClient.doHTTPWithSignedURL("CompleteMultipartUpload", HTTP_POST, signedUrl, actualSignedRequestHeaders, data, output, true)
	if err != nil {
		output = nil
	} else {
		ParseCompleteMultipartUploadOutput(output)
		if output.EncodingType == "url" {
			err = decodeCompleteMultipartUploadOutput(output)
			if err != nil {
				doLog(LEVEL_ERROR, "Failed to get CompleteMultipartUploadOutput with error: %v.", err)
				output = nil
			}
		}
	}
	return
}

// ListPartsWithSignedUrl lists the uploaded parts in a bucket by using the multipart upload ID with the specified signed url and signed request headers
func (OSSClient OSSClient) ListPartsWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *ListPartsOutput, err error) {
	output = &ListPartsOutput{}
	err = OSSClient.doHTTPWithSignedURL("ListParts", HTTP_GET, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	} else if output.EncodingType == "url" {
		err = decodeListPartsOutput(output)
		if err != nil {
			doLog(LEVEL_ERROR, "Failed to get ListPartsOutput with error: %v.", err)
			output = nil
		}
	}
	return
}

// CopyPartWithSignedUrl copy a part to a specified bucket by using a specified multipart upload ID with the specified signed url and signed request headers
func (OSSClient OSSClient) CopyPartWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *CopyPartOutput, err error) {
	output = &CopyPartOutput{}
	err = OSSClient.doHTTPWithSignedURL("CopyPart", HTTP_PUT, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	} else {
		ParseCopyPartOutput(output)
	}
	return
}

// SetBucketRequestPaymentWithSignedUrl sets requester-pays setting for a bucket with the specified signed url and signed request headers and data
func (OSSClient OSSClient) SetBucketRequestPaymentWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header, data io.Reader) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = OSSClient.doHTTPWithSignedURL("SetBucketRequestPayment", HTTP_PUT, signedUrl, actualSignedRequestHeaders, data, output, true)
	if err != nil {
		output = nil
	}
	return
}

// GetBucketRequestPaymentWithSignedUrl gets requester-pays setting of a bucket with the specified signed url and signed request headers
func (OSSClient OSSClient) GetBucketRequestPaymentWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *GetBucketRequestPaymentOutput, err error) {
	output = &GetBucketRequestPaymentOutput{}
	err = OSSClient.doHTTPWithSignedURL("GetBucketRequestPayment", HTTP_GET, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	}
	return
}

// SetBucketEncryptionWithSignedURL sets bucket encryption setting for a bucket with the specified signed url and signed request headers and data
func (OSSClient OSSClient) SetBucketEncryptionWithSignedURL(signedURL string, actualSignedRequestHeaders http.Header, data io.Reader) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = OSSClient.doHTTPWithSignedURL("SetBucketEncryption", HTTP_PUT, signedURL, actualSignedRequestHeaders, data, output, true)
	if err != nil {
		output = nil
	}
	return
}

// GetBucketEncryptionWithSignedURL gets bucket encryption setting of a bucket with the specified signed url and signed request headers
func (OSSClient OSSClient) GetBucketEncryptionWithSignedURL(signedURL string, actualSignedRequestHeaders http.Header) (output *GetBucketEncryptionOutput, err error) {
	output = &GetBucketEncryptionOutput{}
	err = OSSClient.doHTTPWithSignedURL("GetBucketEncryption", HTTP_GET, signedURL, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	}
	return
}

// DeleteBucketEncryptionWithSignedURL deletes bucket encryption setting of a bucket with the specified signed url and signed request headers
func (OSSClient OSSClient) DeleteBucketEncryptionWithSignedURL(signedURL string, actualSignedRequestHeaders http.Header) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = OSSClient.doHTTPWithSignedURL("DeleteBucketEncryption", HTTP_DELETE, signedURL, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	}
	return
}

// AppendObjectWithSignedUrl uploads an object to the specified bucket with the specified signed url and signed request headers and data
func (OSSClient OSSClient) AppendObjectWithSignedURL(signedURL string, actualSignedRequestHeaders http.Header, data io.Reader) (output *AppendObjectOutput, err error) {
	output = &AppendObjectOutput{}
	err = OSSClient.doHTTPWithSignedURL("AppendObject", HTTP_POST, signedURL, actualSignedRequestHeaders, data, output, true)
	if err != nil {
		output = nil
	} else {
		if err = ParseAppendObjectOutput(output); err != nil {
			output = nil
		}
	}
	return
}

// ModifyObjectWithSignedUrl uploads an object to the specified bucket with the specified signed url and signed request headers and data
func (OSSClient OSSClient) ModifyObjectWithSignedURL(signedURL string, actualSignedRequestHeaders http.Header, data io.Reader) (output *ModifyObjectOutput, err error) {
	output = &ModifyObjectOutput{}
	err = OSSClient.doHTTPWithSignedURL("ModifyObject", HTTP_PUT, signedURL, actualSignedRequestHeaders, data, output, true)
	if err != nil {
		output = nil
	} else {
		ParseModifyObjectOutput(output)
	}
	return
}
