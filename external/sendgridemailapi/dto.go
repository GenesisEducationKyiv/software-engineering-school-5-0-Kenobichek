package sendgridemailapi

type ChannelType string

const (
	ChannelEmail ChannelType = "email"
)

type NotificationTarget struct {
	Type    ChannelType
	Address string
}
