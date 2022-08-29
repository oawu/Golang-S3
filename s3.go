/**
 * @author      OA Wu <oawu.tw@gmail.com>
 * @copyright   Copyright (c) 2015 - 2022
 * @license     http://opensource.org/licenses/MIT  MIT License
 * @link        https://www.ioa.tw/
 */

package s3

import (
	_hmac "crypto/hmac"
	_md5 "crypto/md5"
	_sha1 "crypto/sha1"
	_base64 "encoding/base64"
	_hex "encoding/hex"
	_xml "encoding/xml"
	_err "errors"
	_fmt "fmt"
	_bucket "s3/bucket"
	_enum "s3/enum"
	_model "s3/model"
	_req "s3/request"
	_time "time"
)

type S3 struct {
	_access string
	_secret string
}

var (
	_instances = map[string]*S3{}
	_exts      = map[string]string{".jpg": "image/jpeg", ".gif": "image/gif", ".png": "image/png", ".pdf": "application/pdf", ".gz": "application/x-gzip", ".zip": "application/x-zip", ".swf": "application/x-shockwave-flash", ".tar": "application/x-tar", ".bz": "application/x-bzip", ".bz2": "application/x-bzip2", ".txt": "text/plain", ".html": "text/html", ".htm": "text/html", ".ico": "image/x-icon", ".css": "text/css", ".js": "application/x-javascript", ".xml": "text/xml", ".ogg": "application/ogg", ".wav": "audio/x-wav", ".avi": "video/x-msvideo", ".mpg": "video/mpeg", ".mov": "video/quicktime", ".mp3": "audio/mpeg", ".mpeg": "video/mpeg", ".flv": "video/x-flv", ".php": "application/x-httpd-php", ".bin": "application/macbinary", ".psd": "application/x-photoshop", ".ai": "application/postscript", ".ppt": "application/powerpoint", ".wbxml": "application/wbxml", ".tgz": "application/x-tar", ".jpeg": "image/jpeg", ".jpe": "image/jpeg", ".bmp": "image/bmp", ".shtml": "text/html", ".text": "text/plain", ".doc": "application/msword", ".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document", ".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", ".word": "application/msword", ".json": "application/json", ".svg": "image/svg+xml", ".mp2": "audio/mpeg", ".exe": "application/octet-stream", ".tif": "image/tiff", ".tiff": "image/tiff", ".asc": "text/plain", ".xsl": "text/xml", ".hqx": "application/mac-binhex40", ".cpt": "application/mac-compactpro", ".csv": "text/x-comma-separated-values", ".dms": "application/octet-stream", ".lha": "application/octet-stream", ".lzh": "application/octet-stream", ".class": "application/octet-stream", ".so": "application/octet-stream", ".sea": "application/octet-stream", ".dll": "application/octet-stream", ".oda": "application/oda", ".eps": "application/postscript", ".ps": "application/postscript", ".smi": "application/smil", ".smil": "application/smil", ".mif": "application/vnd.mif", ".xls": "application/excel", ".wmlc": "application/wmlc", ".dcr": "application/x-director", ".dir": "application/x-director", ".dxr": "application/x-director", ".dvi": "application/x-dvi", ".gtar": "application/x-gtar", ".php4": "application/x-httpd-php", ".php3": "application/x-httpd-php", ".phtml": "application/x-httpd-php", ".phps": "application/x-httpd-php-source", ".sit": "application/x-stuffit", ".xhtml": "application/xhtml+xml", ".xht": "application/xhtml+xml", ".mid": "audio/midi", ".midi": "audio/midi", ".mpga": "audio/mpeg", ".aif": "audio/x-aiff", ".aiff": "audio/x-aiff", ".aifc": "audio/x-aiff", ".ram": "audio/x-pn-realaudio", ".rm": "audio/x-pn-realaudio", ".rpm": "audio/x-pn-realaudio-plugin", ".ra": "audio/x-realaudio", ".rv": "video/vnd.rn-realvideo", ".log": "text/plain", ".rtx": "text/richtext", ".rtf": "text/rtf", ".mpe": "video/mpeg", ".qt": "video/quicktime", ".movie": "video/x-sgi-movie", ".xl": "application/excel", ".eml": "message/rfc822"}
)

func hashHmacSha1(password string, salt string) []byte {
	h := _hmac.New(_sha1.New, []byte(salt))
	h.Write([]byte(password))
	return h.Sum(nil)
}

func Instance(access string, secret string) *S3 {
	key := _fmt.Sprintf("%s_%s", access, secret)
	hash := _md5.Sum([]byte(key))
	key = _hex.EncodeToString(hash[:])

	if instance, ok := _instances[key]; ok {
		return instance
	}

	_instances[key] = (&S3{_access: access, _secret: secret})

	return _instances[key]
}

func (s3 *S3) GetS3Interface() {}

func (s3 *S3) Signature(str string) string {
	if s3 == nil {
		return ""
	}

	return _fmt.Sprintf("AWS %s:%s", s3._access, _base64.StdEncoding.EncodeToString(hashHmacSha1(str, s3._secret)))
}
func (s3 *S3) Test() bool {
	response := _req.New(s3).Method(_enum.METHOD_GET).UseSSL(true).Response()
	return response.Error == nil && response.StatusCode == 200
}
func (s3 *S3) Info() (*_model.BucketInfo, error) {
	response := _req.New(s3).Method(_enum.METHOD_GET).UseSSL(true).Response()

	if err := response.IsSuccess(); err != nil {
		return nil, err
	}

	if response.Headers["Content-Type"] != "application/xml" {
		return nil, _err.New("錯誤，回應結果非 XML 格式")
	}

	var result *struct {
		Owner struct {
			Id          string `xml:"ID"`
			DisplayName string `xml:"DisplayName"`
		} `xml:"Owner"`

		Buckets []struct {
			Name string `xml:"Name"`
			Date string `xml:"CreationDate"`
		} `xml:"Buckets>Bucket"`
	} = nil

	if err := _xml.Unmarshal(response.BodyBytes, &result); err != nil {
		return nil, _err.New(_fmt.Sprintf("編譯 XML 失敗，Message：%s", err))
	}

	buckets := []_model.BucketInfoBucket{}
	for _, bucket := range result.Buckets {
		time, err := _time.Parse("2006-01-02T15:04:05.999Z", bucket.Date)
		if err != nil {
			return nil, _err.New(_fmt.Sprintf("轉換時間格式失敗，Message：%s", err))
		}

		buckets = append(buckets, _model.BucketInfoBucket{
			Name: bucket.Name,
			Time: uint64(time.Unix()),
		})
	}

	return &_model.BucketInfo{Owner: &_model.BucketInfoOwner{Id: result.Owner.Id, Name: result.Owner.DisplayName}, Buckets: buckets}, nil
}
func (s3 *S3) Buckets() ([]string, error) {
	buckets := []string{}

	info, err := s3.Info()
	if err != nil {
		return buckets, _err.New(_fmt.Sprintf("取得 Info 資訊失敗，Message：%s", err))
	}

	if info == nil {
		return buckets, _err.New("取得 Info 資訊失敗")
	}

	for _, bucket := range info.Buckets {
		buckets = append(buckets, bucket.Name)
	}

	return buckets, nil
}
func (s3 *S3) Bucket(name string) *_bucket.Bucket {
	bucket, _ := _bucket.New(name, s3)
	return bucket
}
