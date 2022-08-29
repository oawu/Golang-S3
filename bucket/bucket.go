/**
 * @author      OA Wu <oawu.tw@gmail.com>
 * @copyright   Copyright (c) 2015 - 2022
 * @license     http://opensource.org/licenses/MIT  MIT License
 * @link        https://www.ioa.tw/
 */

package bucket

import (
	_md5 "crypto/md5"
	_base64 "encoding/base64"
	_xml "encoding/xml"
	_err "errors"
	_fmt "fmt"
	_io "io"
	_http "net/http"
	_os "os"
	_fs "path/filepath"
	_enum "s3/enum"
	_model "s3/model"
	_req "s3/request"
	_resp "s3/request/response"
	_strconv "strconv"
	_str "strings"
	_sync "sync"
	_time "time"
)

var (
	_exts = map[string]string{".jpg": "image/jpeg", ".gif": "image/gif", ".png": "image/png", ".pdf": "application/pdf", ".gz": "application/x-gzip", ".zip": "application/x-zip", ".swf": "application/x-shockwave-flash", ".tar": "application/x-tar", ".bz": "application/x-bzip", ".bz2": "application/x-bzip2", ".txt": "text/plain", ".html": "text/html", ".htm": "text/html", ".ico": "image/x-icon", ".css": "text/css", ".js": "application/x-javascript", ".xml": "text/xml", ".ogg": "application/ogg", ".wav": "audio/x-wav", ".avi": "video/x-msvideo", ".mpg": "video/mpeg", ".mov": "video/quicktime", ".mp3": "audio/mpeg", ".mpeg": "video/mpeg", ".flv": "video/x-flv", ".php": "application/x-httpd-php", ".bin": "application/macbinary", ".psd": "application/x-photoshop", ".ai": "application/postscript", ".ppt": "application/powerpoint", ".wbxml": "application/wbxml", ".tgz": "application/x-tar", ".jpeg": "image/jpeg", ".jpe": "image/jpeg", ".bmp": "image/bmp", ".shtml": "text/html", ".text": "text/plain", ".doc": "application/msword", ".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document", ".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", ".word": "application/msword", ".json": "application/json", ".svg": "image/svg+xml", ".mp2": "audio/mpeg", ".exe": "application/octet-stream", ".tif": "image/tiff", ".tiff": "image/tiff", ".asc": "text/plain", ".xsl": "text/xml", ".hqx": "application/mac-binhex40", ".cpt": "application/mac-compactpro", ".csv": "text/x-comma-separated-values", ".dms": "application/octet-stream", ".lha": "application/octet-stream", ".lzh": "application/octet-stream", ".class": "application/octet-stream", ".so": "application/octet-stream", ".sea": "application/octet-stream", ".dll": "application/octet-stream", ".oda": "application/oda", ".eps": "application/postscript", ".ps": "application/postscript", ".smi": "application/smil", ".smil": "application/smil", ".mif": "application/vnd.mif", ".xls": "application/excel", ".wmlc": "application/wmlc", ".dcr": "application/x-director", ".dir": "application/x-director", ".dxr": "application/x-director", ".dvi": "application/x-dvi", ".gtar": "application/x-gtar", ".php4": "application/x-httpd-php", ".php3": "application/x-httpd-php", ".phtml": "application/x-httpd-php", ".phps": "application/x-httpd-php-source", ".sit": "application/x-stuffit", ".xhtml": "application/xhtml+xml", ".xht": "application/xhtml+xml", ".mid": "audio/midi", ".midi": "audio/midi", ".mpga": "audio/mpeg", ".aif": "audio/x-aiff", ".aiff": "audio/x-aiff", ".aifc": "audio/x-aiff", ".ram": "audio/x-pn-realaudio", ".rm": "audio/x-pn-realaudio", ".rpm": "audio/x-pn-realaudio-plugin", ".ra": "audio/x-realaudio", ".rv": "video/vnd.rn-realvideo", ".log": "text/plain", ".rtx": "text/richtext", ".rtf": "text/rtf", ".mpe": "video/mpeg", ".qt": "video/quicktime", ".movie": "video/x-sgi-movie", ".xl": "application/excel", ".eml": "message/rfc822"}
)

type Bucket struct {
	s3   _S3
	name string
	uri  string
}

type _S3 interface {
	GetS3Interface()
	Bucket(name string) *Bucket
	Signature(str string) string
}
type _Where interface {
	GetWhereInterface()
	PrefixStr() *string
	NextKeyStr() *string
	ExcludeStr() *string
	LimitNum() *uint64
}

func mapTrim(strs []string) []string {
	news := []string{}
	for _, str := range strs {
		str = _str.Trim(str, " ")
		if str != "" {
			news = append(news, str)
		}
	}
	return news
}
func getFileMD5(file string) (string, error) {
	f, err := _os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := _md5.New()
	if _, err := _io.Copy(h, f); err != nil {
		return "", err
	}

	return _base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}
