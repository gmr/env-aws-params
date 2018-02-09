package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type SSMClient struct {
	client *ssm.SSM
}

func NewSSMClient(region string) (*SSMClient, error) {
	var config *aws.Config

	awsSession := session.Must(session.NewSession(
		&aws.Config{Region: aws.String(region)}))
	_, err := awsSession.Config.Credentials.Get()
	if err != nil {
		return nil, err
	}
	config = nil

	endpoint := os.Getenv("SSM_ENDPOINT")
	if endpoint != "" {
		config = &aws.Config{
			Endpoint: &endpoint,
		}
	}

	client := ssm.New(awsSession, config)
	return &SSMClient{client}, nil
}

func (c *SSMClient) GetParametersByPath(path string) (map[string]string, error) {
	if strings.HasSuffix(path, "/") != true {
		path = fmt.Sprintf("%s/", path)
	}
	parameters := make(map[string]string)
	params := &ssm.GetParametersByPathInput{
		Path:           aws.String(path),
		Recursive:      aws.Bool(true),
		WithDecryption: aws.Bool(true),
	}
	response, err := c.client.GetParametersByPath(params)
	for _, p := range response.Parameters {
		parameters[strings.TrimPrefix(*p.Name, path)] = *p.Value
	}
	return parameters, err
}
