package constvar

import (
	"fmt"
)
const (
	email_accout=`"1447895999@qq.com"`
	email_authorization_code=`"slnaeglsltfohceb"`
	email_host=`"smtp.qq.com"`
	email_port="587"
)
var EmailCofig = fmt.Sprintf("{\"username\":%s,\"password\":%s,\"host\":%s,\"port\":%s}",
	email_accout,
	email_authorization_code,
	email_host,email_port,
)
