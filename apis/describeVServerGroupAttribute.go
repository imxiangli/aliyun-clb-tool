package apis

import (
	slb20140515 "github.com/alibabacloud-go/slb-20140515/v4/client"
)

func DescribeVServerGroupAttribute(regionId string, vServerGroupId string, client *slb20140515.Client) (rs *slb20140515.DescribeVServerGroupAttributeResponse, err error) {
	request := &slb20140515.DescribeVServerGroupAttributeRequest{}
	request.SetRegionId(regionId)
	request.SetVServerGroupId(vServerGroupId)
	return client.DescribeVServerGroupAttribute(request)
}
