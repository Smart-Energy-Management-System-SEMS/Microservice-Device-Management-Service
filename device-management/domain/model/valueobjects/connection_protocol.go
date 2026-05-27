package valueobjects

import shared "device-management-service/shared/domain"

type ConnectionProtocol string

const (
	ConnectionProtocolWifi      ConnectionProtocol = "WIFI"
	ConnectionProtocolBluetooth ConnectionProtocol = "BLUETOOTH"
)

func NewConnectionProtocol(value string) (ConnectionProtocol, error) {
	protocol := ConnectionProtocol(value)
	if !protocol.IsValid() {
		return "", shared.NewValidationError("invalid connection protocol")
	}
	return protocol, nil
}

func (p ConnectionProtocol) IsValid() bool {
	return p == ConnectionProtocolWifi || p == ConnectionProtocolBluetooth
}
