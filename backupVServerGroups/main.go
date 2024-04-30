package main

import (
	"aliyun-clb/apis"
	"aliyun-clb/backupVServerGroups/backup"
	"aliyun-clb/client"
	"flag"
	"fmt"
)

func main() {
	c, err := client.CreateClient()
	if err != nil {
		panic(err)
	}
	regionId := flag.String("regionId", "", "传统型负载均衡实例的地域 ID")
	out := flag.String("out", "", "输出目录")
	sourceLoadBalancerId := flag.String("sourceLoadBalancerId", "", "源负载均衡id")
	targetLoadBalancerId := flag.String("targetLoadBalancerId", "", "目标负载均衡id")
	flag.Parse()

	if *regionId == "" {
		panic("regionId 不能为空")
	}

	if *sourceLoadBalancerId == "" {
		panic("loadBalancerId 不能为空")
	}

	if err != nil {
		panic(err)
	}
	rsp, err := apis.DescribeVServerGroups(regionId, sourceLoadBalancerId, c)
	if err != nil {
		panic(err)
	}
	for _, group := range rsp.Body.VServerGroups.VServerGroup {
		fmt.Printf("VServerGroupId：%s，VServerGroupName：%s，targetLoadBalancerId：%s\n", *group.VServerGroupId, *group.VServerGroupName, *targetLoadBalancerId)
		rsp, err := backup.GroupBackup(regionId, group.VServerGroupId, out, targetLoadBalancerId, c)
		if err != nil {
			fmt.Printf("获取结果失败：%s", *group.VServerGroupName)
		} else {
			fmt.Printf("成功：%s", *rsp.Body.VServerGroupName)
		}
	}
}
