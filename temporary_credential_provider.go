package connections

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/spotmaxtech/gokit"
	"log"
	"time"
)

// Provider for role based credential
type TemporaryCredentialProvider struct {
	assumeRoleInput *sts.AssumeRoleInput
	expiresAt       *time.Time
}

// Create a temporary credential with provider
// * assumeRoleInput: sts.AssumeRoleInput, including duration etc.
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
	log.Println("credential retrieved:", gokit.Prettify(result))
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
