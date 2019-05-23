package nick

import (
	"text/template"

	"github.com/prebid/prebid-server/adapters"
	"github.com/prebid/prebid-server/usersync"
)

func NewNickSyncer(temp *template.Template) usersync.Usersyncer {
	return adapters.NewSyncer("nick", 36, temp, adapters.SyncTypeRedirect)
}
