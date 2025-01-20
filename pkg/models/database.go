package models

type Config struct {
	Agent struct {
		Authtoken string `yaml:"authtoken"`
	} `yaml:"agent"`

	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
		Host     string `yaml:"host"`
	} `yaml:"database"`

	Minio struct {
		Endpoint  string `yaml:"endpoint"`
		AccessKey string `yaml:"access_key"`
		SecretKey string `yaml:"secret_key"`
		Bucket    string `yaml:"bucketName"`
		UseSSL    bool   `yaml:"use_ssl"`
	} `yaml:"minio"`
}
