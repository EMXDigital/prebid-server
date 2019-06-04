package emx_digital

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/mxmCherry/openrtb"
	"github.com/prebid/prebid-server/adapters"
	"github.com/prebid/prebid-server/errortypes"
	"github.com/prebid/prebid-server/openrtb_ext"
)

type EmxDigitalAdapter struct {
	endpoint string
}

func buildEndpoint(endpoint string) string {
	return endpoint + "?t=" + strconv.FormatInt(1000, 10) /* strconv.FormatInt(req.TimeoutMillis, 10) */ + "&ts=" + strconv.FormatInt(time.Now().Unix(), 10) + "&src=" + "emx_pbserver"
}

func (a *EmxDigitalAdapter) MakeRequests(request *openrtb.BidRequest) ([]*adapters.RequestData, []error) {
	var errs []error
	var adapterRequests []*adapters.RequestData

	adapterReq, errors := a.makeRequest(request)

	if adapterReq != nil {
		adapterRequests = append(adapterRequests, adapterReq)
	}
	errs = append(errs, errors...)

	return adapterRequests, errors
}

func (a *EmxDigitalAdapter) makeRequest(request *openrtb.BidRequest) (*adapters.RequestData, []error) {
	var errs []error

	if err := preprocess(request); err != nil {
		errs = append(errs, err)
	}

	fmt.Println("makeRequests")
	requestJSON, err2 := json.Marshal(request)
	if err2 != nil {
		errs = append(errs, err2)
		return nil, errs
	}

	os.Stdout.Write(requestJSON)

	reqJSON, err := json.Marshal(request)
	if err != nil {
		errs = append(errs, err)
		return nil, errs
	}

	headers := http.Header{}
	headers.Add("Content-Type", "application/json;charset=utf-8")

	// build endpoint url for rtbx
	// dont really need to put this in it's own function. could be good to error handle?
	var rtbxEndpoint = buildEndpoint(a.endpoint)

	return &adapters.RequestData{
		Method:  "POST",
		Uri:     rtbxEndpoint,
		Body:    reqJSON,
		Headers: headers,
	}, errs
}

// handle request errors and formatting to be sent to EMX
func preprocess(request *openrtb.BidRequest) error {
	for i := 0; i < len(request.Imp); i++ {
		var imp = request.Imp[i]
		var bidderExt adapters.ExtImpBidder

		if err := json.Unmarshal(imp.Ext, &bidderExt); err != nil {
			return &errortypes.BadInput{
				Message: err.Error(),
			}
		}

		var emxExt openrtb_ext.ExtImpEmxDigital

		if err := json.Unmarshal(bidderExt.Bidder, &emxExt); err != nil {
			return &errortypes.BadInput{
				Message: err.Error(),
			}
		}

		emxExtJSON, err := json.Marshal(emxExt)
		if err != nil {
			return &errortypes.BadInput{
				Message: err.Error(),
			}
		}

		request.Imp[i].Ext = emxExtJSON
		request.Imp[i].TagID = emxExt.TagID

		if request.Imp[i].BidFloor != 0 {
			request.Imp[i].BidFloor, err = strconv.ParseFloat(emxExt.BidFloor, 64)
			if err != nil {
				return &errortypes.BadInput{
					Message: err.Error(),
				}
			}
		}

		if request.Imp[i].Banner != nil {
			request.Imp[i].Banner.W = &request.Imp[i].Banner.Format[0].W
			request.Imp[i].Banner.H = &request.Imp[i].Banner.Format[0].H
		}

	}

	return nil
}

// MakeBids make the bids for the bid response.
func (a *EmxDigitalAdapter) MakeBids(internalRequest *openrtb.BidRequest, externalRequest *adapters.RequestData, response *adapters.ResponseData) (*adapters.BidderResponse, []error) {

	if response.StatusCode == http.StatusNoContent {
		// nick dev
		fmt.Println("\n204 No Content!")
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

	// nick dev
	os.Stdout.Write(response.Body)

	var bidResp openrtb.BidResponse

	if err := json.Unmarshal(response.Body, &bidResp); err != nil {
		return nil, []error{err}
	}

	bidResponse := adapters.NewBidderResponseWithBidsCapacity(1)

	for _, sb := range bidResp.SeatBid {
		for i := range sb.Bid {
			//nick dev
			// temporary fix for impid from rtbx
			sb.Bid[i].ImpID = sb.Bid[i].ID

			bidResponse.Bids = append(bidResponse.Bids, &adapters.TypedBid{
				Bid:     &sb.Bid[i],
				BidType: "banner",
			})
		}
	}

	return bidResponse, nil

}

func NewEmxDigitalBidder(endpoint string) *EmxDigitalAdapter {
	return &EmxDigitalAdapter{
		endpoint: endpoint,
	}
}
