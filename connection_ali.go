package connections

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
)

type ConnectionsAli struct {
	ecs *ecs.Client
	ess *ess.Client
	cs  *cs.Client
}

func NewAli(region string, accessKeyId string, accessKeySecret string) *ConnectionsAli {
	conn := ConnectionsAli{}
	conn.ConnectAli(region, accessKeyId, accessKeySecret)
	return &conn
}

func (c *ConnectionsAli) ConnectAli(region string, accessKeyId string, accessKeySecret string) {
	c.ecs, _ = ecs.NewClientWithAccessKey(region, accessKeyId, accessKeySecret)
	c.ess, _ = ess.NewClientWithAccessKey(region, accessKeyId, accessKeySecret)
	c.cs,_ = cs.NewClientWithAccessKey(region,accessKeyId,accessKeySecret)
}