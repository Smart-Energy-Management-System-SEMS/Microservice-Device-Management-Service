package kafka

import (
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	"device-management-service/device-management/infrastructure/configuration"
	kafkago "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
)

type ConnectionOptions struct {
	Brokers          []string
	ClientID         string
	ConsumerGroup    string
	WriteTimeout     time.Duration
	SecurityProtocol string
	SASLMechanism    string
	Username         string
	Password         string
}

func NewConnectionOptions(config configuration.RuntimeConfig) ConnectionOptions {
	return ConnectionOptions{
		Brokers:          config.KafkaBrokers,
		ClientID:         config.KafkaClientID,
		ConsumerGroup:    config.KafkaConsumerGroup,
		WriteTimeout:     time.Duration(config.KafkaWriteTimeoutMS) * time.Millisecond,
		SecurityProtocol: config.KafkaSecurityProtocol,
		SASLMechanism:    config.KafkaSASLMechanism,
		Username:         config.KafkaUsername,
		Password:         config.KafkaPassword,
	}
}

func (o ConnectionOptions) Dialer() (*kafkago.Dialer, error) {
	mechanism, err := o.saslMechanism()
	if err != nil {
		return nil, err
	}

	dialer := &kafkago.Dialer{
		Timeout:       10 * time.Second,
		ClientID:      strings.TrimSpace(o.ClientID),
		SASLMechanism: mechanism,
	}
	if usesTLS(o.SecurityProtocol) {
		dialer.TLS = &tls.Config{MinVersion: tls.VersionTLS12}
	}
	return dialer, nil
}

func (o ConnectionOptions) Transport() (*kafkago.Transport, error) {
	mechanism, err := o.saslMechanism()
	if err != nil {
		return nil, err
	}

	transport := &kafkago.Transport{
		ClientID: strings.TrimSpace(o.ClientID),
		SASL:     mechanism,
	}
	if usesTLS(o.SecurityProtocol) {
		transport.TLS = &tls.Config{MinVersion: tls.VersionTLS12}
	}
	return transport, nil
}

func (o ConnectionOptions) saslMechanism() (sasl.Mechanism, error) {
	mechanismName := strings.ToUpper(strings.TrimSpace(o.SASLMechanism))
	username := strings.TrimSpace(o.Username)
	password := strings.TrimSpace(o.Password)

	if mechanismName == "" {
		return nil, nil
	}
	if username == "" || password == "" {
		return nil, fmt.Errorf("kafka SASL requires KAFKA_USERNAME and KAFKA_PASSWORD")
	}

	switch mechanismName {
	case "PLAIN":
		return plain.Mechanism{
			Username: username,
			Password: password,
		}, nil
	case "SCRAM-SHA-256":
		return scram.Mechanism(scram.SHA256, username, password)
	case "SCRAM-SHA-512":
		return scram.Mechanism(scram.SHA512, username, password)
	default:
		return nil, fmt.Errorf("unsupported kafka SASL mechanism %q", o.SASLMechanism)
	}
}

func usesTLS(securityProtocol string) bool {
	switch strings.ToUpper(strings.TrimSpace(securityProtocol)) {
	case "SSL", "SASL_SSL":
		return true
	default:
		return false
	}
}
