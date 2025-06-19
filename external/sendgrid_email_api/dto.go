package sendgrid_email_api

type ChannelType string

const (
	ChannelEmail ChannelType = "email"
)

type NotificationTarget struct {
	Type    ChannelType
	Address string
}
