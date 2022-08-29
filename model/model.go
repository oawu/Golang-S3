/**
 * @author      OA Wu <oawu.tw@gmail.com>
 * @copyright   Copyright (c) 2015 - 2021
 * @license     http://opensource.org/licenses/MIT  MIT License
 * @link        https://www.ioa.tw/
 */

package model

type BucketInfo struct {
	Owner *BucketInfoOwner

	Buckets []BucketInfoBucket
}
type BucketInfoOwner struct {
	Id   string
	Name string
}
type BucketInfoBucket struct {
	Name string
	Time uint64
}
type File struct {
	Key  string
	Time uint64
	Md5  string
	Size uint64
}
type FileMeta struct {
	ContentLength uint64
	Time          uint64
	Md5           string
	ContentType   string
}
