package all_sample

// 引入依赖包
import (
	"fmt"
	"github.com/dangcingzzw/inspur-go-sdk/OSS"
	"io/ioutil"
	"strings"
	"time"
)

var ak = "*** Provide your Access Key ***"
var sk = "*** Provide your Secret Key ***"
var endPoint = "your endPoint";
var bucketName = "your bucketName"
var objectKey = "your objectKey"
var location = "your location"

// 创建ossClient结构体
var OSSClient, _ = OSS.New(ak, sk, endPoint)

func main() {
	/**
	bucket
	*/
	createBucket()                //创建桶
	getBucketLocation()           //获取桶所在区域
	headBucket()                  //查看桶是否存在
	listBuckets()                 //桶列表
	doBucketVersioningOperation() //获取版本-设置版本Enabled-获取版本-设置版本Suspended-获取版本
	doBucketAclOperation()        //设置只读-查看-设置私有-查看
	doBucketCorsOperation()       //跨域规则-设置-获取-删除
	doBucketPolicy()              //桶策略-设置-获取-删除
	doBucketLifecycleOperation()  //生命周期-设置-获取-删除
	doBucketWebsiteOperation()    //静态网站设置-获取-删除
	doBucketVersioning()          //桶版本设置-获取
	doBucketEncryption()          //桶加密设置-获取-删除
	deleteBucket()                //删除桶(需要确保桶内没东西）
	// doBucketDomain();//桶自定义域名设置-获取-删除
	// pageListBuckets();//分页获取桶列表

	/**
	object
	*/
	createBucket()            //创建桶
	doObjectSampleOperation() //对象上传-获取元数据，下载对象,对象ACL设置-获取,修改元数据，复制对象，删除对象，批量删除对象，列举对象
	createSignedUrl()         //生成预签名链接
	doObjectAcl()             //对象ACL设置-获取
	doPartUpload()            //分片上传-初始化任务-上传分片-完成上传
	listMultipartUploads()    //分片上传任务列表
	listParts()               //分片列表
	aboutMultipartUpload()    //取消分片上传

	//appendObject() //追加对象
	//doesObjectExist();//判断对象是否存在
	//doObjectVersion()//对象版本获取-删除
}

