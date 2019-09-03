package connections

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/spotmaxtech/gokit"
	"testing"
	"time"
)

func TestNewTemporaryCredentials(t *testing.T) {
	cre := NewTemporaryCredentials(&sts.AssumeRoleInput{
		DurationSeconds: aws.Int64(900),
		RoleArn:         aws.String("arn:aws:iam::**********:role/mobvista_connect"),
		RoleSessionName: aws.String("mobvista"),
	})

	t.Log(cre.Get())
	t1, _ := cre.ExpiresAt()
	t.Log(t1, cre.IsExpired())

	time.Sleep(time.Second * 1)

	cre.Expire()
	t.Log(cre.IsExpired())

	t.Log(cre.Get())
	t2, _ := cre.ExpiresAt()
	t.Log(t2, cre.IsExpired())
}

func TestTemporaryCredentialProvider(t *testing.T) {
	cre := NewTemporaryCredentials(&sts.AssumeRoleInput{
		DurationSeconds: aws.Int64(900),
		RoleArn:         aws.String("arn:aws:iam::**********:role/mobvista_connect"),
		RoleSessionName: aws.String("mobvista"),
	})

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: cre,
	})
	if err != nil {
		t.Fatal("session error", err.Error())
	}

	ec2Svc := ec2.New(sess)
	output1, err := ec2Svc.DescribeRegions(&ec2.DescribeRegionsInput{})
	if err != nil {
		t.Log("describe 1 error", err.Error())
	}
	t.Log(output1)

	// sleep to session expire
	timer := gokit.NewSimpleTimer(time.Second * 910)
	for !timer.Timeout() {
		time.Sleep(time.Second * 30)
		t.Log("sleep ...")
	}

	// force expire
	// cre.Expire()

	// again
	output2, err := ec2Svc.DescribeRegions(&ec2.DescribeRegionsInput{})
	if err != nil {
		t.Log("describe 2 error", err.Error())
	}
	t.Log(output2)

}
