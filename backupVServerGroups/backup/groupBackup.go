package backup

import (
	"aliyun-clb/apis"
	"encoding/json"
	"fmt"
	slb20140515 "github.com/alibabacloud-go/slb-20140515/v4/client"
	"os"
	"path"
)

func GroupBackup(regionId *string, groupId *string, out *string, loadBalancerId *string,
	c *slb20140515.Client) (_result *slb20140515.DescribeVServerGroupAttributeResponse, _err error) {
	rs, err := apis.DescribeVServerGroupAttribute(*regionId, *groupId, c)
	if err != nil {
		fmt.Printf("处理失败，原因：%s\n", err.Error())
	}
	obj := &slb20140515.CreateVServerGroupRequest{}
	if *loadBalancerId != "" {
		obj.SetLoadBalancerId(*loadBalancerId)
	} else {
		obj.SetLoadBalancerId(*rs.Body.LoadBalancerId)
	}
	obj.SetRegionId(*regionId)
	obj.SetVServerGroupName(*rs.Body.VServerGroupName)
	servers, _ := json.Marshal(&rs.Body.BackendServers.BackendServer)
	obj.SetBackendServers(string(servers))
	_, e := os.Stat(*out)
	if e != nil {
		e = os.Mkdir(*out, os.FileMode(0755))
		if e != nil {
			return nil, e
		}
	}
	outfile := path.Join(*out, *rs.Body.VServerGroupName+".json")
	a, _ := json.Marshal(obj)
	err = os.WriteFile(outfile, a, os.FileMode(0755))
	if err != nil {
		return nil, e
	}
	return rs, nil
}
