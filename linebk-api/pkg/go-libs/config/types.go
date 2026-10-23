package config

import "time"

// AppConfig holds application metadata configuration settings.
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Description string `mapstructure:"description"`
	Version     string `mapstructure:"version"`
	Mode        string `mapstructure:"mode"`
}

// ServerConfig holds the HTTP and gRPC server configurations.
type ServerConfig struct {
	HTTP HTTPConfig `mapstructure:"http"`
	GRPC GRPCConfig `mapstructure:"grpc"`
}

// HTTPConfig holds settings for the HTTP server.
type HTTPConfig struct {
	Port uint16 `mapstructure:"port"`
}

// GRPCConfig holds settings for the gRPC server.
type GRPCConfig struct {
	Port uint16 `mapstructure:"port"`
}

// DBConfig holds configuration settings for database connections.
//
//nolint:tagliatelle // AuthConfig fields use snake_case by design
type DBConfig struct {
	MySQL         MySQLConfig         `mapstructure:"mysql"`
	Postgres      PostgresConfig      `mapstructure:"postgres"`
	Mongo         MongoConfig         `mapstructure:"mongo"`
	Redis         RedisConfig         `mapstructure:"redis"`
	Cassandra     CassandraConfig     `mapstructure:"cassandra"`
	Elasticsearch ElasticsearchConfig `mapstructure:"elasticsearch"`
}

// MySQLConfig holds MySQL database settings.
type MySQLConfig struct {
	Host        string `mapstructure:"host"`
	Port        uint16 `mapstructure:"port"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	DbName      string `mapstructure:"db_name"`
	SSLMode     string `mapstructure:"ssl_mode"`
	SSLRootCert string `mapstructure:"ssl_root_cert"`
}

// PostgresConfig holds PostgreSQL database settings.
type PostgresConfig struct {
	Host        string `mapstructure:"host"`
	Port        uint16 `mapstructure:"port"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	DbName      string `mapstructure:"db_name"`
	SSLMode     string `mapstructure:"ssl_mode"`
	SSLRootCert string `mapstructure:"ssl_root_cert"`
}

// MongoConfig holds MongoDB connection settings.
type MongoConfig struct {
	URL     string `mapstructure:"url"`
	PoolMax uint64 `mapstructure:"pool_max"`
}

// RedisConfig holds Redis store connection settings.
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// CassandraConfig holds Apache Cassandra database settings.
type CassandraConfig struct {
	Hosts          []string      `mapstructure:"hosts"`
	Port           int           `mapstructure:"port"`
	CertPath       string        `mapstructure:"cert_path"`
	Keyspace       string        `mapstructure:"keyspace"`
	LocalDC        string        `mapstructure:"local_dc"`
	Username       string        `mapstructure:"username"`
	Password       string        `mapstructure:"password"`
	NumConns       int           `mapstructure:"num_conns"`
	ConnectTimeout time.Duration `mapstructure:"connect_timeout"`
	Timeout        time.Duration `mapstructure:"timeout"`
}

// ElasticsearchConfig holds Elasticsearch cluster settings.
type ElasticsearchConfig struct {
	Addresses          []string `mapstructure:"addresses"`
	Username           string   `mapstructure:"username"`
	Password           string   `mapstructure:"password"`
	APIKey             string   `mapstructure:"api_key"`
	CACertPath         string   `mapstructure:"ca_cert_path"`
	InsecureSkipVerify bool     `mapstructure:"insecure_skip_verify"`
}

// LoggerConfig holds logging behavior settings.
type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

// CookieConfig holds HTTP cookie options.
type CookieConfig struct {
	Domain string `mapstructure:"domain"`
}

// AuthConfig holds authentication and token settings.
//
//nolint:tagliatelle // AuthConfig fields use snake_case by design
type AuthConfig struct {
	AccessTTL  time.Duration `mapstructure:"access_ttl"`
	RefreshTTL time.Duration `mapstructure:"refresh_ttl"`
	Audience   []string      `mapstructure:"audience"`
	OAuth      OAuthConfig   `mapstructure:"oauth"`
}

