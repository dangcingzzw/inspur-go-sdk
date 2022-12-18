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

//nolint:structcheck, unused
package oss

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// IReadCloser defines interface with function: setReadCloser
type IReadCloser interface {
	setReadCloser(body io.ReadCloser)
}

func (output *GetObjectOutput) setReadCloser(body io.ReadCloser) {
	output.Body = body
}

func setHeaders(headers map[string][]string, header string, headerValue []string, isOSS bool) {
	if isOSS {
		header = HEADER_PREFIX_OSS + header
		headers[header] = headerValue
	} else {
		header = HEADER_PREFIX + header
		headers[header] = headerValue
	}
}

func setHeadersNext(headers map[string][]string, header string, headerNext string, headerValue []string, isOSS bool) {
	if isOSS {
		headers[header] = headerValue
	} else {
		headers[headerNext] = headerValue
	}
}

// IBaseModel defines interface for base response model
type IBaseModel interface {
	setStatusCode(statusCode int)

	setRequestID(requestID string)

	setResponseHeaders(responseHeaders map[string][]string)
}

// ISerializable defines interface with function: trans
type ISerializable interface {
	trans(isOSS bool) (map[string]string, map[string][]string, interface{}, error)
}

// DefaultSerializable defines default serializable struct
type DefaultSerializable struct {
	params  map[string]string
	headers map[string][]string
	data    interface{}
}

func (input DefaultSerializable) trans(isOSS bool) (map[string]string, map[string][]string, interface{}, error) {
	return input.params, input.headers, input.data, nil
}

var defaultSerializable = &DefaultSerializable{}

func newSubResourceSerial(subResource SubResourceType) *DefaultSerializable {
	return &DefaultSerializable{map[string]string{string(subResource): ""}, nil, nil}
}

func trans(subResource SubResourceType, input interface{}) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{string(subResource): ""}
	data, err = ConvertRequestToIoReader(input)
	return
}

func (baseModel *BaseModel) setStatusCode(statusCode int) {
	baseModel.StatusCode = statusCode
}

func (baseModel *BaseModel) setRequestID(requestID string) {
	baseModel.RequestId = requestID
}

func (baseModel *BaseModel) setResponseHeaders(responseHeaders map[string][]string) {
	baseModel.ResponseHeaders = responseHeaders
}

func (input ListBucketsInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	headers = make(map[string][]string)
	if input.QueryLocation && !isOSS {
		setHeaders(headers, HEADER_LOCATION_AMZ, []string{"true"}, isOSS)
	}
	if input.BucketType != "" {
		setHeaders(headers, HEADER_BUCKET_TYPE, []string{string(input.BucketType)}, true)
	}
	return
}

func (input CreateBucketInput) prepareGrantHeaders(headers map[string][]string, isOSS bool) {
	if grantReadID := input.GrantReadId; grantReadID != "" {
		setHeaders(headers, HEADER_GRANT_READ_OSS, []string{grantReadID}, isOSS)
	}
	if grantWriteID := input.GrantWriteId; grantWriteID != "" {
		setHeaders(headers, HEADER_GRANT_WRITE_OSS, []string{grantWriteID}, isOSS)
	}
	if grantReadAcpID := input.GrantReadAcpId; grantReadAcpID != "" {
		setHeaders(headers, HEADER_GRANT_READ_ACP_OSS, []string{grantReadAcpID}, isOSS)
	}
	if grantWriteAcpID := input.GrantWriteAcpId; grantWriteAcpID != "" {
		setHeaders(headers, HEADER_GRANT_WRITE_ACP_OSS, []string{grantWriteAcpID}, isOSS)
	}
	if grantFullControlID := input.GrantFullControlId; grantFullControlID != "" {
		setHeaders(headers, HEADER_GRANT_FULL_CONTROL_OSS, []string{grantFullControlID}, isOSS)
	}
	if grantReadDeliveredID := input.GrantReadDeliveredId; grantReadDeliveredID != "" {
		setHeaders(headers, HEADER_GRANT_READ_DELIVERED_OSS, []string{grantReadDeliveredID}, true)
	}
	if grantFullControlDeliveredID := input.GrantFullControlDeliveredId; grantFullControlDeliveredID != "" {
		setHeaders(headers, HEADER_GRANT_FULL_CONTROL_DELIVERED_OSS, []string{grantFullControlDeliveredID}, true)
	}
}

