package constvar

import (
	"fmt"
)
/* const (
	DB_HOST     = "106.15.36.199"
	DB_PORT     = "3306"
	DB_USER     = "wxy"
	DB_PASSWORD = "Wxy123.."
	DB_NAME     = "day_day_fresh"
)
 */
const (
	username = "wxy"
	password = "Wxy123.."
	host     = "106.15.36.199"
	port     = "3306"
	dbName   = "day_day_fresh"

)

var MYSQL_CONN_STR = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	username,
	password,
	host,
	port,
	dbName,
)