// Runs http://sendto/, a private secure service to exchange encrypted files and information.
package main

import (
	_ "embed"
	"log"

	"github.com/frankh/sendto"
)

func main() {
	if err := sendto.Run(); err != nil {
		log.Fatal(err)
	}
}
