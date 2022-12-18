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
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func prepareHeaders(headers map[string][]string, meta bool, isOSS bool) map[string][]string {
	_headers := make(map[string][]string, len(headers))
	if headers != nil {
		for key, value := range headers {
			key = strings.TrimSpace(key)
			if key == "" {
				continue
			}
			_key := strings.ToLower(key)
			if _, ok := allowedRequestHTTPHeaderMetadataNames[_key]; !ok && !strings.HasPrefix(key, HEADER_PREFIX) && !strings.HasPrefix(key, HEADER_PREFIX_OSS) {
				if !meta {
					continue
				}
				if !isOSS {
					_key = HEADER_PREFIX_META + _key
				} else {
					_key = HEADER_PREFIX_META_OSS + _key
				}
			} else {
				_key = key
			}
			_headers[_key] = value
		}
	}
	return _headers
}

func (OSSClient OSSClient) doActionWithoutBucket(action, method string, input ISerializable, output IBaseModel, extensions []extensionOptions) error {
	return OSSClient.doAction(action, method, "", "", input, output, true, true, extensions)
}

func (OSSClient OSSClient) doActionWithBucketV2(action, method, bucketName string, input ISerializable, output IBaseModel, extensions []extensionOptions) error {
	if strings.TrimSpace(bucketName) == "" && !OSSClient.conf.cname {
		return errors.New("Bucket is empty")
	}
	return OSSClient.doAction(action, method, bucketName, "", input, output, false, true, extensions)
}

func (OSSClient OSSClient) doActionWithBucket(action, method, bucketName string, input ISerializable, output IBaseModel, extensions []extensionOptions) error {
	if strings.TrimSpace(bucketName) == "" && !OSSClient.conf.cname {
		return errors.New("Bucket is empty")
	}
	return OSSClient.doAction(action, method, bucketName, "", input, output, true, true, extensions)
}

func (OSSClient OSSClient) doActionWithBucketAndKey(action, method, bucketName, objectKey string, input ISerializable, output IBaseModel, extensions []extensionOptions) error {
	return OSSClient._doActionWithBucketAndKey(action, method, bucketName, objectKey, input, output, true, extensions)
}

func (OSSClient OSSClient) doActionWithBucketAndKeyV2(action, method, bucketName, objectKey string, input ISerializable, output IBaseModel, extensions []extensionOptions) error {
	if strings.TrimSpace(bucketName) == "" && !OSSClient.conf.cname {
		return errors.New("Bucket is empty")
	}
	if strings.TrimSpace(objectKey) == "" {
		return errors.New("Key is empty")
	}
	return OSSClient.doAction(action, method, bucketName, objectKey, input, output, false, true, extensions)
}

func (OSSClient OSSClient) doActionWithBucketAndKeyUnRepeatable(action, method, bucketName, objectKey string, input ISerializable, output IBaseModel, extensions []extensionOptions) error {
	return OSSClient._doActionWithBucketAndKey(action, method, bucketName, objectKey, input, output, false, extensions)
}

func (OSSClient OSSClient) _doActionWithBucketAndKey(action, method, bucketName, objectKey string, input ISerializable, output IBaseModel, repeatable bool, extensions []extensionOptions) error {
	if strings.TrimSpace(bucketName) == "" && !OSSClient.conf.cname {
		return errors.New("Bucket is empty")
	}
	if strings.TrimSpace(objectKey) == "" {
		return errors.New("Key is empty")
	}
	return OSSClient.doAction(action, method, bucketName, objectKey, input, output, true, repeatable, extensions)
}

