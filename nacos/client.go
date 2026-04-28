package nacos

import (
	"errors"
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

var (
	configClient config_client.IConfigClient
	namingClient naming_client.INamingClient
)

func InitNacosClient(config *NacosConfig) error {
	if config == nil {
		return errors.New("nacos config is nil")
	}
	if config.Enabled == false {
		return nil
	}

	serverConfigs, clientConfig := config.getConfig()
	var err error
	// 创建配置客户端
	configClient, err = clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		return err
	}

	namingClient, err = clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func GetConfig(group, dataId string) (string, error) {
	return configClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
}

func ListenConfig(group, dataId string, listener func(namespace, group, dataId, data string)) error {
	return configClient.ListenConfig(vo.ConfigParam{
		DataId:   dataId,
		Group:    group,
		OnChange: listener,
	})
}

type ServiceInfo struct {
	Name      string
	GroupName string
	Clusters  string
	Host      string
	Port      uint64
	Weight    float64
	Metadata  map[string]string
}

func RegisterInstance(info ServiceInfo) error {
	if info.Host == "" {
		localIP, err := GetLocalIP()
		if err != nil {
			return errors.New("the service host cannot be empty")
		}
		info.Host = localIP
	}

	success, err := namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          info.Host,
		Port:        info.Port,
		Weight:      info.Weight,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true, // 临时实例，服务断开后会自动删除（推荐）
		Metadata:    info.Metadata,
		ServiceName: info.Name,
		GroupName:   info.GroupName,
		ClusterName: info.Clusters,
	})
	if !success || err != nil {
		return errors.New(fmt.Sprintf("Service registration failed. error: %v", err))
	}
	return nil
}

func DeregisterInstance(name, groupName, clusters string) (bool, error) {
	return namingClient.DeregisterInstance(vo.DeregisterInstanceParam{
		ServiceName: name,
		GroupName:   groupName,
		Cluster:     clusters,
		Ephemeral:   true,
	})
}
