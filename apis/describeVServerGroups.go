package apis

import (
	slb20140515 "github.com/alibabacloud-go/slb-20140515/v4/client"
)

func DescribeVServerGroups(regionId *string, loadBalancerId *string, c *slb20140515.Client) (_result *slb20140515.DescribeVServerGroupsResponse, _err error) {
	request := slb20140515.DescribeVServerGroupsRequest{}
	request.SetRegionId(*regionId)
	request.SetLoadBalancerId(*loadBalancerId)
	return c.DescribeVServerGroups(&request)
}
