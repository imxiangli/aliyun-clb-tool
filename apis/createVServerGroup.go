package apis

import (
	"encoding/json"
	"errors"
	"fmt"
	slb20140515 "github.com/alibabacloud-go/slb-20140515/v4/client"
	"os"
)

func CreateVServerGroup(f string, client *slb20140515.Client) (rs *slb20140515.CreateVServerGroupResponse, err error) {
	bytes, err := os.ReadFile(f)
	if err != nil {
		return nil, errors.New("读取文件失败")
	}
	jsonStr := string(bytes)
	fmt.Printf("读取到文件内容：%s\n", string(bytes))
	request := &slb20140515.CreateVServerGroupRequest{}
	err = json.Unmarshal([]byte(jsonStr), &request)
	if err != nil {
		return nil, err
	}
	return client.CreateVServerGroup(request)
}