func (input CreateBucketInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	headers = make(map[string][]string)
	if acl := string(input.ACL); acl != "" {
		setHeaders(headers, HEADER_ACL, []string{acl}, isOSS)
	}
	if storageClass := string(input.StorageClass); storageClass != "" {
		if !isOSS {
			if storageClass == string(StorageClassWarm) {
				storageClass = string(storageClassStandardIA)
			} else if storageClass == string(StorageClassCold) {
				storageClass = string(storageClassGlacier)
			}
		}
		setHeadersNext(headers, HEADER_STORAGE_CLASS_OSS, HEADER_STORAGE_CLASS, []string{storageClass}, isOSS)
	}
	if epid := input.Epid; epid != "" {
		setHeaders(headers, HEADER_EPID_HEADERS, []string{epid}, isOSS)
	}
	if availableZone := input.AvailableZone; availableZone != "" {
		setHeaders(headers, HEADER_AZ_REDUNDANCY, []string{availableZone}, isOSS)
	}

	input.prepareGrantHeaders(headers, isOSS)
	if input.IsFSFileInterface {
		setHeaders(headers, headerFSFileInterface, []string{"Enabled"}, true)
	}

	if location := strings.TrimSpace(input.Location); location != "" {
		input.Location = location

		xml := make([]string, 0, 3)
		xml = append(xml, "<CreateBucketConfiguration>")
		if isOSS {
			xml = append(xml, fmt.Sprintf("<Location>%s</Location>", input.Location))
		} else {
			xml = append(xml, fmt.Sprintf("<LocationConstraint>%s</LocationConstraint>", input.Location))
		}
		xml = append(xml, "</CreateBucketConfiguration>")

		data = strings.Join(xml, "")
	}
	return
}

func (input SetBucketStoragePolicyInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	xml := make([]string, 0, 1)
	if !isOSS {
		storageClass := "STANDARD"
		if input.StorageClass == StorageClassWarm {
			storageClass = string(storageClassStandardIA)
		} else if input.StorageClass == StorageClassCold {
			storageClass = string(storageClassGlacier)
		}
		params = map[string]string{string(SubResourceStoragePolicy): ""}
		xml = append(xml, fmt.Sprintf("<StoragePolicy><DefaultStorageClass>%s</DefaultStorageClass></StoragePolicy>", storageClass))
	} else {
		if input.StorageClass != StorageClassWarm && input.StorageClass != StorageClassCold {
			input.StorageClass = StorageClassStandard
		}
		params = map[string]string{string(SubResourceStorageClass): ""}
		xml = append(xml, fmt.Sprintf("<StorageClass>%s</StorageClass>", input.StorageClass))
	}
	data = strings.Join(xml, "")
	return
}

func (input ListObjsInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = make(map[string]string)
	if input.Prefix != "" {
		params["prefix"] = input.Prefix
	}
	if input.Delimiter != "" {
		params["delimiter"] = input.Delimiter
	}
	if input.MaxKeys > 0 {
		params["max-keys"] = IntToString(input.MaxKeys)
	}
	if input.EncodingType != "" {
		params["encoding-type"] = input.EncodingType
	}
	headers = make(map[string][]string)
	if origin := strings.TrimSpace(input.Origin); origin != "" {
		headers[HEADER_ORIGIN_CAMEL] = []string{origin}
	}
	if requestHeader := strings.TrimSpace(input.RequestHeader); requestHeader != "" {
		headers[HEADER_ACCESS_CONTROL_REQUEST_HEADER_CAMEL] = []string{requestHeader}
	}
	return
}

func (input ListObjectsInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params, headers, data, err = input.ListObjsInput.trans(isOSS)
	if err != nil {
		return
	}
	if input.Marker != "" {
		params["marker"] = input.Marker
	}
	return
}

func (input ListVersionsInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params, headers, data, err = input.ListObjsInput.trans(isOSS)
	if err != nil {
		return
	}
	params[string(SubResourceVersions)] = ""
	if input.KeyMarker != "" {
		params["key-marker"] = input.KeyMarker
	}
	if input.VersionIdMarker != "" {
		params["version-id-marker"] = input.VersionIdMarker
	}
	return
}

func (input ListMultipartUploadsInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{string(SubResourceUploads): ""}
	if input.Prefix != "" {
		params["prefix"] = input.Prefix
	}
	if input.Delimiter != "" {
		params["delimiter"] = input.Delimiter
	}
	if input.MaxUploads > 0 {
		params["max-uploads"] = IntToString(input.MaxUploads)
	}
	if input.KeyMarker != "" {
		params["key-marker"] = input.KeyMarker
	}
	if input.UploadIdMarker != "" {
		params["upload-id-marker"] = input.UploadIdMarker
	}
	if input.EncodingType != "" {
		params["encoding-type"] = input.EncodingType
	}
	return
}

