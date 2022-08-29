/**
 * @author      OA Wu <oawu.tw@gmail.com>
 * @copyright   Copyright (c) 2015 - 2022
 * @license     http://opensource.org/licenses/MIT  MIT License
 * @link        https://www.ioa.tw/
 */

package request

import (
	_err "errors"
	_fmt "fmt"
	_ioutil "io/ioutil"
	_http "net/http"
	_url "net/url"
	_os "os"
	_regx "regexp"
	_enum "s3/enum"
	_resp "s3/request/response"
	_sort "sort"
	_str "strings"
	_time "time"
)

type _Data struct {
	Str  string
	Size uint64
}
type _File struct {
	Path string
	Size uint64
}
type _Request struct {
	s3         _S3
	method     _enum.Method
	uri        string
	bucket     string
	data       *_Data
	file       *_File
	parameters map[string]string
	useSSL     bool
	headers    map[string]string
	amzHeaders map[string]string
}

type _S3 interface {
	GetS3Interface()
	Signature(str string) string
}

const HOST = "s3.amazonaws.com"

func rawurlencode(str string) string {
	return _str.Replace(_url.QueryEscape(str), "+", "%20", -1)
}
func New(s3 _S3) *_Request {
	if s3 == nil {
		return nil
	}
	parameters := map[string]string{}
	headers := map[string]string{
		"Host":         HOST,
		"Date":         _time.Now().In(_time.FixedZone("GMT", 0)).Format("Mon, 2 Jan 2006 15:04:05 GMT"),
		"Content-MD5":  "",
		"Content-Type": "",
	}
	amzHeaders := map[string]string{}

	req := &_Request{
		s3:         s3,
		method:     _enum.METHOD_GET,
		uri:        "",
		bucket:     "",
		data:       nil,
		file:       nil,
		useSSL:     false,
		parameters: parameters,
		headers:    headers,
		amzHeaders: amzHeaders}

	return req
}

func (req *_Request) SetHeader(key, val string) *_Request {
	if req != nil && val != "" {
		req.headers[key] = val
	}
	return req
}
func (req *_Request) SetAmzHeader(key, val string) *_Request {
	if req == nil && val != "" {
		return req
	}

	regx, err := _regx.Compile(`^x-amz-.*$`)

	if err != nil {
		return req
	}

	if len(regx.FindStringSubmatch(key)) > 0 {
		req.amzHeaders[key] = val
	} else {
		req.amzHeaders[key] = _fmt.Sprintf("x-amz-meta-%s", val)
	}

	return req
}
func (req *_Request) SetXML(str string) *_Request {
	if req != nil {
		req.data = &_Data{Str: str, Size: uint64(len(str))}
	}
	return req.SetHeader("Content-Type", "application/xml")
}
func (req *_Request) SetFile(path string, size uint64) *_Request {
	if req != nil {
		req.file = &_File{Path: path, Size: size}
	}
	return req
}
func (req *_Request) Parameter(key string, val string) *_Request {
	if req == nil || key == "" {
		return req
	}

	req.parameters[key] = val

	return req
}
func (req *_Request) Response() *_resp.Response {
	resp := _resp.New()

	if req == nil {
		return resp.GG(_err.New("錯誤的 Request"))
	}

	uri := req.uri
	resource := req.uri

	queries := []string{}
	for key, val := range req.parameters {
		queries = append(queries, _fmt.Sprintf("%s=%s", key, rawurlencode(val)))
	}

	if query := _str.Join(queries, "&"); query != "" {
		uri = _fmt.Sprintf("%s?%s", req.uri, query)
	}

	sepQueries := []string{}
	for _, key := range []string{"acl", "location", "torrent", "logging"} {
		if val, ok := req.parameters[key]; ok {
			sepQueries = append(sepQueries, _fmt.Sprintf("%s=%s", key, rawurlencode(val)))
		}
	}

	sepQuery := _str.Join(sepQueries, "&")
	switch {
	case req.bucket == "" && sepQuery == "":
		resource = _fmt.Sprintf("/%s", req.uri)
	case req.bucket != "" && sepQuery == "":
		resource = _fmt.Sprintf("/%s/%s", req.bucket, req.uri)
	case req.bucket != "" && sepQuery != "":
		resource = _fmt.Sprintf("/%s/%s?%s", req.bucket, req.uri, sepQuery)
	case req.bucket == "" && sepQuery != "":
		resource = _fmt.Sprintf("/%s?%s", req.uri, sepQuery)
	}

	protocol := "http"
	if req.useSSL {
		protocol = "https"
	}

	return req.send(_fmt.Sprintf("%s://%s/%s", protocol, req.headers["Host"], uri), req.makeHeader(resource), resp)
}
func (req *_Request) Uri(uri string) *_Request {
	if req == nil || uri == "" {
		return req
	}

	req.uri = _str.Trim(_str.Replace(rawurlencode(uri), "%2F", "/", -1), "/")
	return req
}
func (req *_Request) Bucket(bucket string) *_Request {
	if req == nil {
		return req
	}
	req.bucket = _str.ToLower(bucket)

	if req.bucket != "" {
		req.headers["Host"] = _fmt.Sprintf("%s.%s", req.bucket, HOST)
	} else {
		req.headers["Host"] = HOST
	}

	return req
}
func (req *_Request) Method(method _enum.Method) *_Request {
	if req == nil {
		return req
	}
	req.method = method
	return req
}
func (req *_Request) UseSSL(useSSL bool) *_Request {
	if req == nil {
		return req
	}
	req.useSSL = useSSL
	return req
}

