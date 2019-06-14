package emx_digital

import (
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

func TestEMXDigitalSyncer(t *testing.T) {
	temp := template.Must(template.New("sync-template").Parse("https://e1.emxdgt.com/hb_sync/"))
	syncer := NewEMXDigitalSyncer(temp)
	syncInfo, err := syncer.GetUsersyncInfo("", "")
	assert.NoError(t, err)
	assert.Equal(t, "https://e1.emxdgt.com/hb_sync/", syncInfo.URL)
	assert.Equal(t, "iframe", syncInfo.Type)
	// assert.EqualValues(t, 0, syncer.GDPRVendorID())
	// assert.Equal(t, false, syncInfo.SupportCORS)
}
