package gmail

import (
	"fmt"
	"time"
)

func timeAgo(t time.Time) string {
	elapsed := time.Since(t)
	if elapsed < time.Minute {
		return "just now"
	}
	if elapsed < time.Hour {
		return fmt.Sprintf("%dm ago", int(elapsed.Minutes()))
	}
	if elapsed < 24*time.Hour {
		return fmt.Sprintf("%dh ago", int(elapsed.Hours()))
	}
	return fmt.Sprintf("%dd ago", int(elapsed.Hours()/24))
}
