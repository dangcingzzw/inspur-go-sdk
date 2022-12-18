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

package main

import (
	"OSS"
	"examples"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	endpoint   = "https://your-endpoint"
	ak         = "*** Provide your Access Key ***"
	sk         = "*** Provide your Secret Key ***"
	bucketName = "bucket-test"
	objectKey  = "object-test"
	location   = "yourbucketlocation"
)

var OSSClient *OSS.OSSClient

func getOSSClient() *OSS.OSSClient {
	var err error
	if OSSClient == nil {
		OSSClient, err = OSS.New(ak, sk, endpoint)
		if err != nil {
			panic(err)
		}
	}
	return OSSClient
}

func createBucket() {
	input := &OSS.CreateBucketInput{}
	input.Bucket = bucketName
	input.StorageClass = OSS.StorageClassWarm
	input.ACL = OSS.AclPublicRead
	output, err := getOSSClient().CreateBucket(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func listBuckets() {
	input := &OSS.ListBucketsInput{}
	input.QueryLocation = true
	output, err := getOSSClient().ListBuckets(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		fmt.Printf("Owner.DisplayName:%s, Owner.ID:%s\n", output.Owner.DisplayName, output.Owner.ID)
		for index, val := range output.Buckets {
			fmt.Printf("Bucket[%d]-Name:%s,CreationDate:%s,Location:%s\n", index, val.Name, val.CreationDate, val.Location)
		}
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func setBucketStoragePolicy() {
	input := &OSS.SetBucketStoragePolicyInput{}
	input.Bucket = bucketName
	input.StorageClass = OSS.StorageClassCold
	output, err := getOSSClient().SetBucketStoragePolicy(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func getBucketStoragePolicy() {
	output, err := getOSSClient().GetBucketStoragePolicy(bucketName)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		fmt.Printf("StorageClass:%s\n", output.StorageClass)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func deleteBucket() {
	output, err := getOSSClient().DeleteBucket(bucketName)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func listObjects() {
	input := &OSS.ListObjectsInput{}
	input.Bucket = bucketName
	input.MaxKeys = 10
	//	input.Prefix = "src/"
	output, err := getOSSClient().ListObjects(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		for index, val := range output.Contents {
			fmt.Printf("Content[%d]-OwnerId:%s, OwnerName:%s, ETag:%s, Key:%s, LastModified:%s, Size:%d, StorageClass:%s\n",
				index, val.Owner.ID, val.Owner.DisplayName, val.ETag, val.Key, val.LastModified, val.Size, val.StorageClass)
		}
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func listVersions() {
	input := &OSS.ListVersionsInput{}
	input.Bucket = bucketName
	input.MaxKeys = 10
	//	input.Prefix = "src/"
	output, err := getOSSClient().ListVersions(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		for index, val := range output.Versions {
			fmt.Printf("Version[%d]-OwnerId:%s, OwnerName:%s, ETag:%s, Key:%s, VersionId:%s, LastModified:%s, Size:%d, StorageClass:%s\n",
				index, val.Owner.ID, val.Owner.DisplayName, val.ETag, val.Key, val.VersionId, val.LastModified, val.Size, val.StorageClass)
		}
		for index, val := range output.DeleteMarkers {
			fmt.Printf("DeleteMarker[%d]-OwnerId:%s, OwnerName:%s, Key:%s, VersionId:%s, LastModified:%s, StorageClass:%s\n",
				index, val.Owner.ID, val.Owner.DisplayName, val.Key, val.VersionId, val.LastModified, val.StorageClass)
		}
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func setBucketQuota() {
	input := &OSS.SetBucketQuotaInput{}
	input.Bucket = bucketName
	input.Quota = 1024 * 1024 * 1024
	output, err := getOSSClient().SetBucketQuota(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func getBucketQuota() {
	output, err := getOSSClient().GetBucketQuota(bucketName)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		fmt.Printf("Quota:%d\n", output.Quota)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func getBucketStorageInfo() {
	output, err := getOSSClient().GetBucketStorageInfo(bucketName)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		fmt.Printf("Size:%d, ObjectNumber:%d\n", output.Size, output.ObjectNumber)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func getBucketLocation() {
	output, err := getOSSClient().GetBucketLocation(bucketName)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		fmt.Printf("Location:%s\n", output.Location)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func setBucketAcl() {
	input := &OSS.SetBucketAclInput{}
	input.Bucket = bucketName
	//		input.ACL = oss.AclPublicRead
	input.Owner.ID = "ownerid"
	var grants [3]OSS.Grant
	grants[0].Grantee.Type = OSS.GranteeGroup
	grants[0].Grantee.URI = OSS.GroupAuthenticatedUsers
	grants[0].Permission = OSS.PermissionRead

	grants[1].Grantee.Type = OSS.GranteeUser
	grants[1].Grantee.ID = "userid"
	grants[1].Permission = OSS.PermissionWrite

	grants[2].Grantee.Type = OSS.GranteeUser
	grants[2].Grantee.ID = "userid"
	grants[2].Grantee.DisplayName = "username"
	grants[2].Permission = OSS.PermissionRead
	input.Grants = grants[0:3]
	output, err := getOSSClient().SetBucketAcl(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func getBucketAcl() {
	output, err := getOSSClient().GetBucketAcl(bucketName)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		fmt.Printf("Owner.DisplayName:%s, Owner.ID:%s\n", output.Owner.DisplayName, output.Owner.ID)
		for index, grant := range output.Grants {
			fmt.Printf("Grant[%d]-Type:%s, ID:%s, URI:%s, Permission:%s\n", index, grant.Grantee.Type, grant.Grantee.ID, grant.Grantee.URI, grant.Permission)
		}
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func setBucketPolicy() {
	input := &OSS.SetBucketPolicyInput{}
	input.Bucket = bucketName
	input.Policy = "your policy"
	output, err := getOSSClient().SetBucketPolicy(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func getBucketPolicy() {
	output, err := getOSSClient().GetBucketPolicy(bucketName)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		fmt.Printf("Policy:%s\n", output.Policy)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func deleteBucketPolicy() {
	output, err := getOSSClient().DeleteBucketPolicy(bucketName)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func setBucketCors() {
	input := &OSS.SetBucketCorsInput{}
	input.Bucket = bucketName

	var corsRules [2]OSS.CorsRule
	corsRule0 := OSS.CorsRule{}
	corsRule0.ID = "rule1"
	corsRule0.AllowedOrigin = []string{"http://www.a.com", "http://www.b.com"}
	corsRule0.AllowedMethod = []string{"GET", "PUT", "POST", "HEAD"}
	corsRule0.AllowedHeader = []string{"header1", "header2"}
	corsRule0.MaxAgeSeconds = 100
	corsRule0.ExposeHeader = []string{"oss-1", "oss-2"}
	corsRules[0] = corsRule0
	corsRule1 := OSS.CorsRule{}

	corsRule1.ID = "rule2"
	corsRule1.AllowedOrigin = []string{"http://www.c.com", "http://www.d.com"}
	corsRule1.AllowedMethod = []string{"GET", "PUT", "POST", "HEAD"}
	corsRule1.AllowedHeader = []string{"header3", "header4"}
	corsRule1.MaxAgeSeconds = 50
	corsRule1.ExposeHeader = []string{"oss-3", "oss-4"}
	corsRules[1] = corsRule1
	input.CorsRules = corsRules[:]
	output, err := getOSSClient().SetBucketCors(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func getBucketCors() {
	output, err := getOSSClient().GetBucketCors(bucketName)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		for _, corsRule := range output.CorsRules {
			fmt.Printf("ID:%s, AllowedOrigin:%s, AllowedMethod:%s, AllowedHeader:%s, MaxAgeSeconds:%d, ExposeHeader:%s\n",
				corsRule.ID, strings.Join(corsRule.AllowedOrigin, "|"), strings.Join(corsRule.AllowedMethod, "|"),
				strings.Join(corsRule.AllowedHeader, "|"), corsRule.MaxAgeSeconds, strings.Join(corsRule.ExposeHeader, "|"))
		}
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func deleteBucketCors() {
	output, err := getOSSClient().DeleteBucketCors(bucketName)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func setBucketVersioning() {
	input := &OSS.SetBucketVersioningInput{}
	input.Bucket = bucketName
	input.Status = OSS.VersioningStatusEnabled
	output, err := getOSSClient().SetBucketVersioning(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func getBucketVersioning() {
	output, err := getOSSClient().GetBucketVersioning(bucketName)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		fmt.Printf("Status:%s\n", output.Status)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func headBucket() {
	output, err := getOSSClient().HeadBucket(bucketName)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
		} else {
			fmt.Println(err)
		}
	}
}

func getBucketMetadata() {
	input := &OSS.GetBucketMetadataInput{}
	input.Bucket = bucketName
	output, err := getOSSClient().GetBucketMetadata(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		fmt.Printf("StorageClass:%s\n", output.StorageClass)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Printf("StatusCode:%d\n", OSSError.StatusCode)
		} else {
			fmt.Println(err)
		}
	}
}

func setBucketLoggingConfiguration() {
	input := &OSS.SetBucketLoggingConfigurationInput{}
	input.Bucket = bucketName
	input.TargetBucket = "target-bucket"
	input.TargetPrefix = "prefix"
	var grants [3]OSS.Grant
	grants[0].Grantee.Type = OSS.GranteeGroup
	grants[0].Grantee.URI = OSS.GroupAuthenticatedUsers
	grants[0].Permission = OSS.PermissionRead

	grants[1].Grantee.Type = OSS.GranteeUser
	grants[1].Grantee.ID = "userid"
	grants[1].Permission = OSS.PermissionWrite

	grants[2].Grantee.Type = OSS.GranteeUser
	grants[2].Grantee.ID = "userid"
	grants[2].Grantee.DisplayName = "username"
	grants[2].Permission = OSS.PermissionRead
	input.TargetGrants = grants[0:3]
	output, err := getOSSClient().SetBucketLoggingConfiguration(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func getBucketLoggingConfiguration() {
	output, err := getOSSClient().GetBucketLoggingConfiguration(bucketName)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		fmt.Printf("TargetBucket:%s, TargetPrefix:%s\n", output.TargetBucket, output.TargetPrefix)
		for index, grant := range output.TargetGrants {
			fmt.Printf("Grant[%d]-Type:%s, ID:%s, URI:%s, Permission:%s\n", index, grant.Grantee.Type, grant.Grantee.ID, grant.Grantee.URI, grant.Permission)
		}
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func setBucketWebsiteConfiguration() {
	input := &OSS.SetBucketWebsiteConfigurationInput{}
	input.Bucket = bucketName
	//	input.RedirectAllRequestsTo.HostName = "www.a.com"
	//	input.RedirectAllRequestsTo.Protocol = oss.ProtocolHttp
	input.IndexDocument.Suffix = "suffix"
	input.ErrorDocument.Key = "key"

	var routingRules [2]OSS.RoutingRule
	routingRule0 := OSS.RoutingRule{}

	routingRule0.Redirect.HostName = "www.a.com"
	routingRule0.Redirect.Protocol = OSS.ProtocolHttp
	routingRule0.Redirect.ReplaceKeyPrefixWith = "prefix"
	routingRule0.Redirect.HttpRedirectCode = "304"
	routingRules[0] = routingRule0

	routingRule1 := OSS.RoutingRule{}

	routingRule1.Redirect.HostName = "www.b.com"
	routingRule1.Redirect.Protocol = OSS.ProtocolHttps
	routingRule1.Redirect.ReplaceKeyWith = "replaceKey"
	routingRule1.Redirect.HttpRedirectCode = "304"

	routingRule1.Condition.HttpErrorCodeReturnedEquals = "404"
	routingRule1.Condition.KeyPrefixEquals = "prefix"

	routingRules[1] = routingRule1

	input.RoutingRules = routingRules[:]
	output, err := getOSSClient().SetBucketWebsiteConfiguration(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func getBucketWebsiteConfiguration() {
	output, err := getOSSClient().GetBucketWebsiteConfiguration(bucketName)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		fmt.Printf("RedirectAllRequestsTo.HostName:%s,RedirectAllRequestsTo.Protocol:%s\n", output.RedirectAllRequestsTo.HostName, output.RedirectAllRequestsTo.Protocol)
		fmt.Printf("Suffix:%s\n", output.IndexDocument.Suffix)
		fmt.Printf("Key:%s\n", output.ErrorDocument.Key)
		for index, routingRule := range output.RoutingRules {
			fmt.Printf("Condition[%d]-KeyPrefixEquals:%s, HttpErrorCodeReturnedEquals:%s\n", index, routingRule.Condition.KeyPrefixEquals, routingRule.Condition.HttpErrorCodeReturnedEquals)
			fmt.Printf("Redirect[%d]-Protocol:%s, HostName:%s, ReplaceKeyPrefixWith:%s, HttpRedirectCode:%s\n",
				index, routingRule.Redirect.Protocol, routingRule.Redirect.HostName, routingRule.Redirect.ReplaceKeyPrefixWith, routingRule.Redirect.HttpRedirectCode)
		}
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func deleteBucketWebsiteConfiguration() {
	output, err := getOSSClient().DeleteBucketWebsiteConfiguration(bucketName)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func setBucketLifecycleConfiguration() {
	input := &OSS.SetBucketLifecycleConfigurationInput{}
	input.Bucket = bucketName

	var lifecycleRules [2]OSS.LifecycleRule
	lifecycleRule0 := OSS.LifecycleRule{}
	lifecycleRule0.ID = "rule0"
	lifecycleRule0.Prefix = "prefix0"
	lifecycleRule0.Status = OSS.RuleStatusEnabled

	var transitions [2]OSS.Transition
	transitions[0] = OSS.Transition{}
	transitions[0].Days = 30
	transitions[0].StorageClass = OSS.StorageClassWarm

	transitions[1] = OSS.Transition{}
	transitions[1].Days = 60
	transitions[1].StorageClass = OSS.StorageClassCold
	lifecycleRule0.Transitions = transitions[:]

	lifecycleRule0.Expiration.Days = 100
	lifecycleRule0.NoncurrentVersionExpiration.NoncurrentDays = 20

	lifecycleRules[0] = lifecycleRule0

	lifecycleRule1 := OSS.LifecycleRule{}
	lifecycleRule1.Status = OSS.RuleStatusEnabled
	lifecycleRule1.ID = "rule1"
	lifecycleRule1.Prefix = "prefix1"
	lifecycleRule1.Expiration.Date = time.Now().Add(time.Duration(24) * time.Hour)

	var noncurrentTransitions [2]OSS.NoncurrentVersionTransition
	noncurrentTransitions[0] = OSS.NoncurrentVersionTransition{}
	noncurrentTransitions[0].NoncurrentDays = 30
	noncurrentTransitions[0].StorageClass = OSS.StorageClassWarm

	noncurrentTransitions[1] = OSS.NoncurrentVersionTransition{}
	noncurrentTransitions[1].NoncurrentDays = 60
	noncurrentTransitions[1].StorageClass = OSS.StorageClassCold
	lifecycleRule1.NoncurrentVersionTransitions = noncurrentTransitions[:]
	lifecycleRules[1] = lifecycleRule1

	input.LifecycleRules = lifecycleRules[:]

	output, err := getOSSClient().SetBucketLifecycleConfiguration(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func getBucketLifecycleConfiguration() {
	output, err := getOSSClient().GetBucketLifecycleConfiguration(bucketName)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		for index, lifecycleRule := range output.LifecycleRules {
			fmt.Printf("LifecycleRule[%d]:\n", index)
			fmt.Printf("ID:%s, Prefix:%s, Status:%s\n", lifecycleRule.ID, lifecycleRule.Prefix, lifecycleRule.Status)

			date := ""
			for _, transition := range lifecycleRule.Transitions {
				if !transition.Date.IsZero() {
					date = transition.Date.String()
				}
				fmt.Printf("transition.StorageClass:%s, Transition.Date:%s, Transition.Days:%d\n", transition.StorageClass, date, transition.Days)
			}

			date = ""
			if !lifecycleRule.Expiration.Date.IsZero() {
				date = lifecycleRule.Expiration.Date.String()
			}
			fmt.Printf("Expiration.Date:%s, Expiration.Days:%d\n", lifecycleRule.Expiration.Date, lifecycleRule.Expiration.Days)

			for _, noncurrentVersionTransition := range lifecycleRule.NoncurrentVersionTransitions {
				fmt.Printf("noncurrentVersionTransition.StorageClass:%s, noncurrentVersionTransition.NoncurrentDays:%d\n",
					noncurrentVersionTransition.StorageClass, noncurrentVersionTransition.NoncurrentDays)
			}
			fmt.Printf("NoncurrentVersionExpiration.NoncurrentDays:%d\n", lifecycleRule.NoncurrentVersionExpiration.NoncurrentDays)
		}
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func deleteBucketLifecycleConfiguration() {
	output, err := getOSSClient().DeleteBucketLifecycleConfiguration(bucketName)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func setBucketTagging() {
	input := &OSS.SetBucketTaggingInput{}
	input.Bucket = bucketName

	var tags [2]OSS.Tag
	tags[0] = OSS.Tag{Key: "key0", Value: "value0"}
	tags[1] = OSS.Tag{Key: "key1", Value: "value1"}
	input.Tags = tags[:]
	output, err := getOSSClient().SetBucketTagging(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func getBucketTagging() {
	output, err := getOSSClient().GetBucketTagging(bucketName)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		for index, tag := range output.Tags {
			fmt.Printf("Tag[%d]-Key:%s, Value:%s\n", index, tag.Key, tag.Value)
		}
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func deleteBucketTagging() {
	output, err := getOSSClient().DeleteBucketTagging(bucketName)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func setBucketNotification() {
	input := &OSS.SetBucketNotificationInput{}
	input.Bucket = bucketName
	var topicConfigurations [1]OSS.TopicConfiguration
	topicConfigurations[0] = OSS.TopicConfiguration{}
	topicConfigurations[0].ID = "001"
	topicConfigurations[0].Topic = "your topic"
	topicConfigurations[0].Events = []OSS.EventType{OSS.ObjectCreatedAll}

	var filterRules [2]OSS.FilterRule

	filterRules[0] = OSS.FilterRule{Name: "prefix", Value: "smn"}
	filterRules[1] = OSS.FilterRule{Name: "suffix", Value: ".jpg"}
	topicConfigurations[0].FilterRules = filterRules[:]

	input.TopicConfigurations = topicConfigurations[:]
	output, err := getOSSClient().SetBucketNotification(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func getBucketNotification() {
	output, err := getOSSClient().GetBucketNotification(bucketName)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		for index, topicConfiguration := range output.TopicConfigurations {
			fmt.Printf("TopicConfiguration[%d]\n", index)
			fmt.Printf("ID:%s, Topic:%s, Events:%v\n", topicConfiguration.ID, topicConfiguration.Topic, topicConfiguration.Events)
		}
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func setBucketEncryption() {
	input := &OSS.SetBucketEncryptionInput{}
	input.Bucket = bucketName
	input.SSEAlgorithm = OSS.DEFAULT_SSE_KMS_ENCRYPTION

	output, err := getOSSClient().SetBucketEncryption(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func getBucketEncryption() {
	output, err := getOSSClient().GetBucketEncryption(bucketName)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		if output.KMSMasterKeyID == "" {
			fmt.Printf("KMSMasterKeyID: default master key.\n")
		} else {
			fmt.Printf("KMSMasterKeyID: %s\n", output.KMSMasterKeyID)
		}
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func deleteBucketEncryption() {
	output, err := getOSSClient().DeleteBucketEncryption(bucketName)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func listMultipartUploads() {
	input := &OSS.ListMultipartUploadsInput{}
	input.Bucket = bucketName
	input.MaxUploads = 10
	output, err := getOSSClient().ListMultipartUploads(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		for index, upload := range output.Uploads {
			fmt.Printf("Upload[%d]-OwnerId:%s, OwnerName:%s, UploadId:%s, Key:%s, Initiated:%s,StorageClass:%s\n",
				index, upload.Owner.ID, upload.Owner.DisplayName, upload.UploadId, upload.Key, upload.Initiated, upload.StorageClass)
		}
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func deleteObject() {
	input := &OSS.DeleteObjectInput{}
	input.Bucket = bucketName
	input.Key = objectKey
	output, err := getOSSClient().DeleteObject(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		fmt.Printf("VersionId:%s, DeleteMarker:%v\n", output.VersionId, output.DeleteMarker)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func deleteObjects() {
	input := &OSS.DeleteObjectsInput{}
	input.Bucket = bucketName
	var objects [3]OSS.ObjectToDelete
	objects[0] = OSS.ObjectToDelete{Key: "key1"}
	objects[1] = OSS.ObjectToDelete{Key: "key2"}
	objects[2] = OSS.ObjectToDelete{Key: "key3"}

	input.Objects = objects[:]
	output, err := getOSSClient().DeleteObjects(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		for index, deleted := range output.Deleteds {
			fmt.Printf("Deleted[%d]-Key:%s, VersionId:%s\n", index, deleted.Key, deleted.VersionId)
		}
		for index, err := range output.Errors {
			fmt.Printf("Error[%d]-Key:%s, Code:%s\n", index, err.Key, err.Code)
		}
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func setObjectAcl() {
	input := &OSS.SetObjectAclInput{}
	input.Bucket = bucketName
	input.Key = objectKey
	// input.ACL = oss.AclPublicRead
	input.Owner.ID = "ownerid"
	var grants [3]OSS.Grant
	grants[0].Grantee.Type = OSS.GranteeGroup
	grants[0].Grantee.URI = OSS.GroupAuthenticatedUsers
	grants[0].Permission = OSS.PermissionRead

	grants[1].Grantee.Type = OSS.GranteeUser
	grants[1].Grantee.ID = "userid"
	grants[1].Permission = OSS.PermissionWrite

	grants[2].Grantee.Type = OSS.GranteeUser
	grants[2].Grantee.ID = "userid"
	grants[2].Grantee.DisplayName = "username"
	grants[2].Permission = OSS.PermissionRead
	input.Grants = grants[0:3]
	output, err := getOSSClient().SetObjectAcl(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func getObjectAcl() {
	input := &OSS.GetObjectAclInput{}
	input.Bucket = bucketName
	input.Key = objectKey
	output, err := getOSSClient().GetObjectAcl(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		fmt.Printf("Owner.DisplayName:%s, Owner.ID:%s\n", output.Owner.DisplayName, output.Owner.ID)
		for index, grant := range output.Grants {
			fmt.Printf("Grant[%d]-Type:%s, ID:%s, URI:%s, Permission:%s\n", index, grant.Grantee.Type, grant.Grantee.ID, grant.Grantee.URI, grant.Permission)
		}
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func restoreObject() {
	input := &OSS.RestoreObjectInput{}
	input.Bucket = bucketName
	input.Key = objectKey
	input.Days = 1
	input.Tier = OSS.RestoreTierExpedited
	output, err := getOSSClient().RestoreObject(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func getObjectMetadata() {
	input := &OSS.GetObjectMetadataInput{}
	input.Bucket = bucketName
	input.Key = objectKey
	output, err := getOSSClient().GetObjectMetadata(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		fmt.Printf("StorageClass:%s, ETag:%s, ContentType:%s, ContentLength:%d, LastModified:%s\n",
			output.StorageClass, output.ETag, output.ContentType, output.ContentLength, output.LastModified)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Printf("StatusCode:%d\n", OSSError.StatusCode)
		} else {
			fmt.Println(err)
		}
	}
}

func copyObject() {
	input := &OSS.CopyObjectInput{}
	input.Bucket = bucketName
	input.Key = objectKey
	input.CopySourceBucket = bucketName
	input.CopySourceKey = objectKey + "-back"
	input.Metadata = map[string]string{"meta": "value"}
	input.MetadataDirective = OSS.ReplaceMetadata

	output, err := getOSSClient().CopyObject(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		fmt.Printf("ETag:%s, LastModified:%s\n",
			output.ETag, output.LastModified)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func initiateMultipartUpload() {
	input := &OSS.InitiateMultipartUploadInput{}
	input.Bucket = bucketName
	input.Key = objectKey
	input.Metadata = map[string]string{"meta": "value"}
	output, err := getOSSClient().InitiateMultipartUpload(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		fmt.Printf("Bucket:%s, Key:%s, UploadId:%s\n", output.Bucket, output.Key, output.UploadId)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func abortMultipartUpload() {
	input := &OSS.ListMultipartUploadsInput{}
	input.Bucket = bucketName
	output, err := getOSSClient().ListMultipartUploads(input)
	if err == nil {
		for _, upload := range output.Uploads {
			input := &OSS.AbortMultipartUploadInput{Bucket: bucketName}
			input.UploadId = upload.UploadId
			input.Key = upload.Key
			output, err := getOSSClient().AbortMultipartUpload(input)
			if err == nil {
				fmt.Printf("Abort uploadId[%s] successfully\n", input.UploadId)
				fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
			}
		}
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func putObject() {
	input := &OSS.PutObjectInput{}
	input.Bucket = bucketName
	input.Key = objectKey
	input.Metadata = map[string]string{"meta": "value"}
	input.Body = strings.NewReader("Hello oss")
	output, err := getOSSClient().PutObject(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		fmt.Printf("ETag:%s, StorageClass:%s\n", output.ETag, output.StorageClass)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func putFile() {
	input := &OSS.PutFileInput{}
	input.Bucket = bucketName
	input.Key = objectKey
	input.SourceFile = "localfile"
	output, err := getOSSClient().PutFile(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		fmt.Printf("ETag:%s, StorageClass:%s\n", output.ETag, output.StorageClass)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func uploadPart() {
	sourceFile := "localfile"
	var partSize int64 = 1024 * 1024 * 5
	fileInfo, statErr := os.Stat(sourceFile)
	if statErr != nil {
		panic(statErr)
	}
	partCount := fileInfo.Size() / partSize
	if fileInfo.Size()%partSize > 0 {
		partCount++
	}
	var i int64
	for i = 0; i < partCount; i++ {
		input := &OSS.UploadPartInput{}
		input.Bucket = bucketName
		input.Key = objectKey
		input.UploadId = "uploadid"
		input.PartNumber = int(i + 1)
		input.Offset = i * partSize
		if i == partCount-1 {
			input.PartSize = fileInfo.Size()
		} else {
			input.PartSize = partSize
		}
		input.SourceFile = sourceFile
		output, err := getOSSClient().UploadPart(input)
		if err == nil {
			fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
			fmt.Printf("ETag:%s\n", output.ETag)
		} else {
			if OSSError, ok := err.(OSS.OSSError); ok {
				fmt.Println(OSSError.StatusCode)
				fmt.Println(OSSError.Code)
				fmt.Println(OSSError.Message)
			} else {
				fmt.Println(err)
			}
		}
	}
}

func listParts() {
	input := &OSS.ListPartsInput{}
	input.Bucket = bucketName
	input.Key = objectKey
	input.UploadId = "uploadid"
	output, err := getOSSClient().ListParts(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		for index, part := range output.Parts {
			fmt.Printf("Part[%d]-ETag:%s, PartNumber:%d, LastModified:%s, Size:%d\n", index, part.ETag,
				part.PartNumber, part.LastModified, part.Size)
		}
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func completeMultipartUpload() {
	input := &OSS.CompleteMultipartUploadInput{}
	input.Bucket = bucketName
	input.Key = objectKey
	input.UploadId = "uploadid"
	input.Parts = []OSS.Part{
		OSS.Part{PartNumber: 1, ETag: "etag1"},
		OSS.Part{PartNumber: 2, ETag: "etag2"},
		OSS.Part{PartNumber: 3, ETag: "etag3"},
	}
	output, err := getOSSClient().CompleteMultipartUpload(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		fmt.Printf("Location:%s, Bucket:%s, Key:%s, ETag:%s\n", output.Location, output.Bucket, output.Key, output.ETag)
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func copyPart() {

	sourceBucket := "source-bucket"
	sourceKey := "source-key"
	input := &OSS.GetObjectMetadataInput{}
	input.Bucket = sourceBucket
	input.Key = sourceKey
	output, err := getOSSClient().GetObjectMetadata(input)
	if err == nil {
		objectSize := output.ContentLength
		var partSize int64 = 5 * 1024 * 1024
		partCount := objectSize / partSize
		if objectSize%partSize > 0 {
			partCount++
		}
		var i int64
		for i = 0; i < partCount; i++ {
			input := &OSS.CopyPartInput{}
			input.Bucket = bucketName
			input.Key = objectKey
			input.UploadId = "uploadid"
			input.PartNumber = int(i + 1)
			input.CopySourceBucket = sourceBucket
			input.CopySourceKey = sourceKey
			input.CopySourceRangeStart = i * partSize
			if i == partCount-1 {
				input.CopySourceRangeEnd = objectSize - 1
			} else {
				input.CopySourceRangeEnd = (i+1)*partSize - 1
			}
			output, err := getOSSClient().CopyPart(input)
			if err == nil {
				fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
				fmt.Printf("ETag:%s, PartNumber:%d\n", output.ETag, output.PartNumber)
			} else {
				if OSSError, ok := err.(OSS.OSSError); ok {
					fmt.Println(OSSError.StatusCode)
					fmt.Println(OSSError.Code)
					fmt.Println(OSSError.Message)
				} else {
					fmt.Println(err)
				}
			}
		}
	}
}

func getObject() {
	input := &OSS.GetObjectInput{}
	input.Bucket = bucketName
	input.Key = objectKey
	output, err := getOSSClient().GetObject(input)
	if err == nil {
		defer output.Body.Close()
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		fmt.Printf("StorageClass:%s, ETag:%s, ContentType:%s, ContentLength:%d, LastModified:%s\n",
			output.StorageClass, output.ETag, output.ContentType, output.ContentLength, output.LastModified)
		p := make([]byte, 1024)
		var readErr error
		var readCount int
		for {
			readCount, readErr = output.Body.Read(p)
			if readCount > 0 {
				fmt.Printf("%s", p[:readCount])
			}
			if readErr != nil {
				break
			}
		}
	} else {
		if OSSError, ok := err.(OSS.OSSError); ok {
			fmt.Println(OSSError.StatusCode)
			fmt.Println(OSSError.Code)
			fmt.Println(OSSError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func runExamples() {
	examples.RunBucketOperationsSample()
	//	examples.RunObjectOperationsSample()
	//	examples.RunDownloadSample()
	//	examples.RunCreateFolderSample()
	//	examples.RunDeleteObjectsSample()
	//	examples.RunListObjectsSample()
	//	examples.RunListVersionsSample()
	//	examples.RunListObjectsInFolderSample()
	//	examples.RunConcurrentCopyPartSample()
	//	examples.RunConcurrentDownloadObjectSample()
	//	examples.RunConcurrentUploadPartSample()
	//	examples.RunRestoreObjectSample()

	//	examples.RunSimpleMultipartUploadSample()
	//	examples.RunObjectMetaSample()
	//	examples.RunTemporarySignatureSample()
}

func main() {
	//---- init log ----
	defer OSS.CloseLog()
	OSS.InitLog("/temp/oss-SDK.log", 1024*1024*100, 5, OSS.LEVEL_WARN, false)

	//---- run examples----
	//	runExamples()

	//---- bucket related APIs ----
	//	createBucket()
	//  listBuckets()
	//	oss.FlushLog()
	//	setBucketStoragePolicy()
	//	getBucketStoragePolicy()
	//  listObjects()
	//  listVersions()
	//  listMultipartUploads()
	//	setBucketQuota()
	//	getBucketQuota()
	//	getBucketStorageInfo()
	//	getBucketLocation()
	//	setBucketAcl()
	//  getBucketAcl()
	//	setBucketPolicy()
	//  getBucketPolicy()
	//	deleteBucketPolicy()
	//  setBucketCors()
	//  getBucketCors()
	//  deleteBucketCors()
	//  setBucketVersioning()
	//  getBucketVersioning()
	//  headBucket()
	//  getBucketMetadata()
	//  setBucketLoggingConfiguration()
	//  getBucketLoggingConfiguration()
	//  setBucketWebsiteConfiguration()
	//  getBucketWebsiteConfiguration()
	//  deleteBucketWebsiteConfiguration()
	//  setBucketLifecycleConfiguration()
	//  getBucketLifecycleConfiguration()
	//  deleteBucketLifecycleConfiguration()
	//  setBucketTagging()
	//  getBucketTagging()
	//  deleteBucketTagging()
	//  setBucketNotification()
	//  getBucketNotification()
	//  setBucketEncryption()
	//  getBucketEncryption()
	//  deleteBucketEncryption()

	//---- object related APIs ----
	//  deleteObject()
	//  deleteObjects()
	//  setObjectAcl()
	//  getObjectAcl()
	//  restoreObject()
	//  copyObject()
	//  initiateMultipartUpload()
	//  uploadPart()
	//  copyPart()
	//  listParts()
	//  completeMultipartUpload()
	//  abortMultipartUpload()
	//  putObject()
	//  putFile()
	//  getObjectMetadata()
	//  getObject()

	// deleteBucket()
}
