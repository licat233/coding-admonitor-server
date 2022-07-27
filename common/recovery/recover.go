package recovery

import (
	"fmt"
)

func RecoverHandler() {
	msg := recover()
	if msg != nil {
		fmt.Println(msg)
	}
}
