package emx_digital

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	// dev
	"fmt"
	"os"

	"github.com/mxmCherry/openrtb"
	"github.com/prebid/prebid-server/adapters"
	"github.com/prebid/prebid-server/errortypes"
)

// --- globals ---
// var Endpoint = "https://hb.emxdgt.com"

// ---- type definitions ----
type EmxDigitalAdapter struct {
	endpoint string
}

type responseImp struct {
	ID     string `json:"id"`
	Price  string `json:"price"`
	AdM    string `json:"adm"`
	CrID   string `json:"crid"`
	Width  uint64 `json:"w,omitempty"`
	Height uint64 `json:"h,omitempty"`
	Secure uint64 `json:"secure"`
	tid    uint64 `json:"secure"`
}

// --- functionality ---
// unpack pbs requests into http requests to fetch bids
// , req *pbs.PBSRequest
func (a *EmxDigitalAdapter) MakeRequests(request *openrtb.BidRequest) ([]*adapters.RequestData, []error) {
	now := time.Now()
	// create endpoint
	a.endpoint += "?t=" + strconv.FormatInt(1000, 10) /* strconv.FormatInt(req.TimeoutMillis, 10) */ + "&ts=" + strconv.FormatInt(now.Unix(), 10)

	var errors = make([]error, 0)

	reqJSON, err := json.Marshal(request)
	if err != nil {
		errors = append(errors, err)
		return nil, errors
	}

	headers := http.Header{}
	headers.Add("Content-Type", "application/json;charset=utf-8")

	os.Stdout.Write(reqJSON)
	fmt.Printf("%T\n", request)

	return []*adapters.RequestData{{
		Method:  "POST",
		Uri:     a.endpoint,
		Body:    reqJSON,
		Headers: headers,
	}}, errors

}

// unpack server responses into bids.
func (a *EmxDigitalAdapter) MakeBids(internalRequest *openrtb.BidRequest, externalRequest *adapters.RequestData, response *adapters.ResponseData) (*adapters.BidderResponse, []error) {

	os.Stdout.Write(response.Body)
	// fmt.Printf("%T\n", response)

	if response.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	if response.StatusCode == http.StatusBadRequest {
		return nil, []error{&errortypes.BadInput{
			Message: fmt.Sprintf("Unexpected status code: %d. Run with request.debug = 1 for more info", response.StatusCode),
		}}
	}

	if response.StatusCode != http.StatusOK {
		return nil, []error{&errortypes.BadServerResponse{
			Message: fmt.Sprintf("Unexpected status code: %d. Run with request.debug = 1 for more info", response.StatusCode),
		}}
	}

	var bidResp openrtb.BidResponse
	if err := json.Unmarshal(response.Body, &bidResp); err != nil {
		return nil, []error{err}
	}

	bidResponse := adapters.NewBidderResponseWithBidsCapacity(1)

	for _, sb := range bidResp.SeatBid {
		for i := range sb.Bid {
			bidResponse.Bids = append(bidResponse.Bids, &adapters.TypedBid{
				Bid:     &sb.Bid[i],
				BidType: "banner",
			})
		}
	}
	return bidResponse, nil

}

// --- helpers ---

// dev only?
func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in
	}
	return out.String()
}

// logging for errors
func BadInput(msg string) *errortypes.BadInput {
	return &errortypes.BadInput{
		Message: msg,
	}
}

// logging for errors
func BadServerResponse(msg string) *errortypes.BadServerResponse {
	return &errortypes.BadServerResponse{
		Message: msg,
	}
}

// instantiate new Adapter for PBS
func NewEmxDigitalBidder(endpoint string) *EmxDigitalAdapter {
	return &EmxDigitalAdapter{endpoint}
}
