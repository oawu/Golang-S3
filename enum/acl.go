/**
 * @author      OA Wu <oawu.tw@gmail.com>
 * @copyright   Copyright (c) 2015 - 2021
 * @license     http://opensource.org/licenses/MIT  MIT License
 * @link        https://www.ioa.tw/
 */

package enum

type Acl int

const (
	ACL_PRIVATE Acl = iota
	ACL_PUBLIC_READ
	ACL_PUBLIC_READ_WRITE
	ACL_AWS_EXEC_READ
	ACL_AUTHENTICATED_READ
	ACL_BUCKET_OWNER_READ
	ACL_BUCKET_OWNER_FULL_CONTROL
	ACL_LOG_DELIVERY_WRITE
)

func (acl Acl) Str() string {
	switch acl {
	case ACL_PUBLIC_READ:
		return "public-read"
	case ACL_PUBLIC_READ_WRITE:
		return "public-read-write"
	case ACL_AWS_EXEC_READ:
		return "aws-exec-read"
	case ACL_AUTHENTICATED_READ:
		return "authenticated-read"
	case ACL_BUCKET_OWNER_READ:
		return "bucket-owner-read"
	case ACL_BUCKET_OWNER_FULL_CONTROL:
		return "bucket-owner-full-control"
	case ACL_LOG_DELIVERY_WRITE:
		return "log-delivery-write"
	default:
		return "private"
	}
}

// https://docs.aws.amazon.com/AmazonS3/latest/userguide/acl-overview.html

// ACL_PRIVATE                   | Owner gets FULL_CONTROL. No one else has access rights (default).
// ACL_PUBLIC_READ               | Owner gets FULL_CONTROL. The AllUsers group (see Who is a grantee?) gets READ access.
// ACL_PUBLIC_READ_WRITE         | Owner gets FULL_CONTROL. The AllUsers group gets READ and WRITE access. Granting this on a bucket is generally not recommended.
// ACL_AWS_EXEC_READ             | Owner gets FULL_CONTROL. Amazon EC2 gets READ access to GET an Amazon Machine Image (AMI) bundle from Amazon S3.
// ACL_AUTHENTICATED_READ        | Owner gets FULL_CONTROL. The AuthenticatedUsers group gets READ access.
// ACL_BUCKET_OWNER_READ         | Object owner gets FULL_CONTROL. Bucket owner gets READ access. If you specify this canned ACL when creating a bucket, Amazon S3 ignores it.
// ACL_BUCKET_OWNER_FULL_CONTROL | Both the object owner and the bucket owner get FULL_CONTROL over the object. If you specify this canned ACL when creating a bucket, Amazon S3 ignores it.
// ACL_LOG_DELIVERY_WRITE        | The LogDelivery group gets WRITE and READ_ACP permissions on the bucket. For more information about logs, see (Logging requests using server access logging).