func (input SetBucketQuotaInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	return trans(SubResourceQuota, input)
}

func (input SetBucketAclInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{string(SubResourceAcl): ""}
	headers = make(map[string][]string)

	if acl := string(input.ACL); acl != "" {
		setHeaders(headers, HEADER_ACL, []string{acl}, isOSS)
	} else {
		data, _ = convertBucketACLToXML(input.AccessControlPolicy, false, isOSS)
	}
	return
}

func (input SetBucketPolicyInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{string(SubResourcePolicy): ""}
	data = strings.NewReader(input.Policy)
	return
}

func (input SetBucketCorsInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{string(SubResourceCors): ""}
	data, md5, err := ConvertRequestToIoReaderV2(input)
	if err != nil {
		return
	}
	headers = map[string][]string{HEADER_MD5_CAMEL: {md5}}
	return
}

func (input SetBucketVersioningInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	return trans(SubResourceVersioning, input)
}

func (input SetBucketWebsiteConfigurationInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{string(SubResourceWebsite): ""}
	data, _ = ConvertWebsiteConfigurationToXml(input.BucketWebsiteConfiguration, false)
	return
}

func (input GetBucketMetadataInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	headers = make(map[string][]string)
	if origin := strings.TrimSpace(input.Origin); origin != "" {
		headers[HEADER_ORIGIN_CAMEL] = []string{origin}
	}
	if requestHeader := strings.TrimSpace(input.RequestHeader); requestHeader != "" {
		headers[HEADER_ACCESS_CONTROL_REQUEST_HEADER_CAMEL] = []string{requestHeader}
	}
	return
}

func (input SetBucketLoggingConfigurationInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{string(SubResourceLogging): ""}
	data, _ = ConvertLoggingStatusToXml(input.BucketLoggingStatus, false, isOSS)
	return
}

func (input SetBucketLifecycleConfigurationInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{string(SubResourceLifecycle): ""}
	data, md5 := ConvertLifecyleConfigurationToXml(input.BucketLifecyleConfiguration, true, isOSS)
	headers = map[string][]string{HEADER_MD5_CAMEL: {md5}}
	return
}

func (input SetBucketEncryptionInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{string(SubResourceEncryption): ""}
	data, _ = ConvertEncryptionConfigurationToXml(input.BucketEncryptionConfiguration, false, isOSS)
	return
}

func (input SetBucketTaggingInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{string(SubResourceTagging): ""}
	data, md5, err := ConvertRequestToIoReaderV2(input)
	if err != nil {
		return
	}
	headers = map[string][]string{HEADER_MD5_CAMEL: {md5}}
	return
}

func (input SetBucketNotificationInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{string(SubResourceNotification): ""}
	data, _ = ConvertNotificationToXml(input.BucketNotification, false, isOSS)
	return
}

func (input DeleteObjectInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = make(map[string]string)
	if input.VersionId != "" {
		params[PARAM_VERSION_ID] = input.VersionId
	}
	return
}

func (input DeleteObjectsInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{string(SubResourceDelete): ""}
	if strings.ToLower(input.EncodingType) == "url" {
		for index, object := range input.Objects {
			input.Objects[index].Key = url.QueryEscape(object.Key)
		}
	}
	data, md5 := convertDeleteObjectsToXML(input)
	headers = map[string][]string{HEADER_MD5_CAMEL: {md5}}
	return
}

func (input SetObjectAclInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{string(SubResourceAcl): ""}
	if input.VersionId != "" {
		params[PARAM_VERSION_ID] = input.VersionId
	}
	headers = make(map[string][]string)
	if acl := string(input.ACL); acl != "" {
		setHeaders(headers, HEADER_ACL, []string{acl}, isOSS)
	} else {
		data, _ = ConvertAclToXml(input.AccessControlPolicy, false, isOSS)
	}
	return
}

func (input GetObjectAclInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{string(SubResourceAcl): ""}
	if input.VersionId != "" {
		params[PARAM_VERSION_ID] = input.VersionId
	}
	return
}

func (input RestoreObjectInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{string(SubResourceRestore): ""}
	if input.VersionId != "" {
		params[PARAM_VERSION_ID] = input.VersionId
	}
	if !isOSS {
		data, err = ConvertRequestToIoReader(input)
	} else {
		data = ConverntOSSRestoreToXml(input)
	}
	return
}