func getFileContentType(file string) (string, error) {
	if mime, ok := _exts[_fs.Ext(file)]; ok {
		return mime, nil
	}

	f, err := _os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	buf := make([]byte, 512)
	if _, err := f.Read(buf); err != nil {
		return "", err
	}

	return _http.DetectContentType(buf), nil
}
func copy(s3 _S3, src *Bucket, dest *Bucket, args ...interface{}) error {
	if dest == nil {
		return _err.New("目的地的 Bucket 錯誤")
	}
	if dest.uri == "" {
		return _err.New(_fmt.Sprintf("目的地的沒有指定 S3 路徑"))
	}

	if src == nil {
		return _err.New("被複製的 Bucket 錯誤")
	}
	if src.uri == "" {
		return _err.New(_fmt.Sprintf("被複製的沒有指定 S3 路徑"))
	}

	acl := _enum.ACL_PRIVATE
	cache := ""

	for _, arg := range args {
		if val, ok := arg.(_enum.Acl); ok {
			acl = val
		}

		if val, ok := arg.(int); ok && val > 0 {
			cache = _fmt.Sprintf("%d", val)
		}
	}

	return _req.New(s3).Bucket(dest.name).Uri(dest.uri).Method(_enum.METHOD_PUT).SetAmzHeader("x-amz-acl", acl.Str()).SetAmzHeader("x-amz-copy-source", _fmt.Sprintf("/%s/%s", src.name, src.uri)).SetAmzHeader("x-amz-metadata-directive", "COPY").SetHeader("Cache-Control", cache).Response().IsSuccess()
}

func New(name string, s3 _S3) (*Bucket, error) {
	dirs := mapTrim(_str.Split(name, "/"))
	if len(dirs) <= 0 {
		return nil, _err.New("Bucket 名稱或路徑格式錯誤")
	}
	return &Bucket{s3: s3, name: dirs[0], uri: _str.Join(dirs[1:], "/")}, nil
}

