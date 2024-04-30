package main

import (
	"aliyun-clb/apis"
	"aliyun-clb/client"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	slb20140515 "github.com/alibabacloud-go/slb-20140515/v4/client"
	"strconv"
	"strings"
)

func main() {
	sourceRegionId := flag.String("sourceRegionId", "", "源地域id")
	sourceLoadBalancerId := flag.String("sourceLoadBalancerId", "", "源负载均衡id")
	sourceListener := flag.String("sourceListener", "", "源监听端口协议。例如：https:443")
	targetRegionId := flag.String("targetRegionId", "", "目标地域id")
	targetLoadBalancerId := flag.String("targetLoadBalancerId", "", "目标负载均衡id")
	targetListener := flag.String("targetListener", "", "目标监听端口协议。例如：https:443")
	flag.Parse()
	fmt.Println(*sourceRegionId, *sourceLoadBalancerId, *sourceListener, *targetRegionId, *targetLoadBalancerId, *targetListener)

	c, err := client.CreateClient()
	if err != nil {
		panic(err)
	}

	// 复制虚拟服务器组，并保存源虚拟服务器组和新虚拟服务器组的映射map
	mapping := copyGroups(sourceRegionId, sourceLoadBalancerId, targetRegionId, targetLoadBalancerId, c)
	fmt.Println(mapping)

	// 查询源监听规则列表
	sourcePort, sourceProtocol := extractListener(*sourceListener)
	rulesRspBody := listRules(sourceRegionId, sourceLoadBalancerId, &sourcePort, &sourceProtocol, c)
	rules := rulesRspBody.Rules.Rule

	targetPort, targetProtocol := extractListener(*targetListener)
	createRulesRequest := slb20140515.CreateRulesRequest{}
	createRulesRequest.SetRegionId(*targetRegionId)
	createRulesRequest.SetLoadBalancerId(*targetLoadBalancerId)
	createRulesRequest.SetListenerPort(targetPort)
	createRulesRequest.SetListenerProtocol(targetProtocol)
	targetRoles := make([]slb20140515.DescribeRuleAttributeResponseBody, len(rules))
	domainMapping := make(map[string]string)
	for i, rule := range rules {
		r := slb20140515.DescribeRuleAttributeResponseBody{}
		// RuleName
		r.SetRuleName(*rule.RuleName)
		// Domain
		domain, exists := domainMapping[*rule.Domain]
		if !exists {
			fmt.Printf("请输入对应源 %s 的域名：", *rule.Domain)
			_, err := fmt.Scanln(&domain)
			if err != nil {
				panic(err)
			}
			domainMapping[*rule.Domain] = domain
		}
		r.SetDomain(domain)
		// Url
		if rule.Url != nil {
			r.SetUrl(*rule.Url)
		}
		// VServerGroupId
		if rule.VServerGroupId != nil {
			r.SetVServerGroupId(mapping[*rule.VServerGroupId])
		}
		targetRoles[i] = r
	}
	b, err := json.Marshal(targetRoles)
	if err != nil {
		panic(err)
	}
	ruleList := string(b)
	createRulesRequest.SetRuleList(ruleList)
	// 保存目标监听转发规则
	_, err = c.CreateRules(&createRulesRequest)
	if err != nil {
		panic(err)
	}
}

func extractListener(listener string) (port int32, protocol string) {
	parts := strings.Split(listener, ":")
	if len(parts) == 2 {
		port, err := strconv.Atoi(parts[1])
		if err != nil {
			panic(err)
		}
		return int32(port), parts[0]
	}
	panic(errors.New("不能解析监听参数：" + listener))
}

func listRules(regionId *string, loadBalancerId *string,
	listenerPort *int32, listenerProtocol *string, c *slb20140515.Client) *slb20140515.DescribeRulesResponseBody {
	// 查询监听所有转发规则
	request := slb20140515.DescribeRulesRequest{}
	request.SetRegionId(*regionId)
	request.SetLoadBalancerId(*loadBalancerId)
	request.SetListenerPort(*listenerPort)
	request.SetListenerProtocol(*listenerProtocol)
	rsp, err := c.DescribeRules(&request)
	if err != nil {
		panic(err)
	}
	return rsp.Body
}

func copyGroups(sourceRegionId *string, sourceLoadBalancerId *string, targetRegionId *string,
	targetLoadBalancerId *string, c *slb20140515.Client) (_mapping map[string]string) {
	rsp, err := apis.DescribeVServerGroups(sourceRegionId, sourceLoadBalancerId, c)
	if err != nil {
		panic(err)
	}
	groups := rsp.Body.VServerGroups.VServerGroup
	mapping := make(map[string]string, len(groups))
	for _, group := range groups {
		rs, err := apis.DescribeVServerGroupAttribute(*sourceRegionId, *group.VServerGroupId, c)
		if err != nil {
			panic(err)
		}
		createVServerGroupRequest := &slb20140515.CreateVServerGroupRequest{}
		createVServerGroupRequest.SetLoadBalancerId(*targetLoadBalancerId)
		createVServerGroupRequest.SetRegionId(*targetRegionId)
		createVServerGroupRequest.SetVServerGroupName(*group.VServerGroupName)
		if len(rs.Body.BackendServers.BackendServer) > 0 {
			servers, _ := json.Marshal(&rs.Body.BackendServers.BackendServer)
			createVServerGroupRequest.SetBackendServers(string(servers))
		}
		rsp, err := c.CreateVServerGroup(createVServerGroupRequest)
		if err != nil {
			panic(err)
		}
		mapping[*group.VServerGroupId] = *rsp.Body.VServerGroupId
	}
	return mapping
}