func (OSSClient OSSClient) doAction(action, method, bucketName, objectKey string, input ISerializable, output IBaseModel, xmlResult bool, repeatable bool, extensions []extensionOptions) error {

	var resp *http.Response
	var respError error
	doLog(LEVEL_INFO, "Enter method %s...", action)
	start := GetCurrentTimestamp()

	params, headers, data, err := input.trans(OSSClient.conf.signature == SignatureOSS)
	if err != nil {
		return err
	}

	if params == nil {
		params = make(map[string]string)
	}

	if headers == nil {
		headers = make(map[string][]string)
	}

	for _, extension := range extensions {
		if extensionHeader, ok := extension.(extensionHeaders); ok {
			_err := extensionHeader(headers, OSSClient.conf.signature == SignatureOSS)
			if _err != nil {
				doLog(LEVEL_INFO, fmt.Sprintf("set header with error: %v", _err))
			}
		} else {
			doLog(LEVEL_INFO, "Unsupported extensionOptions")
		}
	}

	switch method {
	case HTTP_GET:
		resp, respError = OSSClient.doHTTPGet(bucketName, objectKey, params, headers, data, repeatable)
	case HTTP_POST:
		resp, respError = OSSClient.doHTTPPost(bucketName, objectKey, params, headers, data, repeatable)
	case HTTP_PUT:
		resp, respError = OSSClient.doHTTPPut(bucketName, objectKey, params, headers, data, repeatable)
	case HTTP_DELETE:
		resp, respError = OSSClient.doHTTPDelete(bucketName, objectKey, params, headers, data, repeatable)
	case HTTP_HEAD:
		resp, respError = OSSClient.doHTTPHead(bucketName, objectKey, params, headers, data, repeatable)
	case HTTP_OPTIONS:
		resp, respError = OSSClient.doHTTPOptions(bucketName, objectKey, params, headers, data, repeatable)
	default:
		respError = errors.New("Unexpect http method error")
	}
	if respError == nil && output != nil {
		respError = ParseResponseToBaseModel(resp, output, xmlResult, OSSClient.conf.signature == SignatureOSS)
		if respError != nil {
			doLog(LEVEL_WARN, "Parse response to BaseModel with error: %v", respError)
		}
	} else {
		doLog(LEVEL_WARN, "Do http request with error: %v", respError)
	}

	if isDebugLogEnabled() {
		doLog(LEVEL_DEBUG, "End method %s, OSSclient cost %d ms", action, (GetCurrentTimestamp() - start))
	}

	return respError
}

func (OSSClient OSSClient) doHTTPGet(bucketName, objectKey string, params map[string]string,
	headers map[string][]string, data interface{}, repeatable bool) (*http.Response, error) {
	return OSSClient.doHTTP(HTTP_GET, bucketName, objectKey, params, prepareHeaders(headers, false, OSSClient.conf.signature == SignatureOSS), data, repeatable)
}

func (OSSClient OSSClient) doHTTPHead(bucketName, objectKey string, params map[string]string,
	headers map[string][]string, data interface{}, repeatable bool) (*http.Response, error) {
	return OSSClient.doHTTP(HTTP_HEAD, bucketName, objectKey, params, prepareHeaders(headers, false, OSSClient.conf.signature == SignatureOSS), data, repeatable)
}

func (OSSClient OSSClient) doHTTPOptions(bucketName, objectKey string, params map[string]string,
	headers map[string][]string, data interface{}, repeatable bool) (*http.Response, error) {
	return OSSClient.doHTTP(HTTP_OPTIONS, bucketName, objectKey, params, prepareHeaders(headers, false, OSSClient.conf.signature == SignatureOSS), data, repeatable)
}

func (OSSClient OSSClient) doHTTPDelete(bucketName, objectKey string, params map[string]string,
	headers map[string][]string, data interface{}, repeatable bool) (*http.Response, error) {
	return OSSClient.doHTTP(HTTP_DELETE, bucketName, objectKey, params, prepareHeaders(headers, false, OSSClient.conf.signature == SignatureOSS), data, repeatable)
}

func (OSSClient OSSClient) doHTTPPut(bucketName, objectKey string, params map[string]string,
	headers map[string][]string, data interface{}, repeatable bool) (*http.Response, error) {
	return OSSClient.doHTTP(HTTP_PUT, bucketName, objectKey, params, prepareHeaders(headers, true, OSSClient.conf.signature == SignatureOSS), data, repeatable)
}

func (OSSClient OSSClient) doHTTPPost(bucketName, objectKey string, params map[string]string,
	headers map[string][]string, data interface{}, repeatable bool) (*http.Response, error) {
	return OSSClient.doHTTP(HTTP_POST, bucketName, objectKey, params, prepareHeaders(headers, true, OSSClient.conf.signature == SignatureOSS), data, repeatable)
}

