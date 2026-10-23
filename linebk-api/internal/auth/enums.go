package auth

// ClientType represents the type of client in the system.
type ClientType string

// ClientType constants define the different client types available in the system.
const (
	ClientTypeUser ClientType = "USER"
)

func (c ClientType) String() string {
	return string(c)
}
