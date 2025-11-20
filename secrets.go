package config

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

var (
	Secrets map[string]string
	once    sync.Once
)

// LoadSecrets loads from AWS Secrets Manager only once
func LoadSecrets(secretName, region string) map[string]string {
	once.Do(func() {

		// ---------------------------
		// 1. Load AWS Config
		// ---------------------------
		awsCfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
		if err != nil {
			log.Fatalf("Failed loading AWS config: %v", err)
		}

		svc := secretsmanager.NewFromConfig(awsCfg)

		input := &secretsmanager.GetSecretValueInput{
			SecretId:     aws.String(secretName),
			VersionStage: aws.String("AWSCURRENT"),
		}

		result, err := svc.GetSecretValue(context.TODO(), input)
		if err != nil {
			log.Fatalf("Error fetching secrets from AWS: %v", err)
		}

		// Parse JSON secrets
		var awsSecrets map[string]string
		if err := json.Unmarshal([]byte(*result.SecretString), &awsSecrets); err != nil {
			log.Fatalf("Failed to parse AWS secret JSON: %v", err)
		}

		Secrets = awsSecrets
	})

	return Secrets
}

// LoadSecretsEnv merges:
// 1. Environment Variables (priority) 
// 2. AWS Secrets Manager (fallback)
func LoadSecretsEnv() map[string]string {

	secretName := os.Getenv("secretName")
	region := os.Getenv("region")

	// Container/CI/CD compatibility: If these are missing,
	// we just return empty map and rely ONLY on env vars.
	if secretName != "" && region != "" {
		LoadSecrets(secretName, region)
	} else {
		Secrets = make(map[string]string)
	}

	// List of secrets expected in the application
	keys := []string{
		"MYSQL_USERNAME",
		"MYSQL_PASSWORD",
		"MYSQL_HOST",
		"MYSQL_PORT",
		"MYSQL_DATABASE",
	}

	// Final combined map
	finalSecrets := make(map[string]string)

	for _, key := range keys {
		val := os.Getenv(key)
		if val != "" {
			// Env var overrides AWS
			finalSecrets[key] = val
		} else if Secrets != nil {
			// Fallback to AWS
			if awsVal, exists := Secrets[key]; exists {
				finalSecrets[key] = awsVal
			}
		}
	}

	return finalSecrets
}
