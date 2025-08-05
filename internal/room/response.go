package room


type NotificationResponse struct {
	Notifications   []*Notification `json:"notifications"`
	Unread int             `json:"unread"`
}