// GetEncryption gets the Encryption field value from SseKmsHeader
func (header SseKmsHeader) GetEncryption() string {
	if header.Encryption != "" {
		return header.Encryption
	}
	if !header.isOSS {
		return DEFAULT_SSE_KMS_ENCRYPTION
	}
	return DEFAULT_SSE_KMS_ENCRYPTION_OSS
}

// GetKey gets the Key field value from SseKmsHeader
func (header SseKmsHeader) GetKey() string {
	return header.Key
}

// GetEncryption gets the Encryption field value from SseCHeader
func (header SseCHeader) GetEncryption() string {
	if header.Encryption != "" {
		return header.Encryption
	}
	return DEFAULT_SSE_C_ENCRYPTION
}

// GetKey gets the Key field value from SseCHeader
func (header SseCHeader) GetKey() string {
	return header.Key
}

// GetKeyMD5 gets the KeyMD5 field value from SseCHeader
func (header SseCHeader) GetKeyMD5() string {
	if header.KeyMD5 != "" {
		return header.KeyMD5
	}

	if ret, err := Base64Decode(header.GetKey()); err == nil {
		return Base64Md5(ret)
	}
	return ""
}

func setSseHeader(headers map[string][]string, sseHeader ISseHeader, sseCOnly bool, isOSS bool) {
	if sseHeader != nil {
		if sseCHeader, ok := sseHeader.(SseCHeader); ok {
			setHeaders(headers, HEADER_SSEC_ENCRYPTION, []string{sseCHeader.GetEncryption()}, isOSS)
			setHeaders(headers, HEADER_SSEC_KEY, []string{sseCHeader.GetKey()}, isOSS)
			setHeaders(headers, HEADER_SSEC_KEY_MD5, []string{sseCHeader.GetKeyMD5()}, isOSS)
		} else if sseKmsHeader, ok := sseHeader.(SseKmsHeader); !sseCOnly && ok {
			sseKmsHeader.isOSS = isOSS
			setHeaders(headers, HEADER_SSEKMS_ENCRYPTION, []string{sseKmsHeader.GetEncryption()}, isOSS)
			if sseKmsHeader.GetKey() != "" {
				setHeadersNext(headers, HEADER_SSEKMS_KEY_OSS, HEADER_SSEKMS_KEY_AMZ, []string{sseKmsHeader.GetKey()}, isOSS)
			}
		}
	}
}

func (input GetObjectMetadataInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = make(map[string]string)
	if input.VersionId != "" {
		params[PARAM_VERSION_ID] = input.VersionId
	}
	headers = make(map[string][]string)

	if input.Origin != "" {
		headers[HEADER_ORIGIN_CAMEL] = []string{input.Origin}
	}

	if input.RequestHeader != "" {
		headers[HEADER_ACCESS_CONTROL_REQUEST_HEADER_CAMEL] = []string{input.RequestHeader}
	}
	setSseHeader(headers, input.SseHeader, true, isOSS)
	return
}

func (input SetObjectMetadataInput) prepareContentHeaders(headers map[string][]string) {
	if input.ContentDisposition != "" {
		headers[HEADER_CONTENT_DISPOSITION_CAMEL] = []string{input.ContentDisposition}
	}
	if input.ContentEncoding != "" {
		headers[HEADER_CONTENT_ENCODING_CAMEL] = []string{input.ContentEncoding}
	}
	if input.ContentLanguage != "" {
		headers[HEADER_CONTENT_LANGUAGE_CAMEL] = []string{input.ContentLanguage}
	}

	if input.ContentType != "" {
		headers[HEADER_CONTENT_TYPE_CAML] = []string{input.ContentType}
	}
}

func (input SetObjectMetadataInput) prepareStorageClass(headers map[string][]string, isOSS bool) {
	if storageClass := string(input.StorageClass); storageClass != "" {
		if !isOSS {
			if storageClass == string(StorageClassWarm) {
				storageClass = string(storageClassStandardIA)
			} else if storageClass == string(StorageClassCold) {
				storageClass = string(storageClassGlacier)
			}
		}
		setHeaders(headers, HEADER_STORAGE_CLASS2, []string{storageClass}, isOSS)
	}
}

