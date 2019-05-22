package emx_digital

import (
    “text/template”

    “github.com/prebid/prebid-server/adapters”
    “github.com/prebid/prebid-server/usersync”
)

func New_lSyncer(temp *template.Template) usersync.Usersyncer {
    return adapters.NewSyncer(“emx_digital”, 0, temp, adapters.SyncTypeIframe)
}