func (req *_Request) makeHeader(resource string) map[string]string {
	headers := map[string]string{}
	if req == nil {
		return headers
	}

	for key, header := range req.amzHeaders {
		if header != "" {
			headers[key] = header
		}
	}
	for key, header := range req.headers {
		if header != "" {
			headers[key] = header
		}
	}

	tokens := []string{req.method.Str()}
	for _, key := range []string{"Content-MD5", "Content-Type", "Date"} {
		if val, ok := req.headers[key]; ok {
			tokens = append(tokens, val)
		}
	}

	amzs := []string{}
	for key, header := range req.amzHeaders {
		if header != "" {
			amzs = append(amzs, _fmt.Sprintf("%s:%s", _str.ToLower(key), header))
		}
	}

	if len(amzs) > 0 {
		_sort.Strings(amzs)
		for _, amz := range amzs {
			tokens = append(tokens, amz)
		}
	}

	headers["Authorization"] = req.s3.Signature(_str.Join(append(tokens, resource), "\n"))
	return headers
}

func (req *_Request) send(url string, headers map[string]string, resp *_resp.Response) *_resp.Response {
	if req == nil {
		return resp.GG(_err.New("錯誤的 Request"))
	}

	var r *_http.Request = nil
	var err error
	switch {
	case req.file != nil:
		data, e := _os.Open(req.file.Path)
		if e != nil {
			return resp.GG(e)
		}
		defer data.Close()
		r, err = _http.NewRequest(req.method.Str(), url, data)
	case req.data != nil:
		r, err = _http.NewRequest(req.method.Str(), url, _str.NewReader(req.data.Str))
	default:
		r, err = _http.NewRequest(req.method.Str(), url, nil)
	}

	if err != nil {
		return resp.GG(err)
	}

	r.TransferEncoding = []string{"identity"}

	if req.file != nil {
		r.ContentLength = int64(req.file.Size)
	}
	r.Header = _http.Header{
		"User-Agent": {"S3/Golang"},
	}
	for key, val := range headers {
		r.Header.Add(key, val)
	}

	client := &_http.Client{
		CheckRedirect: func(req *_http.Request, vias []*_http.Request) error {
			if len(vias) >= 10 {
				return _err.New("轉導次數過多")
			}
			if len(vias) == 0 {
				return nil
			}

			via := vias[0]

			req.Header = _http.Header{}
			for attr, val := range via.Header {
				req.Header[attr] = val
			}

			return nil
		},
	}

	result, err := client.Do(r)
	if err != nil {
		return resp.GG(err)
	}
	defer result.Body.Close()
	sitemap, err := _ioutil.ReadAll(result.Body)
	if err != nil {
		return resp.GG(err)
	}

	resp.StatusCode = uint16(result.StatusCode)
	for key, header := range result.Header {
		if len(header) > 0 {
			resp.Headers[key] = header[0]
		}
	}
	resp.BodyBytes = sitemap
	resp.BodyString = string(sitemap)

	if resp.StatusCode != 307 {
		return resp
	}

	location, ok := resp.Headers["Location"]
	if !ok {
		return resp
	}

	return req.send(location, headers, resp.Reset())
}
