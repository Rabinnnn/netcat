package utils

import(
	"bufio"
	"log"
	"os"
)

// GetLogo is a function used to retrieve the linux logo from a file.
// It opens the specified file, scans through it line by line while
// storing the content then returns the resulting logo.
func GetLogo(filename string) string{
	var logo string

	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Error opening %q: %v\n",filename, err)
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan(){
		logo += scanner.Text() + "\n"
	}
	return logo
}