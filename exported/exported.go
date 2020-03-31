package exported

import "fmt"

const IdPathVar = "id"

const NotificationPath = "/notification"
const SubscribeNotificationPath = "/subscribe"
const SpecificNotificationPath = "/notification/%s"

var SpecificNotificationRoute = fmt.Sprintf("/notification/{%s}", IdPathVar)