func (input SetObjectMetadataInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = make(map[string]string)
	params = map[string]string{string(SubResourceMetadata): ""}
	if input.VersionId != "" {
		params[PARAM_VERSION_ID] = input.VersionId
	}
	headers = make(map[string][]string)

	if directive := string(input.MetadataDirective); directive != "" {
		setHeaders(headers, HEADER_METADATA_DIRECTIVE, []string{string(input.MetadataDirective)}, isOSS)
	} else {
		setHeaders(headers, HEADER_METADATA_DIRECTIVE, []string{string(ReplaceNew)}, isOSS)
	}
	if input.CacheControl != "" {
		headers[HEADER_CACHE_CONTROL_CAMEL] = []string{input.CacheControl}
	}
	input.prepareContentHeaders(headers)
	if input.Expires != "" {
		headers[HEADER_EXPIRES_CAMEL] = []string{input.Expires}
	}
	if input.WebsiteRedirectLocation != "" {
		setHeaders(headers, HEADER_WEBSITE_REDIRECT_LOCATION, []string{input.WebsiteRedirectLocation}, isOSS)
	}
	input.prepareStorageClass(headers, isOSS)
	if input.Metadata != nil {
		for key, value := range input.Metadata {
			key = strings.TrimSpace(key)
			setHeadersNext(headers, HEADER_PREFIX_META_OSS+key, HEADER_PREFIX_META+key, []string{value}, isOSS)
		}
	}
	return
}

func (input GetObjectInput) prepareResponseParams(params map[string]string) {
	if input.ResponseCacheControl != "" {
		params[PARAM_RESPONSE_CACHE_CONTROL] = input.ResponseCacheControl
	}
	if input.ResponseContentDisposition != "" {
		params[PARAM_RESPONSE_CONTENT_DISPOSITION] = input.ResponseContentDisposition
	}
	if input.ResponseContentEncoding != "" {
		params[PARAM_RESPONSE_CONTENT_ENCODING] = input.ResponseContentEncoding
	}
	if input.ResponseContentLanguage != "" {
		params[PARAM_RESPONSE_CONTENT_LANGUAGE] = input.ResponseContentLanguage
	}
	if input.ResponseContentType != "" {
		params[PARAM_RESPONSE_CONTENT_TYPE] = input.ResponseContentType
	}
	if input.ResponseExpires != "" {
		params[PARAM_RESPONSE_EXPIRES] = input.ResponseExpires
	}
}

func (input GetObjectInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params, headers, data, err = input.GetObjectMetadataInput.trans(isOSS)
	if err != nil {
		return
	}
	input.prepareResponseParams(params)
	if input.ImageProcess != "" {
		params[PARAM_IMAGE_PROCESS] = input.ImageProcess
	}
	if input.RangeStart >= 0 && input.RangeEnd > input.RangeStart {
		headers[HEADER_RANGE] = []string{fmt.Sprintf("bytes=%d-%d", input.RangeStart, input.RangeEnd)}
	}

	if input.IfMatch != "" {
		headers[HEADER_IF_MATCH] = []string{input.IfMatch}
	}
	if input.IfNoneMatch != "" {
		headers[HEADER_IF_NONE_MATCH] = []string{input.IfNoneMatch}
	}
	if !input.IfModifiedSince.IsZero() {
		headers[HEADER_IF_MODIFIED_SINCE] = []string{FormatUtcToRfc1123(input.IfModifiedSince)}
	}
	if !input.IfUnmodifiedSince.IsZero() {
		headers[HEADER_IF_UNMODIFIED_SINCE] = []string{FormatUtcToRfc1123(input.IfUnmodifiedSince)}
	}
	return
}

func (input ObjectOperationInput) prepareGrantHeaders(headers map[string][]string) {
	if GrantReadID := input.GrantReadId; GrantReadID != "" {
		setHeaders(headers, HEADER_GRANT_READ_OSS, []string{GrantReadID}, true)
	}
	if GrantReadAcpID := input.GrantReadAcpId; GrantReadAcpID != "" {
		setHeaders(headers, HEADER_GRANT_READ_ACP_OSS, []string{GrantReadAcpID}, true)
	}
	if GrantWriteAcpID := input.GrantWriteAcpId; GrantWriteAcpID != "" {
		setHeaders(headers, HEADER_GRANT_WRITE_ACP_OSS, []string{GrantWriteAcpID}, true)
	}
	if GrantFullControlID := input.GrantFullControlId; GrantFullControlID != "" {
		setHeaders(headers, HEADER_GRANT_FULL_CONTROL_OSS, []string{GrantFullControlID}, true)
	}
}

