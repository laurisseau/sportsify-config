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
	"github.com/joho/godotenv"
)


var (
	Secrets map[string]string
	once    sync.Once
)

func LoadSecrets(secretName string, region string) map[string]string {

	once.Do(func() {

		config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
		if err != nil {
			log.Fatal(err)
		}

		// Create Secrets Manager client
		svc := secretsmanager.NewFromConfig(config)

		input := &secretsmanager.GetSecretValueInput{
			SecretId:     aws.String(secretName),
			VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
		}

		result, err := svc.GetSecretValue(context.TODO(), input)
		if err != nil {
			// For a list of exceptions thrown, see
			// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
			log.Fatal(err.Error())
		}

		// Decrypts secret using the associated KMS key.
		//var secretString string = *result.SecretString

    	var secrets map[string]string
		if err := json.Unmarshal([]byte(*result.SecretString), &secrets); err != nil {
			log.Fatalf("Failed to parse secrets JSON: %v", err)
		}
		
		Secrets = secrets
	})

	return Secrets
}

// LoadSecretsEnv loads secretName and region from .env, then calls config.LoadSecrets
func LoadSecretsEnv() map[string]string {

	// Load .env file
	err := godotenv.Load("../config/.env")
	if err != nil {
		log.Println("Warning: No .env file found. Falling back to environment variables.")
	}

	secretName := os.Getenv("secretName")
	region := os.Getenv("region")

	return LoadSecrets(secretName, region)
}