func prepareAgentHeader(clientUserAgent string) string {
	userAgent := USER_AGENT
	if clientUserAgent != "" {
		userAgent = clientUserAgent
	}
	return userAgent
}

func (OSSClient OSSClient) getSignedURLResponse(action string, output IBaseModel, xmlResult bool, resp *http.Response, err error, start int64) (respError error) {
	var msg interface{}
	if err != nil {
		respError = err
		resp = nil
	} else {
		doLog(LEVEL_DEBUG, "Response headers: %v", resp.Header)
		if resp.StatusCode >= 300 {
			respError = ParseResponseToOSSError(resp, OSSClient.conf.signature == SignatureOSS)
			msg = resp.Status
			resp = nil
		} else {
			if output != nil {
				respError = ParseResponseToBaseModel(resp, output, xmlResult, OSSClient.conf.signature == SignatureOSS)
			}
			if respError != nil {
				doLog(LEVEL_WARN, "Parse response to BaseModel with error: %v", respError)
			}
		}
	}

	if msg != nil {
		doLog(LEVEL_ERROR, "Failed to send request with reason:%v", msg)
	}

	if isDebugLogEnabled() {
		doLog(LEVEL_DEBUG, "End method %s, OSSclient cost %d ms", action, (GetCurrentTimestamp() - start))
	}
	return
}

func (OSSClient OSSClient) doHTTPWithSignedURL(action, method string, signedURL string, actualSignedRequestHeaders http.Header, data io.Reader, output IBaseModel, xmlResult bool) (respError error) {
	req, err := http.NewRequest(method, signedURL, data)
	if err != nil {
		return err
	}
	if OSSClient.conf.ctx != nil {
		req = req.WithContext(OSSClient.conf.ctx)
	}
	var resp *http.Response

	var isSecurityToken bool
	var securityToken string
	var query []string
	parmas := strings.Split(signedURL, "?")
	if len(parmas) > 1 {
		query = strings.Split(parmas[1], "&")
		for _, value := range query {
			if strings.HasPrefix(value, HEADER_STS_TOKEN_AMZ+"=") || strings.HasPrefix(value, HEADER_STS_TOKEN_OSS+"=") {
				if value[len(HEADER_STS_TOKEN_AMZ)+1:] != "" {
					securityToken = value[len(HEADER_STS_TOKEN_AMZ)+1:]
					isSecurityToken = true
				}
			}
		}
	}
	logSignedURL := signedURL
	if isSecurityToken {
		logSignedURL = strings.Replace(logSignedURL, securityToken, "******", -1)
	}
	doLog(LEVEL_INFO, "Do %s with signedUrl %s...", action, logSignedURL)

	req.Header = actualSignedRequestHeaders
	if value, ok := req.Header[HEADER_HOST_CAMEL]; ok {
		req.Host = value[0]
		delete(req.Header, HEADER_HOST_CAMEL)
	} else if value, ok := req.Header[HEADER_HOST]; ok {
		req.Host = value[0]
		delete(req.Header, HEADER_HOST)
	}

	if value, ok := req.Header[HEADER_CONTENT_LENGTH_CAMEL]; ok {
		req.ContentLength = StringToInt64(value[0], -1)
		delete(req.Header, HEADER_CONTENT_LENGTH_CAMEL)
	} else if value, ok := req.Header[HEADER_CONTENT_LENGTH]; ok {
		req.ContentLength = StringToInt64(value[0], -1)
		delete(req.Header, HEADER_CONTENT_LENGTH)
	}

	userAgent := prepareAgentHeader(OSSClient.conf.userAgent)
	req.Header[HEADER_USER_AGENT_CAMEL] = []string{userAgent}
	start := GetCurrentTimestamp()
	resp, err = OSSClient.httpClient.Do(req)
	if isInfoLogEnabled() {
		doLog(LEVEL_INFO, "Do http request cost %d ms", (GetCurrentTimestamp() - start))
	}

	respError = OSSClient.getSignedURLResponse(action, output, xmlResult, resp, err, start)

	return
}

