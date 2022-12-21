# 安装

###### 更新时间： 2022-12-10

> 目录
>
> [环境准备](#环境准备)
>
> [下载sdk](#下载sdk)

## 环境准备

* 环境要求

  - 推荐使用的Golang版本：Go 1.14+。

## 下载sdk

* [sdk下载]
* [](https://github.com/dangcingzzw/inspur-go-sdk/archive/refs/heads/main.zip)





**示例代码**

## 简单文件上传

文件上传使用本地文件作为OSS文件的数据源。

以下代码用于简单文件上传：

```GO
// 引入依赖包
import (
       "fmt"
       "OSS"
)
var ak= "*** Provide your Access Key ***"
var secret= "*** Provide your Secret Key ***"
var endpoint= "https://your-endpoint"
    
// 创建OSSClient结构体 
var OSSClient, _ = OSS.New(ak, sk, endpoint)

func main() {
       input := &OSS.PutObjectInput{}
       input.Bucket = "bucketname"
       input.Key = "objectname"
       input.Body = strings.NewReader("Hello OSS")
       output, err := OSSClient.PutObject(input)
       if err == nil {
              fmt.Printf("RequestId:%s\n", output.RequestId)
              fmt.Printf("ETag:%s\n", output.ETag)
       } else if OSSError, ok := err.(OSS.OSSError); ok {
              fmt.Printf("Code:%s\n", OSSError.Code)
              fmt.Printf("Message:%s\n", OSSError.Message)
       }
}
```

## 下载对象

下载指定桶中的对象。

以下代码用于下载指定对象：  

```go
// 引入依赖包
import (
       "fmt"
       "OSS"
)
var ak= "*** Provide your Access Key ***"
var secret= "*** Provide your Secret Key ***"
var endpoint= "https://your-endpoint"
    
// 创建OSSClient结构体 
var OSSClient, _ = OSS.New(ak, sk, endpoint)

func main() {
       input := &OSS.GetObjectInput{}
       input.Bucket = "bucketname"
       input.Key = "objectname"
       // 指定开始和结束范围
       input.RangeStart = 0
       input.RangeEnd = 1000
       output, err := OSSClient.GetObject(input)
       if err == nil {
              defer output.Body.Close()
              p := make([]byte, 1024)
              var readErr error
              var readCount int
              // 获取对象内容
              for {
                     readCount, readErr = output.Body.Read(p)
                     if readCount > 0 {
                           fmt.Printf("%s", p[:readCount])
                     }
                     if readErr != nil {
                           break
                     }
              }
       } else if OSSError, ok := err.(OSS.OSSError); ok {
              fmt.Printf("Code:%s\n", OSSError.Code)
              fmt.Printf("Message:%s\n", OSSError.Message)
       }
}
```

## 删除对象

删除指定桶中的对象

以下代码用于删除指定桶中的对象：

```GO
// 引入依赖包
import (
       "fmt"
       "OSS"
)
var ak= "*** Provide your Access Key ***"
var secret= "*** Provide your Secret Key ***"
var endpoint= "https://your-endpoint"
    
// 创建OSSClient结构体 
var OSSClient, _ = OSS.New(ak, sk, endpoint)

func main() {
       input := &OSS.DeleteObjectInput{}
       input.Bucket = "bucketname"
       input.Key = "objectname"
       output, err := OSSClient.DeleteObject(input)
       if err == nil {
              fmt.Printf("RequestId:%s\n", output.RequestId)
       } else if OSSError, ok := err.(OSS.OSSError); ok {
              fmt.Printf("Code:%s\n", OSSError.Code)
              fmt.Printf("Message:%s\n", OSSError.Message)
       }
}
```

## 简单创建存储桶

以下代码用于简单创建存储桶：

```php
// 引入依赖包
import (
       "fmt"
       "OSS"
)
var ak= "*** Provide your Access Key ***"
var secret= "*** Provide your Secret Key ***"
var endpoint= "https://your-endpoint"
    
// 创建OSSClient结构体 
var OSSClient, _ = OSS.New(ak, sk, endpoint)
    
func main() {
       input := &OSS.CreateBucketInput{}
       input.Bucket = "bucketname"
       input.Location = "bucketlocation"
       input.ACL = OSS.AclPrivate
       input.StorageClass = OSS.StorageClassWarm
       input.AvailableZone = "3az"
       output, err := OSSClient.CreateBucket(input)
       if err == nil {
              fmt.Printf("RequestId:%s\n", output.RequestId)
       } else {
              if OSSError, ok := err.(OSS.OSSError); ok {
                     fmt.Println(OSSError.Code)
                     fmt.Println(OSSError.Message)
              } else {
                     fmt.Println(err)
              }
       }
}
```

存储桶（Bucket）是存储对象（Object）的容器，对象都隶属于存储桶。

本节介绍如何删除存储桶。

<font color="red">⚠   删除存储桶之前，必须先删除存储桶下的所有文件、分片上传产生的碎片。</font>

以下代码用于删除存储桶：

```php
// 引入依赖包
import (
       "fmt"
       "OSS"
)
var ak= "*** Provide your Access Key ***"
var secret= "*** Provide your Secret Key ***"
var endpoint= "https://your-endpoint"
    
// 创建OSSClient结构体 
var OSSClient, _ = OSS.New(ak, sk, endpoint)

func main() {
       output, err := OSSClient.DeleteBucket("bucketname")
       if err == nil {
              fmt.Printf("RequestId:%s\n", output.RequestId)
       } else {
              if OSSError, ok := err.(OSS.OSSError); ok {
                     fmt.Println(OSSError.Code)
                     fmt.Println(OSSError.Message)
              } else {
                     fmt.Println(err)
              }
       }
}

```
