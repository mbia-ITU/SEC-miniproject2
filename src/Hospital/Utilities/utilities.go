package utilities

import (
	"fmt"
)

func PortToString(port int) string {
	return fmt.Sprintf(":%d", port)
}