func prepareData(headers map[string][]string, data interface{}) (io.Reader, error) {
	var _data io.Reader
	if data != nil {
		if dataStr, ok := data.(string); ok {
			doLog(LEVEL_DEBUG, "Do http request with string: %s", dataStr)
			headers["Content-Length"] = []string{IntToString(len(dataStr))}
			_data = strings.NewReader(dataStr)
		} else if dataByte, ok := data.([]byte); ok {
			doLog(LEVEL_DEBUG, "Do http request with byte array")
			headers["Content-Length"] = []string{IntToString(len(dataByte))}
			_data = bytes.NewReader(dataByte)
		} else if dataReader, ok := data.(io.Reader); ok {
			_data = dataReader
		} else {
			doLog(LEVEL_WARN, "Data is not a valid io.Reader")
			return nil, errors.New("Data is not a valid io.Reader")
		}
	}
	return _data, nil
}

func (OSSClient OSSClient) getRequest(redirectURL, requestURL string, redirectFlag bool, _data io.Reader, method,
	bucketName, objectKey string, params map[string]string, headers map[string][]string) (*http.Request, error) {
	if redirectURL != "" {
		if !redirectFlag {
			parsedRedirectURL, err := url.Parse(redirectURL)
			if err != nil {
				return nil, err
			}
			requestURL, err = OSSClient.doAuth(method, bucketName, objectKey, params, headers, parsedRedirectURL.Host)
			if err != nil {
				return nil, err
			}
			if parsedRequestURL, err := url.Parse(requestURL); err != nil {
				return nil, err
			} else if parsedRequestURL.RawQuery != "" && parsedRedirectURL.RawQuery == "" {
				redirectURL += "?" + parsedRequestURL.RawQuery
			}
		}
		requestURL = redirectURL
	} else {
		var err error
		requestURL, err = OSSClient.doAuth(method, bucketName, objectKey, params, headers, "")
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, requestURL, _data)
	if OSSClient.conf.ctx != nil {
		req = req.WithContext(OSSClient.conf.ctx)
	}
	if err != nil {
		return nil, err
	}
	doLog(LEVEL_DEBUG, "Do request with url [%s] and method [%s]", requestURL, method)
	return req, nil
}

func logHeaders(headers map[string][]string, signature SignatureType) {
	if isDebugLogEnabled() {
		auth := headers[HEADER_AUTH_CAMEL]
		delete(headers, HEADER_AUTH_CAMEL)

		var isSecurityToken bool
		var securityToken []string
		if securityToken, isSecurityToken = headers[HEADER_STS_TOKEN_AMZ]; isSecurityToken {
			headers[HEADER_STS_TOKEN_AMZ] = []string{"******"}
		} else if securityToken, isSecurityToken = headers[HEADER_STS_TOKEN_OSS]; isSecurityToken {
			headers[HEADER_STS_TOKEN_OSS] = []string{"******"}
		}
		doLog(LEVEL_DEBUG, "Request headers: %v", headers)
		headers[HEADER_AUTH_CAMEL] = auth
		if isSecurityToken {
			if signature == SignatureOSS {
				headers[HEADER_STS_TOKEN_OSS] = securityToken
			} else {
				headers[HEADER_STS_TOKEN_AMZ] = securityToken
			}
		}
	}
}

func prepareReq(headers map[string][]string, req, lastRequest *http.Request, clientUserAgent string) *http.Request {
	for key, value := range headers {
		if key == HEADER_HOST_CAMEL {
			req.Host = value[0]
			delete(headers, key)
		} else if key == HEADER_CONTENT_LENGTH_CAMEL {
			req.ContentLength = StringToInt64(value[0], -1)
			delete(headers, key)
		} else {
			req.Header[key] = value
		}
	}

	lastRequest = req

	userAgent := prepareAgentHeader(clientUserAgent)
	req.Header[HEADER_USER_AGENT_CAMEL] = []string{userAgent}

	if lastRequest != nil {
		req.Host = lastRequest.Host
		req.ContentLength = lastRequest.ContentLength
	}
	return lastRequest
}

func canNotRetry(repeatable bool, statusCode int) bool {
	if !repeatable || (statusCode >= 400 && statusCode < 500) || statusCode == 304 {
		return true
	}
	return false
}