func aboutMultipartUpload() {
	input := &OSS.AbortMultipartUploadInput{}
	input.Bucket = bucketName
	input.Key = objectKey
	input.UploadId = "uploadid"
	output, err := OSSClient.AbortMultipartUpload(input)
	if err == nil {
		fmt.Printf("RequestId:%s\n", output.RequestId)
	} else {
		if obsError, ok := err.(OSS.OSSError); ok {
			fmt.Println(obsError.Code)
			fmt.Println(obsError.Message)
		} else {
			fmt.Println(err)
		}
	}
}
func listParts() {
	input := &OSS.ListPartsInput{}
	input.Bucket = bucketName
	input.Key = objectKey
	input.UploadId = "uploadid"
	output, err := OSSClient.ListParts(input)
	if err == nil {
		fmt.Printf("RequestId:%s\n", output.RequestId)
		for index, part := range output.Parts {
			fmt.Printf("Part[%d]-ETag:%s, PartNumber:%d, LastModified:%s, Size:%d\n", index, part.ETag,
				part.PartNumber, part.LastModified, part.Size)
		}
	} else {
		if obsError, ok := err.(OSS.OSSError); ok {
			fmt.Println(obsError.Code)
			fmt.Println(obsError.Message)
		} else {
			fmt.Println(err)
		}
	}
}
func listMultipartUploads() {
	input := &OSS.ListMultipartUploadsInput{}
	input.Bucket = bucketName
	input.MaxUploads = 10
	output, err := OSSClient.ListMultipartUploads(input)
	if err == nil {
		fmt.Printf("RequestId:%s\n", output.RequestId)
		for index, upload := range output.Uploads {
			fmt.Printf("Upload[%d]-OwnerId:%s, UploadId:%s, Key:%s, Initiated:%s\n",
				index, upload.Owner.ID, upload.UploadId, upload.Key, upload.Initiated)
		}
	} else {
		if obsError, ok := err.(OSS.OSSError); ok {
			fmt.Println(obsError.Code)
			fmt.Println(obsError.Message)
		} else {
			fmt.Println(err)
		}
	}
}
func doPartUpload() {
	destBucketName := bucketName
	destObjectKey := objectKey + "-back"
	sourceBucketName := bucketName
	sourceObjectKey := objectKey
	// Claim a upload id firstly
	input := &OSS.InitiateMultipartUploadInput{}
	input.Bucket = destBucketName
	input.Key = destObjectKey
	output, err := OSSClient.InitiateMultipartUpload(input)
	if err != nil {
		panic(err)
	}
	uploadId := output.UploadId

	fmt.Printf("Claiming a new upload id %s\n", uploadId)
	fmt.Println()

	// Get size of the object
	getObjectMetadataInput := &OSS.GetObjectMetadataInput{}
	getObjectMetadataInput.Bucket = sourceBucketName
	getObjectMetadataInput.Key = sourceObjectKey
	getObjectMetadataOutput, err := OSSClient.GetObjectMetadata(getObjectMetadataInput)
	if err != nil {
		panic(err)
	}

	objectSize := getObjectMetadataOutput.ContentLength

	// Calculate how many blocks to be divided
	// 5MB
	var partSize int64 = 5 * 1024 * 1024
	partCount := int(objectSize / partSize)

	if objectSize%partSize != 0 {
		partCount++
	}

	fmt.Printf("Total parts count %d\n", partCount)
	fmt.Println()

	//  Upload multiparts by copy mode
	fmt.Println("Begin to upload multiparts to oss by copy mode")

	partChan := make(chan OSS.Part, 5)

	for i := 0; i < partCount; i++ {
		partNumber := i + 1
		rangeStart := int64(i) * partSize
		rangeEnd := rangeStart + partSize - 1
		if i+1 == partCount {
			rangeEnd = objectSize - 1
		}
		go func(start, end int64, index int) {
			copyPartInput := &OSS.CopyPartInput{}
			copyPartInput.Bucket = destBucketName
			copyPartInput.Key = destObjectKey
			copyPartInput.UploadId = uploadId
			copyPartInput.PartNumber = index
			copyPartInput.CopySourceBucket = sourceBucketName
			copyPartInput.CopySourceKey = sourceObjectKey
			copyPartInput.CopySourceRangeStart = start
			copyPartInput.CopySourceRangeEnd = end
			copyPartOutput, errMsg := OSSClient.CopyPart(copyPartInput)
			if errMsg == nil {
				fmt.Printf("%d finished\n", index)
				partChan <- OSS.Part{ETag: copyPartOutput.ETag, PartNumber: copyPartOutput.PartNumber}
			} else {
				panic(errMsg)
			}
		}(rangeStart, rangeEnd, partNumber)
	}

	parts := make([]OSS.Part, 0, partCount)

	for {
		part, ok := <-partChan
		if !ok {
			break
		}
		parts = append(parts, part)
		if len(parts) == partCount {
			close(partChan)
		}
	}

	fmt.Println()
	fmt.Println("Completing to upload multiparts")
	completeMultipartUploadInput := &OSS.CompleteMultipartUploadInput{}
	completeMultipartUploadInput.Bucket = destBucketName
	completeMultipartUploadInput.Key = destObjectKey
	completeMultipartUploadInput.UploadId = uploadId
	completeMultipartUploadInput.Parts = parts

	_, err12 := OSSClient.CompleteMultipartUpload(completeMultipartUploadInput)
	if err12 != nil {
		panic(err)
	}
	fmt.Println("Complete multiparts finished")
}
func doObjectAcl() {
	input := &OSS.PutObjectInput{}
	input.Bucket = bucketName
	input.Key = objectKey
	input.Body = strings.NewReader("Hello oss")
	//input.SourceFile = "C:\\Users\\dangcingzzw\\Pictures\\bbbzzzwsfaSRCfa.jpeg"

	_, err := OSSClient.PutObject(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Put object:%s successfully!\n", objectKey)
	fmt.Println()

	input2 := &OSS.SetObjectAclInput{}
	input2.Bucket = bucketName
	input2.Key = objectKey
	input2.ACL = OSS.AclPublicRead

	_, err2 := OSSClient.SetObjectAcl(input2)
	if err2 != nil {
		panic(err)
	}
	fmt.Println("Set object acl successfully!")
	fmt.Println()

	output2, err3 := OSSClient.GetObjectAcl(&OSS.GetObjectAclInput{Bucket: bucketName, Key: objectKey})
	if err3 != nil {
		panic(err)
	}
	fmt.Printf("Object owner - ownerId:%s, ownerName:%s\n", output2.Owner.ID, output2.Owner.DisplayName)
	for index, grant := range output2.Grants {
		fmt.Printf("Grant[%d]\n", index)
		fmt.Printf("GranteeUri:%s, GranteeId:%s, GranteeName:%s\n", grant.Grantee.URI, grant.Grantee.ID, grant.Grantee.DisplayName)
		fmt.Printf("Permission:%s\n", grant.Permission)
	}
}

func createSignedUrl() {
	input := &OSS.CreateSignedUrlInput{}
	input.Bucket = bucketName
	input.Method = OSS.HttpMethodPut
	input.Expires = 3600
	output, err := OSSClient.CreateSignedUrl(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s using temporary signature url:\n", "CreateBucket")
	fmt.Println(output.SignedUrl)

	data := strings.NewReader(fmt.Sprintf("<CreateBucketConfiguration><LocationConstraint>%s</LocationConstraint></CreateBucketConfiguration>", location))

	_, err = OSSClient.CreateBucketWithSignedUrl(output.SignedUrl, output.ActualSignedRequestHeaders, data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Create bucket:%s successfully!\n", bucketName)
	fmt.Println()
}
func appendObject() {
	input := &OSS.AppendObjectInput{}
	input.Bucket = bucketName
	input.Key = objectKey + "11a"
	// 第一次追加上传，指定传入追加上传位置为0
	input.Position = 0
	input.Body = strings.NewReader("Hello OBS")
	output, err := OSSClient.AppendObject(input)
	if err == nil {
		fmt.Printf("RequestId:%s\n", output.RequestId)
		fmt.Printf("NextAppendPosition:%d\n", output.NextAppendPosition)

		// 第二次追加上传，指定传入追加上传位置为上次追加上传返回的位置信息
		input := &OSS.AppendObjectInput{}
		input.Bucket = bucketName
		input.Key = objectKey + "11a"
		input.Position = output.NextAppendPosition
		input.Body = strings.NewReader("Hello OBS Again")
		output, err := OSSClient.AppendObject(input)
		if err == nil {
			fmt.Printf("RequestId:%s\n", output.RequestId)
			fmt.Printf("NextAppendPosition:%d\n", output.NextAppendPosition)
		} else {
			if obsError, ok := err.(OSS.OSSError); ok {
				fmt.Println(obsError.Code)
				fmt.Println(obsError.Message)
			} else {
				fmt.Println(err)
			}
		}
	} else {
		if obsError, ok := err.(OSS.OSSError); ok {
			fmt.Println(obsError.Code)
			fmt.Println(obsError.Message)
		} else {
			fmt.Println(err)
		}
	}
}
func doObjectSampleOperation() {
	fmt.Printf("put object - %s")
	//input := &OSS.PutFileInput{}
	input := &OSS.PutObjectInput{}
	input.Bucket = bucketName
	input.Key = objectKey
	input.Body = strings.NewReader("Hello oss")
	//input.SourceFile = "C:\\Users\\dangcingzzw\\Pictures\\bbbzzzwsfaSRCfa.jpeg"

	_, err := OSSClient.PutObject(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Put object:%s successfully!\n", objectKey)
	fmt.Println()

	fmt.Printf("get object - %s")
	input2 := &OSS.GetObjectInput{}
	input2.Bucket = bucketName
	input2.Key = objectKey

	output, err := OSSClient.GetObject(input2)
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

	fmt.Printf("set object acl - %s")
	input3 := &OSS.SetObjectAclInput{}
	input3.Bucket = bucketName
	input3.Key = objectKey
	input3.Owner.ID = "ownerid"
	var grants [3]OSS.Grant
	grants[0].Grantee.Type = OSS.GranteeGroup
	grants[0].Grantee.URI = OSS.GroupAuthenticatedUsers
	grants[0].Permission = OSS.PermissionRead

	grants[1].Grantee.Type = OSS.GranteeUser
	grants[1].Grantee.ID = "userid"
	grants[1].Permission = OSS.PermissionWrite

	grants[2].Grantee.Type = OSS.GranteeUser
	grants[2].Grantee.ID = "userid"
	grants[2].Permission = OSS.PermissionRead
	input3.Grants = grants[0:3]
	output3, err := OSSClient.SetObjectAcl(input3)
	if err == nil {
		fmt.Printf("RequestId:%s\n", output3.RequestId)
	} else {
		if ossError, ok := err.(OSS.OSSError); ok {
			fmt.Println(ossError.Code)
			fmt.Println(ossError.Message)
		} else {
			fmt.Println(err)
		}
	}
	fmt.Printf("get object acl - %s")
	input4 := &OSS.GetObjectAclInput{}
	input4.Bucket = bucketName
	input4.Key = objectKey
	output4, err := OSSClient.GetObjectAcl(input4)
	if err == nil {
		fmt.Printf("RequestId:%s\n", output4.RequestId)
		fmt.Printf("Owner.ID:%s\n", output4.Owner.ID)
		for index, grant := range output4.Grants {
			fmt.Printf("Grant[%d]-Type:%s, ID:%s, URI:%s, Permission:%s\n", index, grant.Grantee.Type, grant.Grantee.ID, grant.Grantee.URI, grant.Permission)
		}
	} else {
		if ossError, ok := err.(OSS.OSSError); ok {
			fmt.Println(ossError.Code)
			fmt.Println(ossError.Message)
		} else {
			fmt.Println(err)
		}
	}

	fmt.Printf("set object metadata - %s")
	input55 := &OSS.PutObjectInput{}
	input55.Bucket = bucketName
	input55.Key = objectKey
	input55.Body = strings.NewReader("Hello oss")
	// Setting object mime type
	input55.ContentType = "text/plain"
	// Setting self-defined metadata
	input55.Metadata = map[string]string{"meta1": "value1", "meta2": "value2"}
	_, err55 := OSSClient.PutObject(input55)
	if err55 != nil {
		panic(err)
	}
	fmt.Println("Set object meatdata successfully!")
	fmt.Println()

	fmt.Printf("get object metadata - %s")
	input6 := &OSS.GetObjectMetadataInput{}
	input6.Bucket = bucketName
	input6.Key = objectKey
	output6, err6 := OSSClient.GetObjectMetadata(input6)
	if err6 == nil {
		fmt.Printf("RequestId:%s\n", output6.RequestId)
		fmt.Printf("StorageClass:%s, ETag:%s, ContentType:%s, ContentLength:%d, LastModified:%s\n",
			output6.StorageClass, output6.ETag, output6.ContentType, output6.ContentLength, output6.LastModified)
	} else {
		if ossError, ok := err.(OSS.OSSError); ok {
			fmt.Println(ossError.Code)
			fmt.Println(ossError.Message)
		} else {
			fmt.Println(err)
		}
	}
	input7 := &OSS.CopyObjectInput{}
	input7.Bucket = bucketName
	input7.Key = objectKey + "-back"
	input7.CopySourceBucket = bucketName
	input7.CopySourceKey = objectKey

	_, err7 := OSSClient.CopyObject(input7)
	if err7 != nil {
		panic(err)
	}
	fmt.Println("Copy object successfully!")
	fmt.Println()

	fmt.Println("list object successfully!")
	input8 := &OSS.ListObjectsInput{}
	input8.Bucket = bucketName
	output8, err := OSSClient.ListObjects(input8)
	if err != nil {
		panic(err)
	}
	for index, val := range output8.Contents {
		fmt.Printf("Content[%d]-ETag:%s, Key:%s, Size:%d\n",
			index, val.ETag, val.Key, val.Size)
	}
	fmt.Println()

	fmt.Println("DELETE object successfully!")
	input9 := &OSS.DeleteObjectInput{}
	input9.Bucket = bucketName
	input9.Key = objectKey
	output9, err := OSSClient.DeleteObject(input9)
	if err == nil {
		fmt.Printf("RequestId:%s\n", output9.RequestId)
	} else {
		if ossError, ok := err.(OSS.OSSError); ok {
			fmt.Println(ossError.Code)
			fmt.Println(ossError.Message)
		} else {
			fmt.Println(err)
		}
	}

	fmt.Println("DELETE objectS successfully!")
	input10 := &OSS.DeleteObjectsInput{}
	input10.Bucket = bucketName
	var objects [3]OSS.ObjectToDelete
	objects[0] = OSS.ObjectToDelete{Key: "key1"}
	objects[1] = OSS.ObjectToDelete{Key: "key2"}
	objects[2] = OSS.ObjectToDelete{Key: "key3"}

	input10.Objects = objects[:]
	output10, err := OSSClient.DeleteObjects(input10)
	if err == nil {
		fmt.Printf("RequestId:%s\n", output10.RequestId)
		for index, deleted := range output10.Deleteds {
			fmt.Printf("Deleted[%d]-Key:%s, VersionId:%s\n", index, deleted.Key, deleted.VersionId)
		}
		for index, err := range output10.Errors {
			fmt.Printf("Error[%d]-Key:%s, Code:%s\n", index, err.Key, err.Code)
		}
	} else {
		if ossError, ok := err.(OSS.OSSError); ok {
			fmt.Println(ossError.Code)
			fmt.Println(ossError.Message)
		} else {
			fmt.Println(err)
		}
	}

}

func deleteBucket() {
	_, err := OSSClient.DeleteBucket(bucketName)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Delete bucket %s successfully!\n", bucketName)
	fmt.Println()
}
func doBucketEncryption() {
	fmt.Printf("set bucket encryption - %s")
	input := &OSS.SetBucketEncryptionInput{}
	input.Bucket = bucketName
	// 指定传入加密算法及对应的KMS加密密钥
	input.SSEAlgorithm = "AES256"
	input.KMSMasterKeyID = ""
	output, err := OSSClient.SetBucketEncryption(input)
	if err == nil {
		fmt.Printf("RequestId:%s\n", output.RequestId)
	} else {
		if ossError, ok := err.(OSS.OSSError); ok {
			fmt.Println(ossError.Code)
			fmt.Println(ossError.Message)
		} else {
			fmt.Println(err)
		}
	}
	fmt.Printf("get bucket encryption - %s")
	// 获取指定桶的加密配置信息
	output2, err := OSSClient.GetBucketEncryption(bucketName)
	if err == nil {
		fmt.Printf("Encryption:%s\n", output2.SSEAlgorithm)
		fmt.Printf("KeyID:%s\n", output2.KMSMasterKeyID)
	} else {
		if ossError, ok := err.(OSS.OSSError); ok {
			fmt.Println(ossError.Code)
			fmt.Println(ossError.Message)
		} else {
			fmt.Println(err)
		}
	}
	fmt.Printf("delete bucket encryption - %s")
	// 删除指定桶的加密配置信息
	output3, err := OSSClient.DeleteBucketEncryption(bucketName)
	if err == nil {
		fmt.Printf("RequestId:%s\n", output3.RequestId)
	} else {
		if ossError, ok := err.(OSS.OSSError); ok {
			fmt.Println(ossError.Code)
			fmt.Println(ossError.Message)
		} else {
			fmt.Println(err)
		}
	}
}
func doBucketPolicy() {
	fmt.Printf("set bucket policy - %s")
	input := &OSS.SetBucketPolicyInput{}
	input.Bucket = bucketName
	input.Policy = "{\"Statement\":[{\"Principal\":\"*\",\"Effect\":\"Allow\",\"Action\":\"ListBucket\",\"Resource\":\"" + bucketName + "\"}]}"
	output, err := OSSClient.SetBucketPolicy(input)
	if err == nil {
		fmt.Printf("RequestId:%s\n", output.RequestId)
	} else {
		if error, ok := err.(OSS.OSSError); ok {
			fmt.Println(error.Code)
			fmt.Println(error.Message)
		} else {
			fmt.Println(err)
		}
	}
	fmt.Printf("get bucket policy - %s")

	output2, err2 := OSSClient.GetBucketPolicy(bucketName)
	if err2 == nil {
		fmt.Printf("RequestId:%s\n", output2.RequestId)
		fmt.Printf("Policy:%s\n", output2.Policy)
	} else {
		if error, ok := err.(OSS.OSSError); ok {
			fmt.Println(error.Code)
			fmt.Println(error.Message)
		} else {
			fmt.Println(err)
		}
	}
	fmt.Printf("delete bucket policy - %s")
	output3, err3 := OSSClient.DeleteBucketPolicy(bucketName)
	if err3 == nil {
		fmt.Printf("RequestId:%s\n", output3.RequestId)
	} else {
		if err, ok := err.(OSS.OSSError); ok {
			fmt.Println(err.Code)
			fmt.Println(err.Message)
		} else {
			fmt.Println(err)
		}
	}
}
func doBucketVersioning() {
	output, err := OSSClient.GetBucketVersioning(bucketName)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Initial bucket versioning status - %s", output.Status)
	fmt.Println()

	// Enable bucket versioning
	input := &OSS.SetBucketVersioningInput{}
	input.Bucket = bucketName
	input.Status = OSS.VersioningStatusEnabled
	_, err = OSSClient.SetBucketVersioning(input)
	if err != nil {
		panic(err)
	}

	output, err = OSSClient.GetBucketVersioning(bucketName)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Current bucket versioning status - %s", output.Status)
	fmt.Println()

	// Suspend bucket versioning
	input = &OSS.SetBucketVersioningInput{}
	input.Bucket = bucketName
	input.Status = OSS.VersioningStatusSuspended
	_, err = OSSClient.SetBucketVersioning(input)
	if err != nil {
		panic(err)
	}

	output, err = OSSClient.GetBucketVersioning(bucketName)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Current bucket versioning status - %s", output.Status)
	fmt.Println()
}
func doBucketWebsiteOperation() {
	input := &OSS.SetBucketWebsiteConfigurationInput{}
	input.Bucket = bucketName
	input.IndexDocument.Suffix = "suffix"
	input.ErrorDocument.Key = "key"

	var routingRules [2]OSS.RoutingRule
	routingRule0 := OSS.RoutingRule{}

	routingRule0.Redirect.HostName = "www.a.com"
	routingRule0.Redirect.Protocol = OSS.ProtocolHttp
	routingRule0.Redirect.ReplaceKeyPrefixWith = "prefixWeb0"
	routingRule0.Redirect.HttpRedirectCode = "304"
	routingRules[0] = routingRule0

	routingRule1 := OSS.RoutingRule{}

	routingRule1.Redirect.HostName = "www.b.com"
	routingRule1.Redirect.Protocol = OSS.ProtocolHttps
	routingRule1.Redirect.ReplaceKeyWith = "replaceKey"
	routingRule1.Redirect.HttpRedirectCode = "304"

	routingRule1.Condition.HttpErrorCodeReturnedEquals = "404"
	routingRule1.Condition.KeyPrefixEquals = "prefixWeb1"

	routingRules[1] = routingRule1

	input.RoutingRules = routingRules[:]
	_, err := OSSClient.SetBucketWebsiteConfiguration(input)
	if err != nil {
		panic(err)
	}
	fmt.Println("Set bucket website successfully!")
	fmt.Println()

	output, err := OSSClient.GetBucketWebsiteConfiguration(bucketName)
	if err != nil {
		panic(err)
	}
	fmt.Printf("RedirectAllRequestsTo.HostName:%s,RedirectAllRequestsTo.Protocol:%s\n", output.RedirectAllRequestsTo.HostName, output.RedirectAllRequestsTo.Protocol)
	fmt.Printf("Suffix:%s\n", output.IndexDocument.Suffix)
	fmt.Printf("Key:%s\n", output.ErrorDocument.Key)
	for index, routingRule := range output.RoutingRules {
		fmt.Printf("Condition[%d]-KeyPrefixEquals:%s, HttpErrorCodeReturnedEquals:%s\n", index, routingRule.Condition.KeyPrefixEquals, routingRule.Condition.HttpErrorCodeReturnedEquals)
		fmt.Printf("Redirect[%d]-Protocol:%s, HostName:%s, ReplaceKeyPrefixWith:%s, HttpRedirectCode:%s\n",
			index, routingRule.Redirect.Protocol, routingRule.Redirect.HostName, routingRule.Redirect.ReplaceKeyPrefixWith, routingRule.Redirect.HttpRedirectCode)
	}
	fmt.Println()

	_, err = OSSClient.DeleteBucketWebsiteConfiguration(bucketName)
	if err != nil {
		panic(err)
	}
	fmt.Println("Delete bucket website successfully!")
	fmt.Println()
}
func headBucket() {
	exit, err := OSSClient.HeadBucket(bucketName)
	if err != nil {
		panic(err)
	}
	fmt.Printf("GET bucket EXIT:%s successfully!\n", exit)
	fmt.Println()
}
func doBucketLifecycleOperation() {
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

	_, err := OSSClient.SetBucketLifecycleConfiguration(input)
	if err != nil {
		panic(err)
	}
	fmt.Println("Set bucket lifecycle successfully!")
	fmt.Println()

	output, err := OSSClient.GetBucketLifecycleConfiguration(bucketName)
	if err != nil {
		panic(err)
	}
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
	fmt.Println()

	_, err = OSSClient.DeleteBucketLifecycleConfiguration(bucketName)
	if err != nil {
		panic(err)
	}
	fmt.Println("Delete bucket lifecycle successfully!")
	fmt.Println()
}

func createBucket() {
	input := &OSS.CreateBucketInput{}
	input.Bucket = bucketName
	_, err := OSSClient.CreateBucket(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Create bucket:%s successfully!\n", bucketName)
	fmt.Println()
}
func doBucketCorsOperation() {
	input := &OSS.SetBucketCorsInput{}
	input.Bucket = bucketName
	var corsRules [2]OSS.CorsRule
	corsRule0 := OSS.CorsRule{}
	corsRule0.ID = "rule1"
	corsRule0.AllowedOrigin = []string{"http://www.a.com", "http://www.b.com"}
	corsRule0.AllowedMethod = []string{"GET", "PUT", "POST", "HEAD"}
	corsRule0.AllowedHeader = []string{"header1", "header2"}
	corsRule0.MaxAgeSeconds = 100
	corsRule0.ExposeHeader = []string{"OSS-1", "OSS-2"}
	corsRules[0] = corsRule0
	corsRule1 := OSS.CorsRule{}

	corsRule1.ID = "rule2"
	corsRule1.AllowedOrigin = []string{"http://www.c.com", "http://www.d.com"}
	corsRule1.AllowedMethod = []string{"GET", "PUT", "POST", "HEAD"}
	corsRule1.AllowedHeader = []string{"header3", "header4"}
	corsRule1.MaxAgeSeconds = 50
	corsRule1.ExposeHeader = []string{"OSS-3", "OSS-4"}
	corsRules[1] = corsRule1
	input.CorsRules = corsRules[:]
	// Setting bucket CORS
	_, err := OSSClient.SetBucketCors(input)
	if err != nil {
		panic(err)
	}
	fmt.Println("Set bucket cors successfully!")
	fmt.Println()

	output, err := OSSClient.GetBucketCors(bucketName)
	if err != nil {
		panic(err)
	}
	for index, corsRule := range output.CorsRules {
		fmt.Printf("CorsRule[%d]\n", index)
		fmt.Printf("ID:%s, AllowedOrigin:%s, AllowedMethod:%s, AllowedHeader:%s, MaxAgeSeconds:%d, ExposeHeader:%s\n",
			corsRule.ID, strings.Join(corsRule.AllowedOrigin, "|"), strings.Join(corsRule.AllowedMethod, "|"),
			strings.Join(corsRule.AllowedHeader, "|"), corsRule.MaxAgeSeconds, strings.Join(corsRule.ExposeHeader, "|"))
	}
	fmt.Println()

	_, err = OSSClient.DeleteBucketCors(bucketName)
	if err != nil {
		panic(err)
	}
	fmt.Println("Delete bucket cors successfully!")
	fmt.Println()
}
func doBucketAclOperation() {
	input := &OSS.SetBucketAclInput{}
	input.Bucket = bucketName
	// Setting bucket ACL to public-read
	input.ACL = OSS.AclPublicRead
	_, err := OSSClient.SetBucketAcl(input)
	if err != nil {
		panic(err)
	}
	fmt.Println("Set bucket acl successfully!")
	fmt.Println()

	output, err := OSSClient.GetBucketAcl(bucketName)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Bucket owner - ownerId:%s, ownerName:%s\n", output.Owner.ID, output.Owner.DisplayName)
	for index, grant := range output.Grants {
		fmt.Printf("Grant[%d]\n", index)
		fmt.Printf("GranteeUri:%s, GranteeId:%s, GranteeName:%s\n", grant.Grantee.URI, grant.Grantee.ID, grant.Grantee.DisplayName)
		fmt.Printf("Permission:%s\n", grant.Permission)
	}
	fmt.Println()
}
func doBucketVersioningOperation() {
	output, err := OSSClient.GetBucketVersioning(bucketName)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Initial bucket versioning status - %s", output.Status)
	fmt.Println()

	// Enable bucket versioning
	input := &OSS.SetBucketVersioningInput{}
	input.Bucket = bucketName
	input.Status = OSS.VersioningStatusEnabled
	_, err = OSSClient.SetBucketVersioning(input)
	if err != nil {
		panic(err)
	}

	output, err = OSSClient.GetBucketVersioning(bucketName)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Current bucket versioning status - %s", output.Status)
	fmt.Println()

	// Suspend bucket versioning
	input = &OSS.SetBucketVersioningInput{}
	input.Bucket = bucketName
	input.Status = OSS.VersioningStatusSuspended
	_, err = OSSClient.SetBucketVersioning(input)
	if err != nil {
		panic(err)
	}

	output, err = OSSClient.GetBucketVersioning(bucketName)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Current bucket versioning status - %s", output.Status)
	fmt.Println()
}
func listBuckets() {
	input := &OSS.CreateSignedUrlInput{}
	input.Method = OSS.HttpMethodGet
	input.Expires = 3600
	output, err := OSSClient.CreateSignedUrl(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s using temporary signature url:\n", "ListBuckets")
	fmt.Println(output.SignedUrl)

	listBucketsOutput, err := OSSClient.ListBucketsWithSignedUrl(output.SignedUrl, output.ActualSignedRequestHeaders)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Owner.DisplayName:%s, Owner.ID:%s\n", listBucketsOutput.Owner.DisplayName, listBucketsOutput.Owner.ID)
	for index, val := range listBucketsOutput.Buckets {
		fmt.Printf("Bucket[%d]-Name:%s,CreationDate:%s\n", index, val.Name, val.CreationDate)
	}
	fmt.Println()
}
func getBucketLocation() {
	output, err := OSSClient.GetBucketLocation(bucketName)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Bucket location - %s\n", output)
	fmt.Println()
}
