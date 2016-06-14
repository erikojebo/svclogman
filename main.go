package main

import (
	"os"
	"io"
	"fmt"
	"bufio"
	"flag"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {

	sourcePath := flag.String("src", "", "Path to the source .svclog file")
	destinationPath := flag.String("dest", "", "Path to the destination .xml file")
	xml := flag.Bool("xml", true, "Create a prettified XML file from the svclog")
	serve := flag.Bool("serve", false, "Serve the log files via a web UI")
	port := flag.Int("port", 45678, "HTTP port for web UI")

	flag.Parse()

	if *xml {
		if (*sourcePath == "" || *destinationPath == "") {
			fmt.Println("You must specify both source (src) and destination (dest) paths when formatting the log file as XML")
		}

		formatXml(sourcePath, destinationPath)
	}

	if *serve {
		if (*sourcePath == "") {
			fmt.Println("You must specify source (src) path when serving an XML log file via the web UI")
		}

		fmt.Printf("Serving a web UI on the port %v", port)
	}
}

func formatXml(sourcePath *string, destinationPath *string) {

	sourceFile, err := os.Open(*sourcePath)

	check(err)

	destinationFile, err := os.Create(*destinationPath)
	
	reader := bufio.NewReader(sourceFile)
	writer := bufio.NewWriter(destinationFile)

	writer.WriteString("<svclog>\r\n")
	
	for true {

		rune, _, err := reader.ReadRune()

		if (err == io.EOF) {
			break
		} 

		check(err)
		
		if rune == '<' {
			_, err = writer.WriteString("\r\n")

			check(err)
		}

		_, err = writer.WriteRune(rune)
	}

	writer.WriteString("\r\n</svclog>")

	writer.Flush()

	check(err)	
}
