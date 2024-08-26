package helper

import (
	"fmt"

	"github.com/tubagusmf/payment-service-gb1/internal/config"
)

func GetConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.GetDbUser(),
		config.GetDbPassword(),
		config.GetDbHost(),
		config.GetDbPort(),
		config.GetDbName(),
	)

}