func (input ObjectOperationInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	headers = make(map[string][]string)
	params = make(map[string]string)
	if acl := string(input.ACL); acl != "" {
		setHeaders(headers, HEADER_ACL, []string{acl}, isOSS)
	}
	input.prepareGrantHeaders(headers)
	if storageClass := string(input.StorageClass); storageClass != "" {
		if !isOSS {
			if storageClass == string(StorageClassWarm) {
				storageClass = string(storageClassStandardIA)
			} else if storageClass == string(StorageClassCold) {
				storageClass = string(storageClassGlacier)
			}
		}
		setHeaders(headers, HEADER_STORAGE_CLASS2, []string{storageClass}, isOSS)
	}
	if input.WebsiteRedirectLocation != "" {
		setHeaders(headers, HEADER_WEBSITE_REDIRECT_LOCATION, []string{input.WebsiteRedirectLocation}, isOSS)

	}
	setSseHeader(headers, input.SseHeader, false, isOSS)
	if input.Expires != 0 {
		setHeaders(headers, HEADER_EXPIRES, []string{Int64ToString(input.Expires)}, true)
	}
	if input.Metadata != nil {
		for key, value := range input.Metadata {
			key = strings.TrimSpace(key)
			setHeadersNext(headers, HEADER_PREFIX_META_OSS+key, HEADER_PREFIX_META+key, []string{value}, isOSS)
		}
	}
	return
}

func (input PutObjectBasicInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params, headers, data, err = input.ObjectOperationInput.trans(isOSS)
	if err != nil {
		return
	}

	if input.ContentMD5 != "" {
		headers[HEADER_MD5_CAMEL] = []string{input.ContentMD5}
	}

	if input.ContentLength > 0 {
		headers[HEADER_CONTENT_LENGTH_CAMEL] = []string{Int64ToString(input.ContentLength)}
	}
	if input.ContentType != "" {
		headers[HEADER_CONTENT_TYPE_CAML] = []string{input.ContentType}
	}
	if input.ContentEncoding != "" {
		headers[HEADER_CONTENT_ENCODING_CAMEL] = []string{input.ContentEncoding}
	}
	return
}

func (input PutObjectInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params, headers, data, err = input.PutObjectBasicInput.trans(isOSS)
	if err != nil {
		return
	}
	if input.Body != nil {
		data = input.Body
	}
	return
}

func (input AppendObjectInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params, headers, data, err = input.PutObjectBasicInput.trans(isOSS)
	if err != nil {
		return
	}
	params[string(SubResourceAppend)] = ""
	params["position"] = strconv.FormatInt(input.Position, 10)
	if input.Body != nil {
		data = input.Body
	}
	return
}

func (input ModifyObjectInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	headers = make(map[string][]string)
	params = make(map[string]string)
	params[string(SubResourceModify)] = ""
	params["position"] = strconv.FormatInt(input.Position, 10)
	if input.ContentLength > 0 {
		headers[HEADER_CONTENT_LENGTH_CAMEL] = []string{Int64ToString(input.ContentLength)}
	}

	if input.Body != nil {
		data = input.Body
	}
	return
}

func (input CopyObjectInput) prepareReplaceHeaders(headers map[string][]string) {
	if input.CacheControl != "" {
		headers[HEADER_CACHE_CONTROL] = []string{input.CacheControl}
	}
	if input.ContentDisposition != "" {
		headers[HEADER_CONTENT_DISPOSITION] = []string{input.ContentDisposition}
	}
	if input.ContentEncoding != "" {
		headers[HEADER_CONTENT_ENCODING] = []string{input.ContentEncoding}
	}
	if input.ContentLanguage != "" {
		headers[HEADER_CONTENT_LANGUAGE] = []string{input.ContentLanguage}
	}
	if input.ContentType != "" {
		headers[HEADER_CONTENT_TYPE] = []string{input.ContentType}
	}
	if input.Expires != "" {
		headers[HEADER_EXPIRES] = []string{input.Expires}
	}
}

func (input CopyObjectInput) prepareCopySourceHeaders(headers map[string][]string, isOSS bool) {
	if input.CopySourceIfMatch != "" {
		setHeaders(headers, HEADER_COPY_SOURCE_IF_MATCH, []string{input.CopySourceIfMatch}, isOSS)
	}
	if input.CopySourceIfNoneMatch != "" {
		setHeaders(headers, HEADER_COPY_SOURCE_IF_NONE_MATCH, []string{input.CopySourceIfNoneMatch}, isOSS)
	}
	if !input.CopySourceIfModifiedSince.IsZero() {
		setHeaders(headers, HEADER_COPY_SOURCE_IF_MODIFIED_SINCE, []string{FormatUtcToRfc1123(input.CopySourceIfModifiedSince)}, isOSS)
	}
	if !input.CopySourceIfUnmodifiedSince.IsZero() {
		setHeaders(headers, HEADER_COPY_SOURCE_IF_UNMODIFIED_SINCE, []string{FormatUtcToRfc1123(input.CopySourceIfUnmodifiedSince)}, isOSS)
	}
}

