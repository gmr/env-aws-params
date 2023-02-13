package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	log "github.com/sirupsen/logrus"
)

type SSMClient struct {
	client *ssm.Client
}

func NewSSMClient(region string, profile string) (*SSMClient, error) {

	var cfg aws.Config
	var err error

	ctx := context.TODO()
	if profile != "" {
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithSharedConfigProfile(profile),
		)
	} else {
		cfg, err = config.LoadDefaultConfig(ctx)
	}
	if err != nil {
		return nil, err
	}

	if region != "" {
		cfg.Region = region
	}

	client := ssm.NewFromConfig(cfg)
	return &SSMClient{client}, nil
}

func (c *SSMClient) GetParametersByPath(path string) (map[string]string, error) {
	if strings.HasSuffix(path, "/") != true {
		path = fmt.Sprintf("%s/", path)
	}

	var nextToken *string
	parameters := make(map[string]string)

	for {
		params := &ssm.GetParametersByPathInput{
			Path:           aws.String(path),
			Recursive:      aws.Bool(true),
			WithDecryption: aws.Bool(true),
			MaxResults:     aws.Int32(10),
			NextToken:      nextToken,
		}
		response, err := c.client.GetParametersByPath(context.TODO(), params)

		if err != nil {
			log.Errorf("Error Getting Parameters from SSM: %s", err)
			return nil, err
		}

		for _, p := range response.Parameters {
			parameters[strings.TrimPrefix(*p.Name, path)] = *p.Value
		}

		if response.NextToken == nil {
			break
		}
		nextToken = response.NextToken
	}
	return parameters, nil
}