func (bucket *Bucket) String() string {
	if bucket == nil {
		return ""
	} else {
		return bucket.name
	}
}
func (bucket *Bucket) Create(args ...interface{}) error {
	if bucket == nil {
		return _err.New("錯誤的 Bucket")
	}

	acl := _enum.ACL_PRIVATE
	loc := _enum.LOC_NANO

	for _, arg := range args {
		if val, ok := arg.(_enum.Acl); ok {
			acl = val
		}
		if val, ok := arg.(_enum.Loc); ok {
			loc = val
		}
	}

	req := _req.New(bucket.s3).Bucket(bucket.name).Method(_enum.METHOD_PUT).SetAmzHeader("x-amz-acl", acl.Str())

	if loc != _enum.LOC_NANO {
		encoder, err := _xml.MarshalIndent(struct {
			XMLName  _xml.Name `xml:"CreateBucketConfiguration"`
			Location string    `xml:"LocationConstraint"`
		}{Location: loc.Str()}, "", "  ")

		if err != nil {
			return _err.New(_fmt.Sprintf("產生 XML 失敗，Message：%s", err))
		}

		req.SetXML(_fmt.Sprintf("%s%s", _xml.Header, string(encoder)))
	}

	response := req.Response()
	if err := response.IsSuccess(); err != nil {
		return err
	}

	return nil
}
func (bucket *Bucket) Delete() error {
	if bucket == nil {
		return _err.New("錯誤的 Bucket")
	}

	err := _req.New(bucket.s3).Bucket(bucket.name).Method(_enum.METHOD_DELETE).Response().IsSuccess([]uint16{200, 204})
	return err
}
func (bucket *Bucket) Files(wheres ..._Where) ([]*_model.File, error) {
	files := []*_model.File{}

	if bucket == nil {
		return files, _err.New("錯誤的 Bucket")
	}

	var prefix, nextKey, exclude *string
	var limit *uint64

	if len(wheres) > 0 && wheres[0] != nil {
		prefix = wheres[0].PrefixStr()
		nextKey = wheres[0].NextKeyStr()
		exclude = wheres[0].ExcludeStr()
		limit = wheres[0].LimitNum()
	}

	for {
		req := _req.New(bucket.s3).Bucket(bucket.name).Method(_enum.METHOD_GET)

		if prefix != nil {
			req.Parameter("prefix", *prefix)
		}
		if nextKey != nil {
			req.Parameter("marker", *nextKey)
		}
		if exclude != nil {
			req.Parameter("delimiter", *exclude)
		}
		if limit != nil {
			req.Parameter("max-keys", _fmt.Sprintf("%d", *limit))
		}

		response := req.Response()
		if err := response.IsSuccess(); err != nil {
			return files, err
		}

		if response.Headers["Content-Type"] != "application/xml" {
			return files, _err.New("錯誤，回應結果非 XML 格式")
		}

		var result *struct {
			Name        string  `xml:"Name"`
			Prefix      string  `xml:"Prefix"`
			NextKey     string  `xml:"Marker"`
			Limit       uint64  `xml:"MaxKeys"`
			IsTruncated bool    `xml:"IsTruncated"`
			NextMarker  *string `xml:"NextMarker"`

			Contents []struct {
				Key   string `xml:"Key"`
				Time  string `xml:"LastModified"`
				ETag  string `xml:"ETag"`
				Size  uint64 `xml:"Size"`
				Owner struct {
					Id   string `xml:"ID"`
					Name string `xml:"DisplayName"`
				} `xml:"Owner"`
			} `xml:"Contents"`
		} = nil

		err := _xml.Unmarshal(response.BodyBytes, &result)
		if err != nil {
			return files, _err.New(_fmt.Sprintf("編譯 XML 失敗，Message：%s", err))
		}

		nextKey = nil
		for _, content := range result.Contents {

			time, err := _time.Parse("2006-01-02T15:04:05.999Z", content.Time)
			if err != nil {
				return nil, _err.New(_fmt.Sprintf("轉換時間格式失敗，Message：%s", err))
			}

			file := &_model.File{
				Key:  content.Key,
				Time: uint64(time.Unix()),
				Md5:  _str.Trim(content.ETag, "\""),
				Size: content.Size,
			}
			files = append(files, file)
			nextKey = &file.Key
		}

		if result.NextMarker != nil {
			nextKey = result.NextMarker
		}

		if limit != nil || nextKey == nil || result.IsTruncated == false {
			break
		}
	}

	return files, nil
}
func (bucket *Bucket) Put(path string, args ...interface{}) error {
	if bucket == nil {
		return _err.New("錯誤的 Bucket")
	}

	path, err := _fs.Abs(path)
	if err != nil {
		return _err.New(_fmt.Sprintf("無法取得 %s 檔案的絕對位置，Message：%s", path, err))
	}

	stat, err := _os.Stat(path)
	if err != nil {
		return _err.New(_fmt.Sprintf("無法取得 %s 檔案狀態，Message：%s", path, err))
	}

	if !stat.Mode().IsRegular() {
		return _err.New(_fmt.Sprintf("檔案 %s 不是常規的檔案", path))
	}

	cMd5, err := getFileMD5(path)
	if err != nil {
		return _err.New(_fmt.Sprintf("無法取得 %s 檔案的 MD5 結果，Message：%s", path, err))
	}

	cType, err := getFileContentType(path)
	if err != nil {
		return _err.New(_fmt.Sprintf("無法取得 %s 檔案的 Content Type，Message：%s", path, err))
	}

	if bucket.uri == "" {
		return _err.New(_fmt.Sprintf("S3 沒有指定路徑"))
	}

	acl := _enum.ACL_PRIVATE
	cache := ""

	for _, arg := range args {
		if val, ok := arg.(_enum.Acl); ok {
			acl = val
		}

		if val, ok := arg.(int); ok && val > 0 {
			cache = _fmt.Sprintf("%d", val)
		}
	}

	return _req.New(bucket.s3).Bucket(bucket.name).Uri(bucket.uri).Method(_enum.METHOD_PUT).SetHeader("Content-Type", cType).SetHeader("Content-MD5", cMd5).SetAmzHeader("x-amz-acl", acl.Str()).SetFile(path, uint64(stat.Size())).SetHeader("Cache-Control", cache).Response().IsSuccess()
}
func (bucket *Bucket) Del() error {
	if bucket == nil {
		return _err.New("錯誤的 Bucket")
	}

	if bucket.uri == "" {
		return _err.New(_fmt.Sprintf("S3 沒有指定路徑"))
	}

	return _req.New(bucket.s3).Bucket(bucket.name).Uri(bucket.uri).Method(_enum.METHOD_DELETE).Response().IsSuccess([]uint16{200, 204})
}
func (bucket *Bucket) Meta() (*_model.FileMeta, error) {
	if bucket == nil {
		return nil, _err.New("錯誤的 Bucket")
	}

	if bucket.uri == "" {
		return nil, _err.New(_fmt.Sprintf("S3 沒有指定路徑"))
	}

	response := _req.New(bucket.s3).Bucket(bucket.name).Uri(bucket.uri).Method(_enum.METHOD_HEAD).Response()
	if err := response.IsSuccess(); err != nil {
		return nil, err
	}

	tmp1 := uint64(0)
	if val, ok := response.Headers["Content-Length"]; ok {
		switch num, err := _strconv.ParseInt(val, 10, 64); true {
		case err != nil:
			return nil, _err.New(_fmt.Sprintf("資訊錯誤，Content-Length 格式有誤，Message：%s", err))
		case num < 0:
			return nil, _err.New(_fmt.Sprintf("資訊錯誤，Content-Length 格式有誤，其值 %d < 0：", num))
		default:
			tmp1 = uint64(num)
		}
	} else {
		return nil, _err.New("資訊有缺，缺少 Content-Length")
	}

	tmp2 := uint64(0)
	if val, ok := response.Headers["Last-Modified"]; ok {
		switch time, err := _time.Parse("Mon, 02 Jan 2006 15:04:05 GMT", val); true {
		case err != nil:
			return nil, _err.New(_fmt.Sprintf("Last-Modified 格式有誤，Message：%s", err))
		default:
			tmp2 = uint64(time.Unix())
		}
	} else {
		return nil, _err.New("資訊有缺，缺少 Last-Modified")
	}

	tmp3 := ""
	if val, ok := response.Headers["Etag"]; ok {
		tmp3 = _str.Trim(val, "\"")
	} else {
		return nil, _err.New("資訊有缺，缺少 Etag")
	}

	tmp4 := ""
	if val, ok := response.Headers["Content-Type"]; ok {
		tmp4 = val
	} else {
		return nil, _err.New("資訊有缺，缺少 Etag")
	}

	return &_model.FileMeta{
		ContentLength: tmp1,
		Time:          tmp2,
		Md5:           tmp3,
		ContentType:   tmp4,
	}, nil
}
func (bucket *Bucket) File() (*_resp.Response, error) {
	if bucket == nil {
		return nil, _err.New("錯誤的 Bucket")
	}

	if bucket.uri == "" {
		return nil, _err.New(_fmt.Sprintf("S3 沒有指定路徑"))
	}

	response := _req.New(bucket.s3).Bucket(bucket.name).Uri(bucket.uri).Method(_enum.METHOD_GET).Response()

	if err := response.IsSuccess(); err != nil {
		return nil, err
	}

	return response, nil
}
func (bucket *Bucket) Save(path string, modes ..._os.FileMode) error {
	resp, err := bucket.File()
	if err != nil {
		return err
	}

	path, err = _fs.Abs(path)
	if err != nil {
		return _err.New(_fmt.Sprintf("無法取得 %s 檔案的絕對位置，Message：%s", path, err))
	}

	file, err := _os.Create(path)
	if err != nil {
		return _err.New(_fmt.Sprintf("無法取得 %s 檔案的資源，Message：%s", path, err))
	}

	defer file.Close()

	_, err = file.Write(resp.BodyBytes)
	if err != nil {
		return _err.New(_fmt.Sprintf("%s 檔案寫入失敗，Message：%s", path, err))
	}
	file.Sync()

	var mode _os.FileMode = 0644
	if len(modes) > 0 {
		mode = modes[0]
	}
	err = _os.Chmod(path, mode)
	if err != nil {
		return _err.New(_fmt.Sprintf("%s 檔案變更權限失敗，Message：%s", path, err))
	}

	return nil
}
func (bucket *Bucket) CopyTo(dest string, args ...interface{}) error {
	return copy(bucket.s3, bucket, bucket.s3.Bucket(dest))
}
func (bucket *Bucket) CopyFrom(src string, args ...interface{}) error {
	return copy(bucket.s3, bucket.s3.Bucket(src), bucket)
}
func (bucket *Bucket) Clean(cpus ...uint8) []error {
	files, err := bucket.Files()
	if err != nil {
		return []error{err}
	}

	cpu := 1
	if len(cpus) > 0 {
		cpu = int(cpus[0])
	}

	total := len(files)
	wg := new(_sync.WaitGroup)
	ins := make(chan string, total)
	ous := make(chan error, total)

	for i := 0; i < cpu; i++ {
		wg.Add(1)
		go func(s3 _S3, wg *_sync.WaitGroup, ins <-chan string, ous chan<- error) {
			defer wg.Done()

			for in := range ins {
				if err := s3.Bucket(in).Del(); err != nil {
					ous <- _err.New(_fmt.Sprintf("刪除檔案 %s 時發生錯誤，Message：%s", in, err))
				} else {
					ous <- nil
				}
			}
		}(bucket.s3, wg, ins, ous)
	}

	for _, file := range files {
		ins <- _fmt.Sprintf("%s/%s", bucket.name, file.Key)
	}
	close(ins)
	wg.Wait()

	errs := []error{}
	for i := 0; i < total; i++ {
		if err := <-ous; err != nil {
			errs = append(errs, err)
		}
	}

	close(ous)

	return errs
}