func (input CopyObjectInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params, headers, data, err = input.ObjectOperationInput.trans(isOSS)
	if err != nil {
		return
	}

	var copySource string
	if input.CopySourceVersionId != "" {
		copySource = fmt.Sprintf("%s/%s?versionId=%s", input.CopySourceBucket, UrlEncode(input.CopySourceKey, false), input.CopySourceVersionId)
	} else {
		copySource = fmt.Sprintf("%s/%s", input.CopySourceBucket, UrlEncode(input.CopySourceKey, false))
	}
	setHeaders(headers, HEADER_COPY_SOURCE, []string{copySource}, isOSS)

	if directive := string(input.MetadataDirective); directive != "" {
		setHeaders(headers, HEADER_METADATA_DIRECTIVE, []string{directive}, isOSS)
	}

	if input.MetadataDirective == ReplaceMetadata {
		input.prepareReplaceHeaders(headers)
	}

	input.prepareCopySourceHeaders(headers, isOSS)
	if input.SourceSseHeader != nil {
		if sseCHeader, ok := input.SourceSseHeader.(SseCHeader); ok {
			setHeaders(headers, HEADER_SSEC_COPY_SOURCE_ENCRYPTION, []string{sseCHeader.GetEncryption()}, isOSS)
			setHeaders(headers, HEADER_SSEC_COPY_SOURCE_KEY, []string{sseCHeader.GetKey()}, isOSS)
			setHeaders(headers, HEADER_SSEC_COPY_SOURCE_KEY_MD5, []string{sseCHeader.GetKeyMD5()}, isOSS)
		}
	}
	if input.SuccessActionRedirect != "" {
		headers[HEADER_SUCCESS_ACTION_REDIRECT] = []string{input.SuccessActionRedirect}
	}
	return
}

func (input AbortMultipartUploadInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{"uploadId": input.UploadId}
	return
}

func (input InitiateMultipartUploadInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params, headers, data, err = input.ObjectOperationInput.trans(isOSS)
	if err != nil {
		return
	}
	if input.ContentType != "" {
		headers[HEADER_CONTENT_TYPE_CAML] = []string{input.ContentType}
	}
	params[string(SubResourceUploads)] = ""
	if input.EncodingType != "" {
		params["encoding-type"] = input.EncodingType
	}
	return
}

func (input UploadPartInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{"uploadId": input.UploadId, "partNumber": IntToString(input.PartNumber)}
	headers = make(map[string][]string)
	setSseHeader(headers, input.SseHeader, true, isOSS)
	if input.ContentMD5 != "" {
		headers[HEADER_MD5_CAMEL] = []string{input.ContentMD5}
	}
	if input.Body != nil {
		data = input.Body
	}
	return
}

func (input CompleteMultipartUploadInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{"uploadId": input.UploadId}
	if input.EncodingType != "" {
		params["encoding-type"] = input.EncodingType
	}
	data, _ = ConvertCompleteMultipartUploadInputToXml(input, false)
	return
}

func (input ListPartsInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{"uploadId": input.UploadId}
	if input.MaxParts > 0 {
		params["max-parts"] = IntToString(input.MaxParts)
	}
	if input.PartNumberMarker > 0 {
		params["part-number-marker"] = IntToString(input.PartNumberMarker)
	}
	if input.EncodingType != "" {
		params["encoding-type"] = input.EncodingType
	}
	return
}