func isRedirectErr(location string, redirectCount, maxRedirectCount int) bool {
	if location != "" && redirectCount < maxRedirectCount {
		return true
	}
	return false
}

func setRedirectFlag(statusCode int, method string) (redirectFlag bool) {
	if statusCode == 302 && method == HTTP_GET {
		redirectFlag = true
	} else {
		redirectFlag = false
	}
	return
}

func prepareRetry(resp *http.Response, headers map[string][]string, _data io.Reader, msg interface{}) (io.Reader, *http.Response, error) {
	if resp != nil {
		_err := resp.Body.Close()
		checkAndLogErr(_err, LEVEL_WARN, "Failed to close resp body")
		resp = nil
	}
	if _, ok := headers[HEADER_AUTH_CAMEL]; ok {
		delete(headers, HEADER_AUTH_CAMEL)
	}
	doLog(LEVEL_WARN, "Failed to send request with reason:%v, will try again", msg)
	if r, ok := _data.(*strings.Reader); ok {
		_, err := r.Seek(0, 0)
		if err != nil {
			return nil, nil, err
		}
	} else if r, ok := _data.(*bytes.Reader); ok {
		_, err := r.Seek(0, 0)
		if err != nil {
			return nil, nil, err
		}
	} else if r, ok := _data.(*fileReaderWrapper); ok {
		fd, err := os.Open(r.filePath)
		if err != nil {
			return nil, nil, err
		}
		fileReaderWrapper := &fileReaderWrapper{filePath: r.filePath}
		fileReaderWrapper.mark = r.mark
		fileReaderWrapper.reader = fd
		fileReaderWrapper.totalCount = r.totalCount
		_data = fileReaderWrapper
		_, err = fd.Seek(r.mark, 0)
		if err != nil {
			errMsg := fd.Close()
			checkAndLogErr(errMsg, LEVEL_WARN, "Failed to close with reason: %v", errMsg)
			return nil, nil, err
		}
	} else if r, ok := _data.(*readerWrapper); ok {
		_, err := r.seek(0, 0)
		if err != nil {
			return nil, nil, err
		}
		r.readedCount = 0
	}
	return _data, resp, nil
}

func (OSSClient OSSClient) doHTTP(method, bucketName, objectKey string, params map[string]string,
	headers map[string][]string, data interface{}, repeatable bool) (resp *http.Response, respError error) {

	bucketName = strings.TrimSpace(bucketName)

	method = strings.ToUpper(method)

	var redirectURL string
	var requestURL string
	maxRetryCount := OSSClient.conf.maxRetryCount
	maxRedirectCount := OSSClient.conf.maxRedirectCount

	_data, _err := prepareData(headers, data)
	if _err != nil {
		return nil, _err
	}

	var lastRequest *http.Request
	redirectFlag := false
	for i, redirectCount := 0, 0; i <= maxRetryCount; i++ {
		req, err := OSSClient.getRequest(redirectURL, requestURL, redirectFlag, _data,
			method, bucketName, objectKey, params, headers)
		if err != nil {
			return nil, err
		}

		logHeaders(headers, OSSClient.conf.signature)

		lastRequest = prepareReq(headers, req, lastRequest, OSSClient.conf.userAgent)

		start := GetCurrentTimestamp()
		resp, err = OSSClient.httpClient.Do(req)
		doLog(LEVEL_INFO, "Do http request cost %d ms", (GetCurrentTimestamp() - start))

		var msg interface{}
		if err != nil {
			msg = err
			respError = err
			resp = nil
			if !repeatable {
				break
			}
		} else {
			doLog(LEVEL_DEBUG, "Response headers: %v", resp.Header)
			if resp.StatusCode < 300 {
				respError = nil
				break
			} else if canNotRetry(repeatable, resp.StatusCode) {
				respError = ParseResponseToOSSError(resp, OSSClient.conf.signature == SignatureOSS)
				resp = nil
				break
			} else if resp.StatusCode >= 300 && resp.StatusCode < 400 {
				location := resp.Header.Get(HEADER_LOCATION_CAMEL)
				if isRedirectErr(location, redirectCount, maxRedirectCount) {
					redirectURL = location
					doLog(LEVEL_WARN, "Redirect request to %s", redirectURL)
					msg = resp.Status
					maxRetryCount++
					redirectCount++
					redirectFlag = setRedirectFlag(resp.StatusCode, method)
				} else {
					respError = ParseResponseToOSSError(resp, OSSClient.conf.signature == SignatureOSS)
					resp = nil
					break
				}
			} else {
				msg = resp.Status
			}
		}
		if i != maxRetryCount {
			_data, resp, err = prepareRetry(resp, headers, _data, msg)
			if err != nil {
				return nil, err
			}
			if r, ok := _data.(*fileReaderWrapper); ok {
				if _fd, _ok := r.reader.(*os.File); _ok {
					defer func() {
						errMsg := _fd.Close()
						checkAndLogErr(errMsg, LEVEL_WARN, "Failed to close with reason: %v", errMsg)
					}()
				}
			}
			time.Sleep(time.Duration(float64(i+2) * rand.Float64() * float64(time.Second)))
		} else {
			doLog(LEVEL_ERROR, "Failed to send request with reason:%v", msg)
			if resp != nil {
				respError = ParseResponseToOSSError(resp, OSSClient.conf.signature == SignatureOSS)
				resp = nil
			}
		}
	}
	return
}

