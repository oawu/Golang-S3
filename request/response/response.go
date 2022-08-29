/**
 * @author      OA Wu <oawu.tw@gmail.com>
 * @copyright   Copyright (c) 2015 - 2022
 * @license     http://opensource.org/licenses/MIT  MIT License
 * @link        https://www.ioa.tw/
 */

package response

import (
	_err "errors"
	_fmt "fmt"
	_str "strings"
)

type Response struct {
	StatusCode uint16
	BodyBytes  []byte
	BodyString string
	Error      error
	Headers    map[string]string
}

func New() *Response {
	resp := &Response{}
	return resp.Reset()
}
func (resp *Response) GG(err error) *Response {
	if resp == nil {
		resp = &Response{StatusCode: 0, BodyBytes: []byte{}, BodyString: "", Error: err, Headers: map[string]string{}}
	}

	resp.Error = err
	return resp
}
func (resp *Response) Reset() *Response {
	if resp == nil {
		return resp
	}
	resp.StatusCode = uint16(0)
	resp.BodyBytes = []byte{}
	resp.BodyString = ""
	resp.Error = nil
	resp.Headers = map[string]string{}
	return resp
}

func (resp *Response) IsSuccess(statuss ...[]uint16) error {
	if resp == nil {
		return _err.New("錯誤的 Response")
	}

	if resp.Error != nil {
		return resp.Error
	}

	status := []uint16{200}
	if len(statuss) > 0 {
		status = statuss[0]
	}

	for _, s := range status {
		if resp.StatusCode == s {
			return nil
		}
	}

	tmps := []string{}
	for _, s := range status {
		tmps = append(tmps, _fmt.Sprintf("%d", s))
	}

	return _err.New(_fmt.Sprintf("錯誤，狀態非 %s", _str.Join(tmps, "、")))
}
