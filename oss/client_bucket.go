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
	"fmt"
	"strings"
)

// ListBuckets lists buckets.
//
// You can use this API to obtain the bucket list. In the list, bucket names are displayed in lexicographical order.
func (OSSClient OSSClient) ListBuckets(input *ListBucketsInput, extensions ...extensionOptions) (output *ListBucketsOutput, err error) {
	if input == nil {
		input = &ListBucketsInput{}
	}
	output = &ListBucketsOutput{}
	err = OSSClient.doActionWithoutBucket("ListBuckets", HTTP_GET, input, output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// CreateBucket creates a bucket.
//
// You can use this API to create a bucket and name it as you specify. The created bucket name must be unique in oss.
func (OSSClient OSSClient) CreateBucket(input *CreateBucketInput, extensions ...extensionOptions) (output *BaseModel, err error) {
	if input == nil {
		return nil, errors.New("CreateBucketInput is nil")
	}
	output = &BaseModel{}
	err = OSSClient.doActionWithBucket("CreateBucket", HTTP_PUT, input.Bucket, input, output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// DeleteBucket deletes a bucket.
//
// You can use this API to delete a bucket. The bucket to be deleted must be empty
// (containing no objects, noncurrent object versions, or part fragments).
func (OSSClient OSSClient) DeleteBucket(bucketName string, extensions ...extensionOptions) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = OSSClient.doActionWithBucket("DeleteBucket", HTTP_DELETE, bucketName, defaultSerializable, output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// SetBucketStoragePolicy sets bucket storage class.
//
// You can use this API to set storage class for bucket.
func (OSSClient OSSClient) SetBucketStoragePolicy(input *SetBucketStoragePolicyInput, extensions ...extensionOptions) (output *BaseModel, err error) {
	if input == nil {
		return nil, errors.New("SetBucketStoragePolicyInput is nil")
	}
	output = &BaseModel{}
	err = OSSClient.doActionWithBucket("SetBucketStoragePolicy", HTTP_PUT, input.Bucket, input, output, extensions)
	if err != nil {
		output = nil
	}
	return
}
func (OSSClient OSSClient) getBucketStoragePolicyS3(bucketName string, extensions []extensionOptions) (output *GetBucketStoragePolicyOutput, err error) {
	output = &GetBucketStoragePolicyOutput{}
	var outputS3 *getBucketStoragePolicyOutputS3
	outputS3 = &getBucketStoragePolicyOutputS3{}
	err = OSSClient.doActionWithBucket("GetBucketStoragePolicy", HTTP_GET, bucketName, newSubResourceSerial(SubResourceStoragePolicy), outputS3, extensions)
	if err != nil {
		output = nil
		return
	}
	output.BaseModel = outputS3.BaseModel
	output.StorageClass = fmt.Sprintf("%s", outputS3.StorageClass)
	return
}

func (OSSClient OSSClient) getBucketStoragePolicyOSS(bucketName string, extensions []extensionOptions) (output *GetBucketStoragePolicyOutput, err error) {
	output = &GetBucketStoragePolicyOutput{}
	var outputOSS *getBucketStoragePolicyOutputOSS
	outputOSS = &getBucketStoragePolicyOutputOSS{}
	err = OSSClient.doActionWithBucket("GetBucketStoragePolicy", HTTP_GET, bucketName, newSubResourceSerial(SubResourceStorageClass), outputOSS, extensions)
	if err != nil {
		output = nil
		return
	}
	output.BaseModel = outputOSS.BaseModel
	output.StorageClass = outputOSS.StorageClass
	return
}

// GetBucketStoragePolicy gets bucket storage class.
//
// You can use this API to obtain the storage class of a bucket.
func (OSSClient OSSClient) GetBucketStoragePolicy(bucketName string, extensions ...extensionOptions) (output *GetBucketStoragePolicyOutput, err error) {
	if OSSClient.conf.signature == SignatureOSS {
		return OSSClient.getBucketStoragePolicyOSS(bucketName, extensions)
	}
	return OSSClient.getBucketStoragePolicyS3(bucketName, extensions)
}

// SetBucketQuota sets the bucket quota.
//
// You can use this API to set the bucket quota. A bucket quota must be expressed in bytes and the maximum value is 2^63-1.
func (OSSClient OSSClient) SetBucketQuota(input *SetBucketQuotaInput, extensions ...extensionOptions) (output *BaseModel, err error) {
	if input == nil {
		return nil, errors.New("SetBucketQuotaInput is nil")
	}
	output = &BaseModel{}
	err = OSSClient.doActionWithBucket("SetBucketQuota", HTTP_PUT, input.Bucket, input, output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// GetBucketQuota gets the bucket quota.
//
// You can use this API to obtain the bucket quota. Value 0 indicates that no upper limit is set for the bucket quota.
func (OSSClient OSSClient) GetBucketQuota(bucketName string, extensions ...extensionOptions) (output *GetBucketQuotaOutput, err error) {
	output = &GetBucketQuotaOutput{}
	err = OSSClient.doActionWithBucket("GetBucketQuota", HTTP_GET, bucketName, newSubResourceSerial(SubResourceQuota), output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// HeadBucket checks whether a bucket exists.
//
// You can use this API to check whether a bucket exists.
func (OSSClient OSSClient) HeadBucket(bucketName string, extensions ...extensionOptions) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = OSSClient.doActionWithBucket("HeadBucket", HTTP_HEAD, bucketName, defaultSerializable, output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// GetBucketMetadata gets the metadata of a bucket.
//
// You can use this API to send a HEAD request to a bucket to obtain the bucket
// metadata such as the storage class and CORS rules (if set).
func (OSSClient OSSClient) GetBucketMetadata(input *GetBucketMetadataInput, extensions ...extensionOptions) (output *GetBucketMetadataOutput, err error) {
	output = &GetBucketMetadataOutput{}
	err = OSSClient.doActionWithBucket("GetBucketMetadata", HTTP_HEAD, input.Bucket, input, output, extensions)
	if err != nil {
		output = nil
	} else {
		ParseGetBucketMetadataOutput(output)
	}
	return
}

func (OSSClient OSSClient) GetBucketFSStatus(input *GetBucketFSStatusInput, extensions ...extensionOptions) (output *GetBucketFSStatusOutput, err error) {
	output = &GetBucketFSStatusOutput{}
	err = OSSClient.doActionWithBucket("GetBucketFSStatus", HTTP_HEAD, input.Bucket, input, output, extensions)
	if err != nil {
		output = nil
	} else {
		ParseGetBucketFSStatusOutput(output)
	}
	return
}

// GetBucketStorageInfo gets storage information about a bucket.
//
// You can use this API to obtain storage information about a bucket, including the
// bucket size and number of objects in the bucket.
func (OSSClient OSSClient) GetBucketStorageInfo(bucketName string, extensions ...extensionOptions) (output *GetBucketStorageInfoOutput, err error) {
	output = &GetBucketStorageInfoOutput{}
	err = OSSClient.doActionWithBucket("GetBucketStorageInfo", HTTP_GET, bucketName, newSubResourceSerial(SubResourceStorageInfo), output, extensions)
	if err != nil {
		output = nil
	}
	return
}

func (OSSClient OSSClient) getBucketLocationS3(bucketName string, extensions []extensionOptions) (output *GetBucketLocationOutput, err error) {
	output = &GetBucketLocationOutput{}
	var outputS3 *getBucketLocationOutputS3
	outputS3 = &getBucketLocationOutputS3{}
	err = OSSClient.doActionWithBucket("GetBucketLocation", HTTP_GET, bucketName, newSubResourceSerial(SubResourceLocation), outputS3, extensions)
	if err != nil {
		output = nil
	} else {
		output.BaseModel = outputS3.BaseModel
		output.Location = outputS3.Location
	}
	return
}
func (OSSClient OSSClient) getBucketLocationOSS(bucketName string, extensions []extensionOptions) (output *GetBucketLocationOutput, err error) {
	output = &GetBucketLocationOutput{}
	var outputOSS *getBucketLocationOutputOSS
	outputOSS = &getBucketLocationOutputOSS{}
	err = OSSClient.doActionWithBucket("GetBucketLocation", HTTP_GET, bucketName, newSubResourceSerial(SubResourceLocation), outputOSS, extensions)
	if err != nil {
		output = nil
	} else {
		output.BaseModel = outputOSS.BaseModel
		output.Location = outputOSS.Location
	}
	return
}

// GetBucketLocation gets the location of a bucket.
//
// You can use this API to obtain the bucket location.
func (OSSClient OSSClient) GetBucketLocation(bucketName string, extensions ...extensionOptions) (output *GetBucketLocationOutput, err error) {
	if OSSClient.conf.signature == SignatureOSS {
		return OSSClient.getBucketLocationOSS(bucketName, extensions)
	}
	return OSSClient.getBucketLocationS3(bucketName, extensions)
}

// SetBucketAcl sets the bucket ACL.
//
// You can use this API to set the ACL for a bucket.
func (OSSClient OSSClient) SetBucketAcl(input *SetBucketAclInput, extensions ...extensionOptions) (output *BaseModel, err error) {
	if input == nil {
		return nil, errors.New("SetBucketAclInput is nil")
	}
	output = &BaseModel{}
	err = OSSClient.doActionWithBucket("SetBucketAcl", HTTP_PUT, input.Bucket, input, output, extensions)
	if err != nil {
		output = nil
	}
	return
}
func (OSSClient OSSClient) getBucketACLOSS(bucketName string, extensions []extensionOptions) (output *GetBucketAclOutput, err error) {
	output = &GetBucketAclOutput{}
	var outputOSS *getBucketACLOutputOSS
	outputOSS = &getBucketACLOutputOSS{}
	err = OSSClient.doActionWithBucket("GetBucketAcl", HTTP_GET, bucketName, newSubResourceSerial(SubResourceAcl), outputOSS, extensions)
	if err != nil {
		output = nil
	} else {
		output.BaseModel = outputOSS.BaseModel
		output.Owner = outputOSS.Owner
		output.Grants = make([]Grant, 0, len(outputOSS.Grants))
		for _, valGrant := range outputOSS.Grants {
			tempOutput := Grant{}
			tempOutput.Delivered = valGrant.Delivered
			tempOutput.Permission = valGrant.Permission
			tempOutput.Grantee.DisplayName = valGrant.Grantee.DisplayName
			tempOutput.Grantee.ID = valGrant.Grantee.ID
			tempOutput.Grantee.Type = valGrant.Grantee.Type
			tempOutput.Grantee.URI = GroupAllUsers

			output.Grants = append(output.Grants, tempOutput)
		}
	}
	return
}

// GetBucketAcl gets the bucket ACL.
//
// You can use this API to obtain a bucket ACL.
func (OSSClient OSSClient) GetBucketAcl(bucketName string, extensions ...extensionOptions) (output *GetBucketAclOutput, err error) {
	output = &GetBucketAclOutput{}
	if OSSClient.conf.signature == SignatureOSS {
		return OSSClient.getBucketACLOSS(bucketName, extensions)
	}
	err = OSSClient.doActionWithBucket("GetBucketAcl", HTTP_GET, bucketName, newSubResourceSerial(SubResourceAcl), output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// SetBucketPolicy sets the bucket policy.
//
// You can use this API to set a bucket policy. If the bucket already has a policy, the
// policy will be overwritten by the one specified in this request.
func (OSSClient OSSClient) SetBucketPolicy(input *SetBucketPolicyInput, extensions ...extensionOptions) (output *BaseModel, err error) {
	if input == nil {
		return nil, errors.New("SetBucketPolicy is nil")
	}
	output = &BaseModel{}
	err = OSSClient.doActionWithBucket("SetBucketPolicy", HTTP_PUT, input.Bucket, input, output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// GetBucketPolicy gets the bucket policy.
//
// You can use this API to obtain the policy of a bucket.
func (OSSClient OSSClient) GetBucketPolicy(bucketName string, extensions ...extensionOptions) (output *GetBucketPolicyOutput, err error) {
	output = &GetBucketPolicyOutput{}
	err = OSSClient.doActionWithBucketV2("GetBucketPolicy", HTTP_GET, bucketName, newSubResourceSerial(SubResourcePolicy), output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// DeleteBucketPolicy deletes the bucket policy.
//
// You can use this API to delete the policy of a bucket.
func (OSSClient OSSClient) DeleteBucketPolicy(bucketName string, extensions ...extensionOptions) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = OSSClient.doActionWithBucket("DeleteBucketPolicy", HTTP_DELETE, bucketName, newSubResourceSerial(SubResourcePolicy), output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// SetBucketCors sets CORS rules for a bucket.
//
// You can use this API to set CORS rules for a bucket to allow client browsers to send cross-origin requests.
func (OSSClient OSSClient) SetBucketCors(input *SetBucketCorsInput, extensions ...extensionOptions) (output *BaseModel, err error) {
	if input == nil {
		return nil, errors.New("SetBucketCorsInput is nil")
	}
	output = &BaseModel{}
	err = OSSClient.doActionWithBucket("SetBucketCors", HTTP_PUT, input.Bucket, input, output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// GetBucketCors gets CORS rules of a bucket.
//
// You can use this API to obtain the CORS rules of a specified bucket.
func (OSSClient OSSClient) GetBucketCors(bucketName string, extensions ...extensionOptions) (output *GetBucketCorsOutput, err error) {
	output = &GetBucketCorsOutput{}
	err = OSSClient.doActionWithBucket("GetBucketCors", HTTP_GET, bucketName, newSubResourceSerial(SubResourceCors), output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// DeleteBucketCors deletes CORS rules of a bucket.
//
// You can use this API to delete the CORS rules of a specified bucket.
func (OSSClient OSSClient) DeleteBucketCors(bucketName string, extensions ...extensionOptions) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = OSSClient.doActionWithBucket("DeleteBucketCors", HTTP_DELETE, bucketName, newSubResourceSerial(SubResourceCors), output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// SetBucketVersioning sets the versioning status for a bucket.
//
// You can use this API to set the versioning status for a bucket.
func (OSSClient OSSClient) SetBucketVersioning(input *SetBucketVersioningInput, extensions ...extensionOptions) (output *BaseModel, err error) {
	if input == nil {
		return nil, errors.New("SetBucketVersioningInput is nil")
	}
	output = &BaseModel{}
	err = OSSClient.doActionWithBucket("SetBucketVersioning", HTTP_PUT, input.Bucket, input, output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// GetBucketVersioning gets the versioning status of a bucket.
//
// You can use this API to obtain the versioning status of a bucket.
func (OSSClient OSSClient) GetBucketVersioning(bucketName string, extensions ...extensionOptions) (output *GetBucketVersioningOutput, err error) {
	output = &GetBucketVersioningOutput{}
	err = OSSClient.doActionWithBucket("GetBucketVersioning", HTTP_GET, bucketName, newSubResourceSerial(SubResourceVersioning), output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// SetBucketWebsiteConfiguration sets website hosting for a bucket.
//
// You can use this API to set website hosting for a bucket.
func (OSSClient OSSClient) SetBucketWebsiteConfiguration(input *SetBucketWebsiteConfigurationInput, extensions ...extensionOptions) (output *BaseModel, err error) {
	if input == nil {
		return nil, errors.New("SetBucketWebsiteConfigurationInput is nil")
	}
	output = &BaseModel{}
	err = OSSClient.doActionWithBucket("SetBucketWebsiteConfiguration", HTTP_PUT, input.Bucket, input, output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// GetBucketWebsiteConfiguration gets the website hosting settings of a bucket.
//
// You can use this API to obtain the website hosting settings of a bucket.
func (OSSClient OSSClient) GetBucketWebsiteConfiguration(bucketName string, extensions ...extensionOptions) (output *GetBucketWebsiteConfigurationOutput, err error) {
	output = &GetBucketWebsiteConfigurationOutput{}
	err = OSSClient.doActionWithBucket("GetBucketWebsiteConfiguration", HTTP_GET, bucketName, newSubResourceSerial(SubResourceWebsite), output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// DeleteBucketWebsiteConfiguration deletes the website hosting settings of a bucket.
//
// You can use this API to delete the website hosting settings of a bucket.
func (OSSClient OSSClient) DeleteBucketWebsiteConfiguration(bucketName string, extensions ...extensionOptions) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = OSSClient.doActionWithBucket("DeleteBucketWebsiteConfiguration", HTTP_DELETE, bucketName, newSubResourceSerial(SubResourceWebsite), output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// SetBucketLoggingConfiguration sets the bucket logging.
//
// You can use this API to configure access logging for a bucket.
func (OSSClient OSSClient) SetBucketLoggingConfiguration(input *SetBucketLoggingConfigurationInput, extensions ...extensionOptions) (output *BaseModel, err error) {
	if input == nil {
		return nil, errors.New("SetBucketLoggingConfigurationInput is nil")
	}
	output = &BaseModel{}
	err = OSSClient.doActionWithBucket("SetBucketLoggingConfiguration", HTTP_PUT, input.Bucket, input, output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// GetBucketLoggingConfiguration gets the logging settings of a bucket.
//
// You can use this API to obtain the access logging settings of a bucket.
func (OSSClient OSSClient) GetBucketLoggingConfiguration(bucketName string, extensions ...extensionOptions) (output *GetBucketLoggingConfigurationOutput, err error) {
	output = &GetBucketLoggingConfigurationOutput{}
	err = OSSClient.doActionWithBucket("GetBucketLoggingConfiguration", HTTP_GET, bucketName, newSubResourceSerial(SubResourceLogging), output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// SetBucketLifecycleConfiguration sets lifecycle rules for a bucket.
//
// You can use this API to set lifecycle rules for a bucket, to periodically transit
// storage classes of objects and delete objects in the bucket.
func (OSSClient OSSClient) SetBucketLifecycleConfiguration(input *SetBucketLifecycleConfigurationInput, extensions ...extensionOptions) (output *BaseModel, err error) {
	if input == nil {
		return nil, errors.New("SetBucketLifecycleConfigurationInput is nil")
	}
	output = &BaseModel{}
	err = OSSClient.doActionWithBucket("SetBucketLifecycleConfiguration", HTTP_PUT, input.Bucket, input, output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// GetBucketLifecycleConfiguration gets lifecycle rules of a bucket.
//
// You can use this API to obtain the lifecycle rules of a bucket.
func (OSSClient OSSClient) GetBucketLifecycleConfiguration(bucketName string, extensions ...extensionOptions) (output *GetBucketLifecycleConfigurationOutput, err error) {
	output = &GetBucketLifecycleConfigurationOutput{}
	err = OSSClient.doActionWithBucket("GetBucketLifecycleConfiguration", HTTP_GET, bucketName, newSubResourceSerial(SubResourceLifecycle), output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// DeleteBucketLifecycleConfiguration deletes lifecycle rules of a bucket.
//
// You can use this API to delete all lifecycle rules of a bucket.
func (OSSClient OSSClient) DeleteBucketLifecycleConfiguration(bucketName string, extensions ...extensionOptions) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = OSSClient.doActionWithBucket("DeleteBucketLifecycleConfiguration", HTTP_DELETE, bucketName, newSubResourceSerial(SubResourceLifecycle), output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// SetBucketEncryption sets the default server-side encryption for a bucket.
//
// You can use this API to create or update the default server-side encryption for a bucket.
func (OSSClient OSSClient) SetBucketEncryption(input *SetBucketEncryptionInput, extensions ...extensionOptions) (output *BaseModel, err error) {
	if input == nil {
		return nil, errors.New("SetBucketEncryptionInput is nil")
	}
	output = &BaseModel{}
	err = OSSClient.doActionWithBucket("SetBucketEncryption", HTTP_PUT, input.Bucket, input, output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// GetBucketEncryption gets the encryption configuration of a bucket.
//
// You can use this API to obtain obtain the encryption configuration of a bucket.
func (OSSClient OSSClient) GetBucketEncryption(bucketName string, extensions ...extensionOptions) (output *GetBucketEncryptionOutput, err error) {
	output = &GetBucketEncryptionOutput{}
	err = OSSClient.doActionWithBucket("GetBucketEncryption", HTTP_GET, bucketName, newSubResourceSerial(SubResourceEncryption), output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// DeleteBucketEncryption deletes the encryption configuration of a bucket.
//
// You can use this API to delete the encryption configuration of a bucket.
func (OSSClient OSSClient) DeleteBucketEncryption(bucketName string, extensions ...extensionOptions) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = OSSClient.doActionWithBucket("DeleteBucketEncryption", HTTP_DELETE, bucketName, newSubResourceSerial(SubResourceEncryption), output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// SetBucketTagging sets bucket tags.
//
// You can use this API to set bucket tags.
func (OSSClient OSSClient) SetBucketTagging(input *SetBucketTaggingInput, extensions ...extensionOptions) (output *BaseModel, err error) {
	if input == nil {
		return nil, errors.New("SetBucketTaggingInput is nil")
	}
	output = &BaseModel{}
	err = OSSClient.doActionWithBucket("SetBucketTagging", HTTP_PUT, input.Bucket, input, output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// GetBucketTagging gets bucket tags.
//
// You can use this API to obtain the tags of a specified bucket.
func (OSSClient OSSClient) GetBucketTagging(bucketName string, extensions ...extensionOptions) (output *GetBucketTaggingOutput, err error) {
	output = &GetBucketTaggingOutput{}
	err = OSSClient.doActionWithBucket("GetBucketTagging", HTTP_GET, bucketName, newSubResourceSerial(SubResourceTagging), output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// DeleteBucketTagging deletes bucket tags.
//
// You can use this API to delete the tags of a specified bucket.
func (OSSClient OSSClient) DeleteBucketTagging(bucketName string, extensions ...extensionOptions) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = OSSClient.doActionWithBucket("DeleteBucketTagging", HTTP_DELETE, bucketName, newSubResourceSerial(SubResourceTagging), output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// SetBucketNotification sets event notification for a bucket.
//
// You can use this API to configure event notification for a bucket. You will be notified of all
// specified operations performed on the bucket.
func (OSSClient OSSClient) SetBucketNotification(input *SetBucketNotificationInput, extensions ...extensionOptions) (output *BaseModel, err error) {
	if input == nil {
		return nil, errors.New("SetBucketNotificationInput is nil")
	}
	output = &BaseModel{}
	err = OSSClient.doActionWithBucket("SetBucketNotification", HTTP_PUT, input.Bucket, input, output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// GetBucketNotification gets event notification settings of a bucket.
//
// You can use this API to obtain the event notification configuration of a bucket.
func (OSSClient OSSClient) GetBucketNotification(bucketName string, extensions ...extensionOptions) (output *GetBucketNotificationOutput, err error) {
	if OSSClient.conf.signature != SignatureOSS {
		return OSSClient.getBucketNotificationS3(bucketName, extensions)
	}
	output = &GetBucketNotificationOutput{}
	err = OSSClient.doActionWithBucket("GetBucketNotification", HTTP_GET, bucketName, newSubResourceSerial(SubResourceNotification), output, extensions)
	if err != nil {
		output = nil
	}
	return
}

func (OSSClient OSSClient) getBucketNotificationS3(bucketName string, extensions []extensionOptions) (output *GetBucketNotificationOutput, err error) {
	outputS3 := &getBucketNotificationOutputS3{}
	err = OSSClient.doActionWithBucket("GetBucketNotification", HTTP_GET, bucketName, newSubResourceSerial(SubResourceNotification), outputS3, extensions)
	if err != nil {
		return nil, err
	}

	output = &GetBucketNotificationOutput{}
	output.BaseModel = outputS3.BaseModel
	topicConfigurations := make([]TopicConfiguration, 0, len(outputS3.TopicConfigurations))
	for _, topicConfigurationS3 := range outputS3.TopicConfigurations {
		topicConfiguration := TopicConfiguration{}
		topicConfiguration.ID = topicConfigurationS3.ID
		topicConfiguration.Topic = topicConfigurationS3.Topic
		topicConfiguration.FilterRules = topicConfigurationS3.FilterRules

		events := make([]EventType, 0, len(topicConfigurationS3.Events))
		for _, event := range topicConfigurationS3.Events {
			events = append(events, ParseStringToEventType(event))
		}
		topicConfiguration.Events = events
		topicConfigurations = append(topicConfigurations, topicConfiguration)
	}
	output.TopicConfigurations = topicConfigurations
	return
}

// SetBucketRequestPayment sets requester-pays setting for a bucket.
func (OSSClient OSSClient) SetBucketRequestPayment(input *SetBucketRequestPaymentInput, extensions ...extensionOptions) (output *BaseModel, err error) {
	if input == nil {
		return nil, errors.New("SetBucketRequestPaymentInput is nil")
	}
	output = &BaseModel{}
	err = OSSClient.doActionWithBucket("SetBucketRequestPayment", HTTP_PUT, input.Bucket, input, output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// GetBucketRequestPayment gets requester-pays setting of a bucket.
func (OSSClient OSSClient) GetBucketRequestPayment(bucketName string, extensions ...extensionOptions) (output *GetBucketRequestPaymentOutput, err error) {
	output = &GetBucketRequestPaymentOutput{}
	err = OSSClient.doActionWithBucket("GetBucketRequestPayment", HTTP_GET, bucketName, newSubResourceSerial(SubResourceRequestPayment), output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// SetBucketFetchPolicy sets the bucket fetch policy.
//
// You can use this API to set a bucket fetch policy.
func (OSSClient OSSClient) SetBucketFetchPolicy(input *SetBucketFetchPolicyInput, extensions ...extensionOptions) (output *BaseModel, err error) {
	if input == nil {
		return nil, errors.New("SetBucketFetchPolicyInput is nil")
	}
	if strings.TrimSpace(string(input.Status)) == "" {
		return nil, errors.New("Fetch policy status is empty")
	}
	if strings.TrimSpace(input.Agency) == "" {
		return nil, errors.New("Fetch policy agency is empty")
	}
	output = &BaseModel{}
	err = OSSClient.doActionWithBucketAndKey("SetBucketFetchPolicy", HTTP_PUT, input.Bucket, string(objectKeyExtensionPolicy), input, output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// GetBucketFetchPolicy gets the bucket fetch policy.
//
// You can use this API to obtain the fetch policy of a bucket.
func (OSSClient OSSClient) GetBucketFetchPolicy(input *GetBucketFetchPolicyInput, extensions ...extensionOptions) (output *GetBucketFetchPolicyOutput, err error) {
	if input == nil {
		return nil, errors.New("GetBucketFetchPolicyInput is nil")
	}
	output = &GetBucketFetchPolicyOutput{}
	err = OSSClient.doActionWithBucketAndKeyV2("GetBucketFetchPolicy", HTTP_GET, input.Bucket, string(objectKeyExtensionPolicy), input, output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// DeleteBucketFetchPolicy deletes the bucket fetch policy.
//
// You can use this API to delete the fetch policy of a bucket.
func (OSSClient OSSClient) DeleteBucketFetchPolicy(input *DeleteBucketFetchPolicyInput, extensions ...extensionOptions) (output *BaseModel, err error) {
	if input == nil {
		return nil, errors.New("DeleteBucketFetchPolicyInput is nil")
	}
	output = &BaseModel{}
	err = OSSClient.doActionWithBucketAndKey("DeleteBucketFetchPolicy", HTTP_DELETE, input.Bucket, string(objectKeyExtensionPolicy), input, output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// SetBucketFetchJob sets the bucket fetch job.
//
// You can use this API to set a bucket fetch job.
func (OSSClient OSSClient) SetBucketFetchJob(input *SetBucketFetchJobInput, extensions ...extensionOptions) (output *SetBucketFetchJobOutput, err error) {
	if input == nil {
		return nil, errors.New("SetBucketFetchJobInput is nil")
	}
	if strings.TrimSpace(input.URL) == "" {
		return nil, errors.New("URL is empty")
	}
	output = &SetBucketFetchJobOutput{}
	err = OSSClient.doActionWithBucketAndKeyV2("SetBucketFetchJob", HTTP_POST, input.Bucket, string(objectKeyAsyncFetchJob), input, output, extensions)
	if err != nil {
		output = nil
	}
	return
}

// GetBucketFetchJob gets the bucket fetch job.
//
// You can use this API to obtain the fetch job of a bucket.
func (OSSClient OSSClient) GetBucketFetchJob(input *GetBucketFetchJobInput, extensions ...extensionOptions) (output *GetBucketFetchJobOutput, err error) {
	if input == nil {
		return nil, errors.New("GetBucketFetchJobInput is nil")
	}
	if strings.TrimSpace(input.JobID) == "" {
		return nil, errors.New("JobID is empty")
	}
	output = &GetBucketFetchJobOutput{}
	err = OSSClient.doActionWithBucketAndKeyV2("GetBucketFetchJob", HTTP_GET, input.Bucket, string(objectKeyAsyncFetchJob)+"/"+input.JobID, input, output, extensions)
	if err != nil {
		output = nil
	}
	return
}
