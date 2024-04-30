package client

import (
	"errors"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	slb20140515 "github.com/alibabacloud-go/slb-20140515/v4/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/spf13/viper"
)

// CreateClient /**
func CreateClient() (_result *slb20140515.Client, _err error) {
	appConfig, err := getConfigs()
	if err != nil {
		return nil, err
	}
	accessKeyId := appConfig.GetString("app.access-key-id")
	accessKeySecret := appConfig.GetString("app.access-key-secret")
	config := &openapi.Config{
		AccessKeyId:     &accessKeyId,
		AccessKeySecret: &accessKeySecret,
	}
	config.Endpoint = tea.String("slb.aliyuncs.com")
	_result = &slb20140515.Client{}
	_result, _err = slb20140515.NewClient(config)
	return _result, _err
}

func getConfigs() (config *viper.Viper, _err error) {
	config = viper.New()
	config.AddConfigPath("./config/")
	config.SetConfigName("app")
	config.SetConfigType("yaml")
	if err := config.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			return nil, errors.New("找不到配置文件")
		} else {
			return nil, errors.New("读取配置文件出错")
		}
	}
	return config, nil
}
