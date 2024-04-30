package main

import (
	"aliyun-clb/client"
	"encoding/json"
	"flag"
	"fmt"
	slb20140515 "github.com/alibabacloud-go/slb-20140515/v4/client"
	"os"
	"path"
)

func main() {
	regionId := flag.String("regionId", "", "传统型负载均衡实例的地域 ID")
	out := flag.String("out", "", "输出目录")
	sourceLoadBalancerId := flag.String("sourceLoadBalancerId", "", "源负载均衡id")
	targetLoadBalancerId := flag.String("targetLoadBalancerId", "", "目标负载均衡id")
	listenerPort := flag.Int("listenerPort", 0, "负载均衡实例前端使用的监听端口，取值范围：1~65535。")
	listenerProtocol := flag.String("listenerProtocol", "", "负载均衡实例前端使用的协议。")
	flag.Parse()

	c, err := client.CreateClient()
	if err != nil {
		panic(err)
	}

	// 查询监听所有转发规则
	request := slb20140515.DescribeRulesRequest{}
	request.SetRegionId(*regionId)
	request.SetLoadBalancerId(*sourceLoadBalancerId)
	request.SetListenerPort(int32(*listenerPort))
	request.SetListenerProtocol(*listenerProtocol)
	rsp, err := c.DescribeRules(&request)
	if err != nil {
		panic(err)
	}

	createRulesRequest := slb20140515.CreateRulesRequest{}
	createRulesRequest.SetRegionId(*regionId)
	createRulesRequest.SetLoadBalancerId(*targetLoadBalancerId)
	createRulesRequest.SetListenerPort(int32(*listenerPort))
	createRulesRequest.SetListenerProtocol(*listenerProtocol)

	ruleLen := len(rsp.Body.Rules.Rule)
	roles := make([]slb20140515.DescribeRuleAttributeResponseBody, ruleLen)
	// 遍历监听转发规则，查询详情并保存
	// groupPath := path.Join(*out, "groups")
	for _, rule := range rsp.Body.Rules.Rule {
		detailRequest := slb20140515.DescribeRuleAttributeRequest{}
		detailRequest.SetRegionId(*regionId)
		detailRequest.SetRuleId(*rule.RuleId)
		detailRsp, err := c.DescribeRuleAttribute(&detailRequest)
		if err != nil {
			panic(err)
		}
		r := slb20140515.DescribeRuleAttributeResponseBody{}
		// RuleName
		r.SetRuleName(*detailRsp.Body.RuleName)
		// Domain
		r.SetDomain(*detailRsp.Body.Domain)
		// Url
		r.SetUrl(*detailRsp.Body.Url)
		// Domain
		r.SetDomain(*detailRsp.Body.Domain)
		// VServerGroupId
		// 获取并保存虚拟服务器组
		// rsp, err := backup.GroupBackup(regionId, r.VServerGroupId, &groupPath, targetLoadBalancerId, c)
		r.SetVServerGroupId(*detailRsp.Body.Domain)
		roles = append(roles, r)
	}
	b, err := json.Marshal(roles)
	if err != nil {
		panic(err)
	}
	createRulesRequest.SetRuleList(string(b))

	_, e := os.Stat(*out)
	if e != nil {
		e = os.Mkdir(*out, os.FileMode(0755))
		if e != nil {
			panic(e)
		}
	}
	outfile := path.Join(*out, fmt.Sprintf("%s-%d.json", *listenerProtocol, *listenerPort))
	a, _ := json.Marshal(createRulesRequest)
	err = os.WriteFile(outfile, a, os.FileMode(0755))
	if err != nil {
		panic(err)
	}
	fmt.Printf("保存成功：%s", outfile)
}
