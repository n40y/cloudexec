package engines

import (
	"context"
	"fmt"

	"cloudexec/pkg/templates"
	"cloudexec/pkg/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func ExecuteAWSEngine(tmpl *templates.Template, accessKey, secretKey, region string) error {
	ctx := context.TODO()
	var cfg aws.Config
	var err error

	if accessKey != "" && secretKey != "" {
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		)
	} else {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(region))
	}

	if err != nil {
		return fmt.Errorf("impossible de charger la configuration AWS : %w", err)
	}

	// 1. Phase systématique de reconnaissance d'identité (Whoami)
	stsClient := sts.NewFromConfig(cfg)
	_, err = stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return fmt.Errorf("échec de l'authentification AWS : %w", err)
	}

	// 2. Routage dynamique selon l'action demandée par le template YAML
	switch tmpl.Action {
	case "sts:GetCallerIdentity":
		utils.LogSuccess("MATCH TRACE [%s] -> %s", tmpl.ID, tmpl.Info.Name)

	case "s3:ListBuckets":
		utils.LogInfo("[%s] Exécution de l'énumération des buckets S3...", tmpl.ID)
		s3Client := s3.NewFromConfig(cfg)

		result, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
		if err != nil {
			return fmt.Errorf("droits insuffisants pour l'action s3:ListBuckets : %w", err)
		}

		utils.LogSuccess("!!! MATCH TROUVÉ !!! [%s] -> %s", tmpl.ID, tmpl.Info.Name)
		if len(result.Buckets) == 0 {
			fmt.Println("    └─ Aucun bucket trouvé dans ce compte.")
		}
		for _, bucket := range result.Buckets {
			fmt.Printf("    └─ "+utils.Cyan+"Bucket:"+utils.Reset+" %s (Créé le: %v)\n", *bucket.Name, bucket.CreationDate)
		}

	default:
		utils.LogWarning("Action '%s' non supportée par le moteur AWS pour le moment.", tmpl.Action)
	}

	return nil
}
