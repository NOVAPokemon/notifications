package exported

import "fmt"

const IdPathVar = "id"

const NotificationPath = "/notification"
const SubscribeNotificationPath = "/subscribe"

var SpecificNotificationPath = fmt.Sprintf("/notification/{%s}", IdPathVar)
