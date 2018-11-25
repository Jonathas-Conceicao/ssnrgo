package ssnrgo

import ()

func TestNotification() *Notification {
	return NewAnonymousNotification(313, "Things are good")
}

func TestListingRequest() *Listing {
	return NewListingRequestAll()
}
