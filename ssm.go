package main

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type SSMClient struct {
	client *ssm.Client
}

func NewSSMClient(ctx context.Context, region string, profile string) (*SSMClient, error) {
	opts := []func(*config.LoadOptions) error{}
	if profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(profile))
	}

	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, err
	}

	if region != "" {
		cfg.Region = region
	}

	return &SSMClient{client: ssm.NewFromConfig(cfg)}, nil
}

func (c *SSMClient) GetParametersByPath(ctx context.Context, path string) (map[string]string, error) {
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	var nextToken *string
	parameters := make(map[string]string)

	for {
		// MaxResults is capped at 10 by SSM for GetParametersByPath.
		response, err := c.client.GetParametersByPath(ctx, &ssm.GetParametersByPathInput{
			Path:           aws.String(path),
			Recursive:      aws.Bool(true),
			WithDecryption: aws.Bool(true),
			MaxResults:     aws.Int32(10),
			NextToken:      nextToken,
		})
		if err != nil {
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
