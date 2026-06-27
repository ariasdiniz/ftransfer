package menu

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ariasdiniz/ftransfer/src/transfer"
)

func ShowMenu(metadata transfer.Metadata) transfer.Metadata {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Will you send or receive a file?")
	fmt.Println("Press 1 for Send and any other key for Send. Then press ENTER")
	scanner.Scan()

	isSender := scanner.Text()
	if isSender == "1" {
		metadata.Conn = "sender"

	FileRetry:
		fmt.Println("Which file will you send? Write the relative path to this folder")
		scanner.Scan()
		file := scanner.Text()

		_, err := os.Stat(file)
		if err != nil {
			fmt.Printf("File %s does not exist\n", file)
			goto FileRetry
		}
		metadata.Fname = file
		return metadata
	}

	fmt.Println("What is the IP of the file sender? Write in this format: 0.0.0.0")
	scanner.Scan()
	host := scanner.Text()

	if host != "" {
		metadata.Host = host
	}
	return metadata
}
