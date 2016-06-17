package main

import (
	"fmt"
	"flag"
	"github.com/erikojebo/svclogman/format"
)

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

		format.FormatXml(sourcePath, destinationPath)
	}

	if *serve {
		if (*sourcePath == "") {
			fmt.Println("You must specify source (src) path when serving an XML log file via the web UI")
		}

		fmt.Printf("Serving a web UI on the port %v", port)
	}
}

