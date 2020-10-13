package dto

type PushSubscriptionDTO struct {
	Endpoint       string
	ExpirationTime uint
	Keys Keys
}

type Keys struct {
	P256dh string
	Auth string
}