type connDelegate struct {
	conn          net.Conn
	socketTimeout time.Duration
	finalTimeout  time.Duration
}

func getConnDelegate(conn net.Conn, socketTimeout int, finalTimeout int) *connDelegate {
	return &connDelegate{
		conn:          conn,
		socketTimeout: time.Second * time.Duration(socketTimeout),
		finalTimeout:  time.Second * time.Duration(finalTimeout),
	}
}

func (delegate *connDelegate) Read(b []byte) (n int, err error) {
	setReadDeadlineErr := delegate.SetReadDeadline(time.Now().Add(delegate.socketTimeout))
	flag := isDebugLogEnabled()

	if setReadDeadlineErr != nil && flag {
		doLog(LEVEL_DEBUG, "Failed to set read deadline with reason: %v, but it's ok", setReadDeadlineErr)
	}

	n, err = delegate.conn.Read(b)
	setReadDeadlineErr = delegate.SetReadDeadline(time.Now().Add(delegate.finalTimeout))
	if setReadDeadlineErr != nil && flag {
		doLog(LEVEL_DEBUG, "Failed to set read deadline with reason: %v, but it's ok", setReadDeadlineErr)
	}
	return n, err
}

func (delegate *connDelegate) Write(b []byte) (n int, err error) {
	setWriteDeadlineErr := delegate.SetWriteDeadline(time.Now().Add(delegate.socketTimeout))
	flag := isDebugLogEnabled()
	if setWriteDeadlineErr != nil && flag {
		doLog(LEVEL_DEBUG, "Failed to set write deadline with reason: %v, but it's ok", setWriteDeadlineErr)
	}

	n, err = delegate.conn.Write(b)
	finalTimeout := time.Now().Add(delegate.finalTimeout)
	setWriteDeadlineErr = delegate.SetWriteDeadline(finalTimeout)
	if setWriteDeadlineErr != nil && flag {
		doLog(LEVEL_DEBUG, "Failed to set write deadline with reason: %v, but it's ok", setWriteDeadlineErr)
	}
	setReadDeadlineErr := delegate.SetReadDeadline(finalTimeout)
	if setReadDeadlineErr != nil && flag {
		doLog(LEVEL_DEBUG, "Failed to set read deadline with reason: %v, but it's ok", setReadDeadlineErr)
	}
	return n, err
}

func (delegate *connDelegate) Close() error {
	return delegate.conn.Close()
}

func (delegate *connDelegate) LocalAddr() net.Addr {
	return delegate.conn.LocalAddr()
}

func (delegate *connDelegate) RemoteAddr() net.Addr {
	return delegate.conn.RemoteAddr()
}

func (delegate *connDelegate) SetDeadline(t time.Time) error {
	return delegate.conn.SetDeadline(t)
}

func (delegate *connDelegate) SetReadDeadline(t time.Time) error {
	return delegate.conn.SetReadDeadline(t)
}

func (delegate *connDelegate) SetWriteDeadline(t time.Time) error {
	return delegate.conn.SetWriteDeadline(t)
}
