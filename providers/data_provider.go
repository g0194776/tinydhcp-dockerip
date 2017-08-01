package providers

import "fmt"

type DataProvider interface {
	Initialize(ips []string, env_id int, otherargs string, needInitData bool) error
	GetAvailableIP(nodeIp string, envId int, owner, desc string) (string, error)
}

func GetDataProvider(provider string) (DataProvider, error) {
	switch provider {
	case "mysql":
		return &MySqlDataProvider{}, nil
	default:
		return nil, fmt.Errorf("Unsupported data provider type: %s", provider)
	}
}
