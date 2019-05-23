package emx_digital

import (
	"testing"

	"github.com/prebid/prebid-server/adapters/adapterstest"
)

func TestJsonSamples(t *testing.T) {
	adapterstest.RunJSONBidderTest(t, "emx_digitaltest", NewEMXDigitalBidder("https://hb.emxdgt.com"))
}
