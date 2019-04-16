package openrtb_ext

// ExtImpEmxdigital defines the contract for bidrequest.imp[i].ext.emxdigital
type ExtImpEmxdigital struct {
	AppId    string  `json:"appId"`
	TagId float64 `json:"tagid"`
	BidFloor float64 `json:"bidfloor"`
}
