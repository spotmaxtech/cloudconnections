package connections

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"time"
)

// Provider for role based credential
type TemporaryCredentialProvider struct {
	assumeRoleInput *sts.AssumeRoleInput
	expiresAt       *time.Time
}

// Create a temporary credential with provider
// * duration: timer duration
// * autoRefresh: enable auto refresh or not
func NewTemporaryCredentials(assumeRoleInput *sts.AssumeRoleInput) *credentials.Credentials {
	return credentials.NewCredentials(&TemporaryCredentialProvider{
		assumeRoleInput: assumeRoleInput,
	})
}

func (p *TemporaryCredentialProvider) Retrieve() (credentials.Value, error) {
	ses, err := session.NewSession(&aws.Config{})
	svc := sts.New(ses)

	result, err := svc.AssumeRole(p.assumeRoleInput)
	if err != nil {
		return credentials.Value{ProviderName: "TemporaryCredentialProvider"}, err
	}

	p.expiresAt = result.Credentials.Expiration
	return credentials.Value{
		AccessKeyID:     *result.Credentials.AccessKeyId,
		SecretAccessKey: *result.Credentials.SecretAccessKey,
		SessionToken:    *result.Credentials.SessionToken,
		ProviderName:    "TemporaryCredentialProvider",
	}, nil
}

func (p *TemporaryCredentialProvider) ExpiresAt() time.Time {
	return *p.expiresAt
}

func (p *TemporaryCredentialProvider) IsExpired() bool {
	return time.Now().After(*p.expiresAt)
}
