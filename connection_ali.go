package connections

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
)

type ConnectionsAli struct {
	ECS *ecs.Client
	ESS *ess.Client
	CS  *cs.Client
	RAM *ram.Client
}

func NewAli(region string, accessKeyId string, accessKeySecret string) *ConnectionsAli {
	conn := ConnectionsAli{}
	conn.ConnectAli(region, accessKeyId, accessKeySecret)
	return &conn
}

func (c *ConnectionsAli) ConnectAli(region string, accessKeyId string, accessKeySecret string) {
	c.ECS, _ = ecs.NewClientWithAccessKey(region, accessKeyId, accessKeySecret)
	c.ESS, _ = ess.NewClientWithAccessKey(region, accessKeyId, accessKeySecret)
	c.CS,_ = cs.NewClientWithAccessKey(region,accessKeyId,accessKeySecret)
	c.RAM,_ = ram.NewClientWithAccessKey(region, accessKeyId, accessKeySecret)
}
