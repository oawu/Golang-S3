/**
 * @author      OA Wu <oawu.tw@gmail.com>
 * @copyright   Copyright (c) 2015 - 2021
 * @license     http://opensource.org/licenses/MIT  MIT License
 * @link        https://www.ioa.tw/
 */

package enum

type Method int

const (
	METHOD_GET Method = iota
	METHOD_POST
	METHOD_PUT
	METHOD_DELETE
	METHOD_HEAD
)

func (method Method) Str() string {
	switch method {
	case METHOD_HEAD:
		return "HEAD"
	case METHOD_POST:
		return "POST"
	case METHOD_PUT:
		return "PUT"
	case METHOD_DELETE:
		return "DELETE"
	default:
		return "GET"
	}
}
