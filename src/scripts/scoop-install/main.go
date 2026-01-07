package main

import (
	"dotfiles/src/helpers/scoop"
	"fmt"
	"time"
)

func main() {
	exports, err := scoop.GetScoopExports()
	if err != nil {
		fmt.Println("Error getting scoop exports:", err)
		return
	}

	for _, app := range exports.Apps {
		fmt.Println(app.Updated.Format(time.DateTime))
	}
}