func (input CopyPartInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{"uploadId": input.UploadId, "partNumber": IntToString(input.PartNumber)}
	headers = make(map[string][]string, 1)
	var copySource string
	if input.CopySourceVersionId != "" {
		copySource = fmt.Sprintf("%s/%s?versionId=%s", input.CopySourceBucket, UrlEncode(input.CopySourceKey, false), input.CopySourceVersionId)
	} else {
		copySource = fmt.Sprintf("%s/%s", input.CopySourceBucket, UrlEncode(input.CopySourceKey, false))
	}
	setHeaders(headers, HEADER_COPY_SOURCE, []string{copySource}, isOSS)
	if input.CopySourceRangeStart >= 0 && input.CopySourceRangeEnd > input.CopySourceRangeStart {
		setHeaders(headers, HEADER_COPY_SOURCE_RANGE, []string{fmt.Sprintf("bytes=%d-%d", input.CopySourceRangeStart, input.CopySourceRangeEnd)}, isOSS)
	}

	setSseHeader(headers, input.SseHeader, true, isOSS)
	if input.SourceSseHeader != nil {
		if sseCHeader, ok := input.SourceSseHeader.(SseCHeader); ok {
			setHeaders(headers, HEADER_SSEC_COPY_SOURCE_ENCRYPTION, []string{sseCHeader.GetEncryption()}, isOSS)
			setHeaders(headers, HEADER_SSEC_COPY_SOURCE_KEY, []string{sseCHeader.GetKey()}, isOSS)
			setHeaders(headers, HEADER_SSEC_COPY_SOURCE_KEY_MD5, []string{sseCHeader.GetKeyMD5()}, isOSS)
		}

	}
	return
}

func (input HeadObjectInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = make(map[string]string)
	if input.VersionId != "" {
		params[PARAM_VERSION_ID] = input.VersionId
	}
	return
}

func (input SetBucketRequestPaymentInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	return trans(SubResourceRequestPayment, input)
}

type partSlice []Part

func (parts partSlice) Len() int {
	return len(parts)
}

func (parts partSlice) Less(i, j int) bool {
	return parts[i].PartNumber < parts[j].PartNumber
}

func (parts partSlice) Swap(i, j int) {
	parts[i], parts[j] = parts[j], parts[i]
}

type readerWrapper struct {
	reader      io.Reader
	mark        int64
	totalCount  int64
	readedCount int64
}

func (rw *readerWrapper) seek(offset int64, whence int) (int64, error) {
	if r, ok := rw.reader.(*strings.Reader); ok {
		return r.Seek(offset, whence)
	} else if r, ok := rw.reader.(*bytes.Reader); ok {
		return r.Seek(offset, whence)
	} else if r, ok := rw.reader.(*os.File); ok {
		return r.Seek(offset, whence)
	}
	return offset, nil
}

func (rw *readerWrapper) Read(p []byte) (n int, err error) {
	if rw.totalCount == 0 {
		return 0, io.EOF
	}
	if rw.totalCount > 0 {
		n, err = rw.reader.Read(p)
		readedOnce := int64(n)
		remainCount := rw.totalCount - rw.readedCount
		if remainCount > readedOnce {
			rw.readedCount += readedOnce
			return n, err
		}
		rw.readedCount += remainCount
		return int(remainCount), io.EOF
	}
	return rw.reader.Read(p)
}

type fileReaderWrapper struct {
	readerWrapper
	filePath string
}

func (input SetBucketFetchPolicyInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	contentType, _ := mimeTypes["json"]
	headers = make(map[string][]string, 2)
	headers[HEADER_CONTENT_TYPE] = []string{contentType}
	setHeaders(headers, headerOefMarker, []string{"yes"}, isOSS)
	data, err = convertFetchPolicyToJSON(input)
	return
}

func (input GetBucketFetchPolicyInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	headers = make(map[string][]string, 1)
	setHeaders(headers, headerOefMarker, []string{"yes"}, isOSS)
	return
}

func (input DeleteBucketFetchPolicyInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	headers = make(map[string][]string, 1)
	setHeaders(headers, headerOefMarker, []string{"yes"}, isOSS)
	return
}

func (input SetBucketFetchJobInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	contentType, _ := mimeTypes["json"]
	headers = make(map[string][]string, 2)
	headers[HEADER_CONTENT_TYPE] = []string{contentType}
	setHeaders(headers, headerOefMarker, []string{"yes"}, isOSS)
	data, err = convertFetchJobToJSON(input)
	return
}

func (input GetBucketFetchJobInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	headers = make(map[string][]string, 1)
	setHeaders(headers, headerOefMarker, []string{"yes"}, isOSS)
	return
}

func (input RenameFileInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{string(SubResourceRename): ""}
	params["name"] = input.NewObjectKey
	headers = make(map[string][]string)
	if requestPayer := string(input.RequestPayer); requestPayer != "" {
		headers[HEADER_REQUEST_PAYER] = []string{requestPayer}
	}
	return
}

func (input RenameFolderInput) trans(isOSS bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{string(SubResourceRename): ""}
	params["name"] = input.NewObjectKey
	headers = make(map[string][]string)
	if requestPayer := string(input.RequestPayer); requestPayer != "" {
		headers[HEADER_REQUEST_PAYER] = []string{requestPayer}
	}
	return
}
