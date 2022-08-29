/**
 * @author      OA Wu <oawu.tw@gmail.com>
 * @copyright   Copyright (c) 2015 - 2021
 * @license     http://opensource.org/licenses/MIT  MIT License
 * @link        https://www.ioa.tw/
 */

package s3

type Where struct {
	Prefix  string
	NextKey string
	Exclude string
	Limit   uint64
}
func (where Where) GetWhereInterface() {}
func (where Where) PrefixStr() *string {
	if where.Prefix == "" {
		return nil
	}
	return &where.Prefix
}
func (where Where) NextKeyStr() *string {
	if where.NextKey == "" {
		return nil
	}
	return &where.NextKey
}
func (where Where) ExcludeStr() *string {
	if where.Exclude == "" {
		return nil
	}
	return &where.Exclude
}
func (where Where) LimitNum() *uint64 {
	if where.Limit == 0 {
		return nil
	}
	return &where.Limit
}