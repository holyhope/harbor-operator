package v1alpha2

import (
	corev1 "k8s.io/api/core/v1"
)

const (
	// SecretTypeHTPasswd contains data needed for authenticate users against password.
	//
	// Required field:
	// - Secret.Data["auth"] - File containing user:password lines
	SecretTypeHTPasswd corev1.SecretType = "goharbor.io/htpasswd" // nolint:gosec

	// HTPasswdFileName is the file of the users required for SecretTypeHTPasswd secrets
	HTPasswdFileName = "htpasswd" // https://kubernetes.github.io/ingress-nginx/examples/auth/basic/#basic-authentication
)

const (
	// SecretTypeSharedSecret contains a single secret to share
	//
	// Required field:
	// - Secret.Data["secret"] - secret to shared
	SecretTypeSingle corev1.SecretType = "goharbor.io/single"

	// SharedSecretKey is the password required for SecretTypeSingle
	SharedSecretKey = "secret"
)

const (
	// SecretTypeRedis contains a password to connect to redis
	//
	// Required field:
	// - Secret.Data["redis-password"] - password to connect to redis, may be empty
	SecretTypeRedis corev1.SecretType = "goharbor.io/redis"

	// RedisPasswordKey is the password to connect to redis
	RedisPasswordKey = "redis-password" // https://github.com/bitnami/charts/blob/master/bitnami/redis/templates/secret.yaml#L14
)

const (
	// SecretTypeSharedSecret contains password for a postgresql user
	//
	// Required field:
	// - Secret.Data["redis-password"] - password to connect to redis, may be empty
	SecretTypePostgresql corev1.SecretType = "goharbor.io/postgresql"

	// PostgresqlPasswordKey is the password to connect to postgresql
	PostgresqlPasswordKey = "postgresql-password" // nolint:gosec
)

const (
	// SecretTypeCSRF contains data needed for CSRF security
	//
	// Required field:
	// - Secret.Data["key"] - CSRF key
	SecretTypeCSRF corev1.SecretType = "goharbor.io/csrf" // nolint:gosec

	// CSRFSecretKey is the key for SecretTypeCSRF
	CSRFSecretKey = "key"
)

const (
	// SecretTypeNotarySignerAliases contains aliases for encryption keys.
	// Only "default" key is required
	// Keys must match [a-zA-Z]([a-zA-Z0-9_]*[a-zA-Z0-9])?
	// Passwords may be any string
	//
	// Required field:
	// - Secret.Data["default"] - The default password
	SecretTypeNotarySignerAliases corev1.SecretType = "goharbor.io/notary-signer-aliases"

	// SharedSecretKey is the default password to use
	DefaultAliasSecretKey = "default"
)