// OAuthConfig holds third-party OAuth provider settings.
type OAuthConfig struct {
	Google GoogleOAuthConfig `mapstructure:"google"`
}

// GoogleOAuthConfig holds Google OAuth client options.
type GoogleOAuthConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

// StorageConfig holds object storage and CDN settings.
type StorageConfig struct {
	S3  S3Config  `mapstructure:"s3"`
	CDN CDNConfig `mapstructure:"cdn"`
}

// S3Config holds S3-compatible object storage options.
type S3Config struct {
	Bucket    string `mapstructure:"bucket"`
	Region    string `mapstructure:"region"`
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
}

// CDNConfig holds CDN distribution options.
type CDNConfig struct {
	BaseURL string `mapstructure:"base_url"`
}

// SecurityConfig holds CORS, JWT, and Vault configurations.
type SecurityConfig struct {
	CORS  CORSConfig  `mapstructure:"cors"`
	JWT   JWTConfig   `mapstructure:"jwt"`
	Vault VaultConfig `mapstructure:"vault"`
}

// CORSConfig holds Cross-Origin Resource Sharing rules.
type CORSConfig struct {
	AllowedOrigins   []string `mapstructure:"allowed_origins"`
	AllowedMethods   []string `mapstructure:"allowed_methods"`
	AllowedHeaders   []string `mapstructure:"allowed_headers"`
	ExposedHeaders   []string `mapstructure:"exposed_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	MaxAge           int      `mapstructure:"max_age"`
}

// JWTConfig holds JSON Web Token options.
type JWTConfig struct {
	Issuer string `mapstructure:"issuer"`
}

// VaultConfig holds HashiCorp Vault integration settings.
type VaultConfig struct {
	Address   string `mapstructure:"address"`
	Token     string `mapstructure:"token"`
	Namespace string `mapstructure:"namespace"`
	Role      string `mapstructure:"role"`
}

// PaymentConfig holds payment gateway configurations.
type PaymentConfig struct {
	Stripe StripeConfig `mapstructure:"stripe"`
}

// StripeConfig holds Stripe API credentials and webhook secrets.
type StripeConfig struct {
	SecretKey     string `mapstructure:"secret_key"`
	WebhookSecret string `mapstructure:"webhook_secret"`
}

// EmailConfig holds email service provider settings.
//
//nolint:tagliatelle // EmailConfig fields use snake_case by design
type EmailConfig struct {
	SendGridApiKey string `mapstructure:"sendgrid_api_key"`
	SenderEmail    string `mapstructure:"sender_email"`
	SenderName     string `mapstructure:"sender_name"`
}

// MessagingConfig holds message broker configurations.
type MessagingConfig struct {
	Kafka KafkaConfig `mapstructure:"kafka"`
}

// KafkaClusterConfig holds shared broker & authentication settings.
type KafkaClusterConfig struct {
	Brokers          string `mapstructure:"brokers"`
	SecurityProtocol string `mapstructure:"security_protocol"`
	SaslMechanism    string `mapstructure:"sasl_mechanism"`
	Username         string `mapstructure:"username"`
	Password         string `mapstructure:"password"`
}

// ConsumerConfig holds consumer-specific behaviors.
type ConsumerConfig struct {
	GroupID         string   `mapstructure:"group_id"`
	Topics          []string `mapstructure:"topics"`
	AutoOffsetReset string   `mapstructure:"auto_offset_reset"`
}

// KafkaConfig is the root wrapper for all Kafka configurations.
type KafkaConfig struct {
	Cluster  KafkaClusterConfig `mapstructure:"cluster"`
	Consumer ConsumerConfig     `mapstructure:"consumer"`
}

// CloudConfig holds cloud provider settings.
type CloudConfig struct {
	GCP GCPConfig `mapstructure:"gcp"`
}

// GCPConfig holds Google Cloud Platform configurations.
type GCPConfig struct {
	ProjectID string `mapstructure:"project_id"`
}
