package emxdigital

import (
	"text/template"

	"github.com/prebid/prebid-server/adapters"
	"github.com/prebid/prebid-server/usersync"
)

func NewEmxdigitalSyncer(temp *template.Template) usersync.Usersyncer {
	return adapters.NewSyncer("emxdigital", 0, temp, adapters.SyncTypeIframe)
}
