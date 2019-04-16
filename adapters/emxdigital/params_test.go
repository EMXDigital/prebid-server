package emxdigital

import (
	"encoding/json"
	"testing"

	"github.com/prebid/prebid-server/openrtb_ext"
)

var validParams = []string{
	`{"appId":"11bc5dd5-7421-4dd8-c926-40fa653bec76", "tagid": 25251, "bidfloor":0.01}`,
}

// var invalidParams = []string{
// 	`{"appId":1176, "bidfloor":0.01}`,
// 	`{"appId":"11bc5dd5-7421-4dd8-c926-40fa653bec76"}`,
// }
