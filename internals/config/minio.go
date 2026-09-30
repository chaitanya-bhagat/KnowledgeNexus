package config

const (
	MINIO_ROOT_USER     = "MINIO_ROOT_USER"
	MINIO_ROOT_PASSWORD = "MINIO_ROOT_PASSWORD"
	MINIO_ENDPOINT      = "MINIO_ENDPOINT"
	MINIO_BUCKET        = "MINIO_BUCKET"
	MINIO_USE_SSL       = "MINIO_USE_SSL"
)

type MinIOConfig struct {
	AccessKey string
	SecretKey string
	Endpoint  string
	Bucket    string
	UseSSL    bool
}

func loadMinIOConfig() (MinIOConfig, error) {
	accessKey, err := requiredEnv(MINIO_ROOT_USER)
	if err != nil {
		return MinIOConfig{}, err
	}
	secretKey, err := requiredEnv(MINIO_ROOT_PASSWORD)
	if err != nil {
		return MinIOConfig{}, err
	}
	useSSL, err := envBoolOrDefault(MINIO_USE_SSL, false)
	if err != nil {
		return MinIOConfig{}, err
	}

	return MinIOConfig{
		AccessKey: accessKey,
		SecretKey: secretKey,
		Endpoint:  envOrDefault(MINIO_ENDPOINT, "localhost:9000"),
		Bucket:    envOrDefault(MINIO_BUCKET, "knowledge-nexus"),
		UseSSL:    useSSL,
	}, nil

}
