/**
 * @author      OA Wu <oawu.tw@gmail.com>
 * @copyright   Copyright (c) 2015 - 2022
 * @license     http://opensource.org/licenses/MIT  MIT License
 * @link        https://www.ioa.tw/
 */

package enum

type Loc int

const (
	LOC_NANO Loc = iota
	LOC_AF_SOUTH_1
	LOC_AP_EAST_1
	LOC_AP_NORTHEAST_1
	LOC_AP_NORTHEAST_2
	LOC_AP_NORTHEAST_3
	LOC_AP_SOUTH_1
	LOC_AP_SOUTHEAST_1
	LOC_AP_SOUTHEAST_2
	LOC_CA_CENTRAL_1
	LOC_CN_NORTH_1
	LOC_CN_NORTHWEST_1
	LOC_EU
	LOC_EU_CENTRAL_1
	LOC_EU_NORTH_1
	LOC_EU_SOUTH_1
	LOC_EU_WEST_1
	LOC_EU_WEST_2
	LOC_EU_WEST_3
	LOC_ME_SOUTH_1
	LOC_SA_EAST_1
	LOC_US_EAST_2
	LOC_US_GOV_EAST_1
	LOC_US_GOV_WEST_1
	LOC_US_WEST_1
	LOC_US_WEST_2
)

func (loc Loc) Str() string {
	switch loc {
	case LOC_AF_SOUTH_1:
		return "af-south-1"
	case LOC_AP_EAST_1:
		return "ap-east-1"
	case LOC_AP_NORTHEAST_1:
		return "ap-northeast-1"
	case LOC_AP_NORTHEAST_2:
		return "ap-northeast-2"
	case LOC_AP_NORTHEAST_3:
		return "ap-northeast-3"
	case LOC_AP_SOUTH_1:
		return "ap-south-1"
	case LOC_AP_SOUTHEAST_1:
		return "ap-southeast-1"
	case LOC_AP_SOUTHEAST_2:
		return "ap-southeast-2"
	case LOC_CA_CENTRAL_1:
		return "ca-central-1"
	case LOC_CN_NORTH_1:
		return "cn-north-1"
	case LOC_CN_NORTHWEST_1:
		return "cn-northwest-1"
	case LOC_EU:
		return "EU"
	case LOC_EU_CENTRAL_1:
		return "eu-central-1"
	case LOC_EU_NORTH_1:
		return "eu-north-1"
	case LOC_EU_SOUTH_1:
		return "eu-south-1"
	case LOC_EU_WEST_1:
		return "eu-west-1"
	case LOC_EU_WEST_2:
		return "eu-west-2"
	case LOC_EU_WEST_3:
		return "eu-west-3"
	case LOC_ME_SOUTH_1:
		return "me-south-1"
	case LOC_SA_EAST_1:
		return "sa-east-1"
	case LOC_US_EAST_2:
		return "us-east-2"
	case LOC_US_GOV_EAST_1:
		return "us-gov-east-1"
	case LOC_US_GOV_WEST_1:
		return "us-gov-west-1"
	case LOC_US_WEST_1:
		return "us-west-1"
	case LOC_US_WEST_2:
		return "us-west-2"
	default:
		return ""
	}
}
