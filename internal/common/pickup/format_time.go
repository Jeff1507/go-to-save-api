package pickup

import "fmt"

func FormatTime(hour, minute int) string {
	return fmt.Sprintf("%02d:%02d", hour, minute)
}
