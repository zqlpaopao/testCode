

package main

import (
	"fmt"
	"github.com/ThomasRooney/gexpect"
	"log"
)

func main() {
	cmd := "ssh root@11.91.161.27:22"
	pwd := "@Noah0b2"

	child, err := gexpect.Spawn(cmd)
	if err != nil {
		log.Fatal("Spawn cmd error ", err)
	}

	//if err := child.ExpectTimeout("password: ", 10*time.Second); err != nil {
	//	log.Fatal("Expect timieout error ", err)
	//}

	if err := child.SendLine(pwd); err != nil {
		log.Fatal("SendLine password error ", err)
	}

	if err := child.Wait(); err != nil {
		log.Fatal("Wait error: ", err)
	}

	fmt.Println("Success")
}