package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// Config represents the application's root configuration tree.
type Config struct {
	// 1. Application Metadata & Operational Mode
	App AppConfig `mapstructure:"app"`

	// 2. Server Infrastructure & Entrypoints
	Server ServerConfig `mapstructure:"server"`
	Cookie CookieConfig `mapstructure:"cookie"`

	// 3. Security, Authentication & Access Control
	Auth     AuthConfig     `mapstructure:"auth"`
	Security SecurityConfig `mapstructure:"security"`

	// 4. Observability & Telemetry
	Logger LoggerConfig `mapstructure:"logger"`

	// 5. Primary Persistence & Storage
	DB      DBConfig      `mapstructure:"db"`
	Storage StorageConfig `mapstructure:"storage"`

	// 6. Messaging & Asynchronous Event Streaming
	Messaging MessagingConfig `mapstructure:"messaging"` // (Kafka, RabbitMQ, etc.)

	// 7. External Cloud Services & Provider Integrations
	Cloud   CloudConfig   `mapstructure:"cloud"`   // (AWS, GCP, Cloudflare, etc.)
	Email   EmailConfig   `mapstructure:"email"`   // (SendGrid, SES, Mailgun, etc.)
	Payment PaymentConfig `mapstructure:"payment"` // (Stripe, Omise, etc.)
}

// Load loads application configuration from configuration files and environment variable bindings.
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config/")

	if err := bindEnv(); err != nil {
		return nil, err
	}

	if err := viper.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return nil, fmt.Errorf("config error: %w", err)
		}
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}

func bindEnv() error {
	envBindings := map[string]string{
		"app.mode":                                  "APP_MODE",
		"logger.level":                              "LOGGER_LEVEL",
		"server.http.port":                          "SERVER_HTTP_PORT",
		"server.grpc.port":                          "SERVER_GRPC_PORT",
		"cookie.domain":                             "COOKIE_DOMAIN",
		"auth.audience":                             "AUTH_AUDIENCE",
		"auth.oauth.google.client_id":               "AUTH_OAUTH_GOOGLE_CLIENT_ID",
		"auth.oauth.google.client_secret":           "AUTH_OAUTH_GOOGLE_CLIENT_SECRET",
		"db.mysql.host":                             "DB_MYSQL_HOST",
		"db.mysql.port":                             "DB_MYSQL_PORT",
		"db.mysql.user":                             "DB_MYSQL_USER",
		"db.mysql.password":                         "DB_MYSQL_PASSWORD",
		"db.mysql.db_name":                          "DB_MYSQL_DB_NAME",
		"db.postgres.host":                          "DB_POSTGRES_HOST",
		"db.postgres.port":                          "DB_POSTGRES_PORT",
		"db.postgres.user":                          "DB_POSTGRES_USER",
		"db.postgres.password":                      "DB_POSTGRES_PASSWORD",
		"db.postgres.db_name":                       "DB_POSTGRES_DB_NAME",
		"db.mongo.url":                              "DB_MONGO_URL",
		"db.redis.addr":                             "DB_REDIS_ADDR",
		"db.redis.username":                         "DB_REDIS_USERNAME",
		"db.redis.password":                         "DB_REDIS_PASSWORD",
		"db.redis.db":                               "DB_REDIS_DB",
		"db.cassandra.hosts":                        "DB_CASSANDRA_HOSTS",
		"db.cassandra.port":                         "DB_CASSANDRA_PORT",
		"db.cassandra.username":                     "DB_CASSANDRA_USERNAME",
		"db.cassandra.password":                     "DB_CASSANDRA_PASSWORD",
		"db.cassandra.keyspace":                     "DB_CASSANDRA_KEYSPACE",
		"db.cassandra.local_dc":                     "DB_CASSANDRA_LOCAL_DC",
		"db.elasticsearch.addresses":                "DB_ELASTICSEARCH_ADDRESSES",
		"db.elasticsearch.api_key":                  "DB_ELASTICSEARCH_API_KEY",
		"storage.s3.bucket":                         "STORAGE_S3_BUCKET",
		"storage.s3.region":                         "STORAGE_S3_REGION",
		"storage.s3.endpoint":                       "STORAGE_S3_ENDPOINT",
		"storage.s3.access_key":                     "STORAGE_S3_ACCESS_KEY",
		"storage.s3.secret_key":                     "STORAGE_S3_SECRET_KEY",
		"storage.cdn.base_url":                      "STORAGE_CDN_BASE_URL",
		"messaging.kafka.cluster.brokers":           "MESSAGING_KAFKA_CLUSTER_BROKERS",
		"messaging.kafka.cluster.security_protocol": "MESSAGING_KAFKA_CLUSTER_SECURITY_PROTOCOL",
		"messaging.kafka.cluster.sasl_mechanism":    "MESSAGING_KAFKA_CLUSTER_SASL_MECHANISM",
		"messaging.kafka.cluster.username":          "MESSAGING_KAFKA_CLUSTER_USERNAME",
		"messaging.kafka.cluster.password":          "MESSAGING_KAFKA_CLUSTER_PASSWORD",
		"security.cors.allowed_origins":             "SECURITY_CORS_ALLOWED_ORIGINS",
		"security.jwt.issuer":                       "SECURITY_JWT_ISSUER",
		"security.vault.address":                    "SECURITY_VAULT_ADDRESS",
		"security.vault.token":                      "SECURITY_VAULT_TOKEN",
		"security.vault.namespace":                  "SECURITY_VAULT_NAMESPACE",
		"security.vault.role":                       "SECURITY_VAULT_ROLE",
		"payment.stripe.secret_key":                 "PAYMENT_STRIPE_SECRET_KEY",
		"payment.stripe.webhook_secret":             "PAYMENT_STRIPE_WEBHOOK_SECRET",
		"email.sendgrid_api_key":                    "EMAIL_SENDGRID_API_KEY",
	}

	for configKey, envVar := range envBindings {
		if err := viper.BindEnv(configKey, envVar); err != nil {
			return fmt.Errorf("config error: %w", err)
		}
	}

	return nil
}
