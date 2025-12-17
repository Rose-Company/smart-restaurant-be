package digital_ocean_storage

import (
	"app-noti/config"
	awss3 "app-noti/pkg/awsS3"
)

func NewDOStorage(prefix string) (error, *awss3.S3Storage) {
	cfg := awss3.StorageConfigureParams{
		AccessKey: config.Config.DigitalOcean.StorageAccessKey,
		SecretKey: config.Config.DigitalOcean.StorageSecretKey,
		Endpoint:  config.Config.DigitalOcean.StorageEndPoint,
		Region:    config.Config.DigitalOcean.StorageRegion,
	}

	storage := awss3.S3Storage{}
	err := storage.Configure(prefix, cfg)
	if err != nil {
		return err, nil
	}

	return nil, &storage
}
