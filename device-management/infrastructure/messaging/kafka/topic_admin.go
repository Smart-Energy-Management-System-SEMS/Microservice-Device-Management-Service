package kafka

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"device-management-service/device-management/infrastructure/configuration"
	kafkago "github.com/segmentio/kafka-go"
)

// EnsureTopics makes sure the service topics exist before the application
// starts publishing events. Missing topics are created once; existing topics
// are left untouched.
func EnsureTopics(ctx context.Context, options ConnectionOptions, topics configuration.KafkaTopics) error {
	broker := firstBroker(options.Brokers)
	if broker == "" {
		return fmt.Errorf("no kafka brokers configured")
	}

	dialer, err := options.Dialer()
	if err != nil {
		return err
	}

	conn, err := dialer.DialContext(ctx, "tcp", broker)
	if err != nil {
		return fmt.Errorf("dial broker %s: %w", broker, err)
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		return fmt.Errorf("read partitions: %w", err)
	}

	existingTopics := make(map[string]struct{}, len(partitions))
	for _, partition := range partitions {
		if strings.TrimSpace(partition.Topic) == "" {
			continue
		}
		existingTopics[partition.Topic] = struct{}{}
	}

	missingTopics := make([]string, 0, 6)
	for _, topic := range uniqueTopics(topics) {
		if _, exists := existingTopics[topic]; !exists {
			missingTopics = append(missingTopics, topic)
		}
	}
	if len(missingTopics) == 0 {
		log.Printf("kafka topics already available: %v", uniqueTopics(topics))
		return nil
	}

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("discover controller: %w", err)
	}

	controllerAddress := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))
	controllerConn, err := dialer.DialContext(ctx, "tcp", controllerAddress)
	if err != nil {
		return fmt.Errorf("dial controller %s: %w", controllerAddress, err)
	}
	defer controllerConn.Close()

	configs := make([]kafkago.TopicConfig, 0, len(missingTopics))
	for _, topic := range missingTopics {
		configs = append(configs, kafkago.TopicConfig{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		})
	}

	if err := controllerConn.CreateTopics(configs...); err != nil {
		return fmt.Errorf("create topics %v: %w", missingTopics, err)
	}

	log.Printf("kafka topics created: %v", missingTopics)
	return nil
}

func uniqueTopics(topics configuration.KafkaTopics) []string {
	ordered := []string{
		strings.TrimSpace(topics.DeviceRegistered),
		strings.TrimSpace(topics.DeviceStatusUpdated),
		strings.TrimSpace(topics.DeviceLinked),
		strings.TrimSpace(topics.DeviceUnlinked),
		strings.TrimSpace(topics.DeviceConfigurationUpdated),
		strings.TrimSpace(topics.DeviceEventRecorded),
	}

	seen := make(map[string]struct{}, len(ordered))
	result := make([]string, 0, len(ordered))
	for _, topic := range ordered {
		if topic == "" {
			continue
		}
		if _, exists := seen[topic]; exists {
			continue
		}
		seen[topic] = struct{}{}
		result = append(result, topic)
	}
	return result
}

func firstBroker(brokers []string) string {
	for _, broker := range brokers {
		trimmed := strings.TrimSpace(broker)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}
