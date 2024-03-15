package connections

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/aws/aws-sdk-go/service/costexplorer/costexploreriface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/pricing"
	"github.com/aws/aws-sdk-go/service/pricing/pricingiface"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
)

type Connections struct {
	Session        *session.Session
	AutoScaling    autoscalingiface.AutoScalingAPI
	EC2            ec2iface.EC2API
	ELB            elbiface.ELBAPI
	ELBV2          elbv2iface.ELBV2API
	CostExplorer   costexploreriface.CostExplorerAPI
	CloudFormation cloudformationiface.CloudFormationAPI
	CloudWatch     cloudwatchiface.CloudWatchAPI
	Pricing        pricingiface.PricingAPI
	Region         string
	IAM            iamiface.IAMAPI
	Route53        route53iface.Route53API
}

func New(region string) *Connections {
	conn := Connections{}
	conn.Connect(region)
	return &conn
}

func NewAWS(region string, accessKeyId string, secretAccessKey string) *Connections {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKeyId, secretAccessKey, ""),
	}))
	conn := Connections{Session: sess}
	conn.Connect(*conn.Session.Config.Region)
	return &conn
}

// new from session when session provided
func NewFromSession(sess *session.Session) *Connections {
	conn := Connections{Session: sess}
	conn.Connect(*conn.Session.Config.Region)
	return &conn
}

func (c *Connections) SetSession(region string) {
	c.Session = session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))
}

// AWS Price List Service API provides the following two endpoints:
// https://api.pricing.us-east-1.amazonaws.com
// https://api.pricing.ap-south-1.amazonaws.com
func (c *Connections) PricingSession() *session.Session {
	return session.Must(session.NewSession(&aws.Config{Region: aws.String("us-east-1")}))
}

func (c *Connections) Connect(region string) {
	if c.Session == nil {
		c.SetSession(region)
	}

	asConn := make(chan *autoscaling.AutoScaling)
	ec2Conn := make(chan *ec2.EC2)
	elbConn := make(chan *elb.ELB)
	elbv2Conn := make(chan *elbv2.ELBV2)
	costExplorerConn := make(chan *costexplorer.CostExplorer)
	cloudFormationConn := make(chan *cloudformation.CloudFormation)
	cloudWatchConn := make(chan *cloudwatch.CloudWatch)
	pricingConn := make(chan *pricing.Pricing)
	iamConn := make(chan *iam.IAM)
	route53Conn := make(chan *route53.Route53)

	go func() { asConn <- autoscaling.New(c.Session) }()
	go func() { ec2Conn <- ec2.New(c.Session) }()
	go func() { elbConn <- elb.New(c.Session) }()
	go func() { elbv2Conn <- elbv2.New(c.Session) }()
	go func() { costExplorerConn <- costexplorer.New(c.Session) }()
	go func() { cloudFormationConn <- cloudformation.New(c.Session) }()
	go func() { cloudWatchConn <- cloudwatch.New(c.Session) }()
	go func() { pricingConn <- pricing.New(c.PricingSession()) }()
	go func() { iamConn <- iam.New(c.Session) }()
	go func() { route53Conn <- route53.New(c.Session) }()

	c.AutoScaling, c.EC2, c.ELB, c.ELBV2, c.CostExplorer, c.CloudFormation, c.CloudWatch, c.Pricing, c.Region, c.IAM, c.Route53 =
		<-asConn, <-ec2Conn, <-elbConn, <-elbv2Conn, <-costExplorerConn, <-cloudFormationConn, <-cloudWatchConn, <-pricingConn, region, <-iamConn, <-route53Conn
}
