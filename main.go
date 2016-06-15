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

	indentationLevel := 0
	currentContext := ""
	previousContext := ""

	// <html><head>foo</head><body><form><!-- a button --><input type="submit" value="Submit" /></form></body></html>

	// <html>
	//   <head>
	//     foo
	//   </head>
	//   <body>
	//     <form>
	//       <!-- a button -->
	//       <input type="submit" value="Submit" />
	//     </form>
	//   </body>
    // </html>

	var nextRune rune
	isContextSwitch := false

	for true {
		isContextSwitch = false

		rune, _, err := reader.ReadRune()

		if (err == io.EOF) {
			break
		} 

		check(err)

		if (rune == '/' || rune == '!') && currentContext == "startTag" {
			
			nextRune, _, err = reader.ReadRune()

			if (err == io.EOF) {
				break
			}

			// Adjust when identifying that what looked like a start tag
			// was actually a self closing tag or a comment (which we can label as an end tag)
			if (rune == '/' && nextRune == '>') || (rune == '!' && nextRune == '-') {
				currentContext = "endTag"
			}

			reader.UnreadRune()
			
		} else if rune == '>' {
			previousContext = currentContext
			currentContext = ""
		} else if rune == '<' {

			nextRune, _, err = reader.ReadRune()

			if (err == io.EOF) {
				break
			}

			if currentContext != "" {
				previousContext = currentContext				
			}


			if nextRune == '/' {
				currentContext = "endTag"
			} else {
				currentContext = "startTag"
			}

			isContextSwitch = true

			reader.UnreadRune()

		} else if currentContext == "" {
			currentContext = "content"
			isContextSwitch = true
		}
		
		if isContextSwitch {

			// star+start => newline + indent++
			// start+alpha => newline + indent++
			// alpha + end => newline + indent--
			// end + end >= newline + indent--
			// end + start => newline

			_, err = writer.WriteString("\r\n")
			check(err)

			if previousContext == "startTag" && currentContext == "startTag" {
				indentationLevel += 1
			} else if previousContext == "startTag" && currentContext == "content" {
				indentationLevel += 1
			} else if previousContext == "content" && currentContext == "endTag" {
				indentationLevel -= 1
			} else if previousContext == "endTag" && currentContext == "endTag" {
				indentationLevel -= 1
			} else if previousContext == "endTag" && currentContext == "startTag" {
				// No indent change
			}

			//fmt.Printf("%v - %v\n", previousContext, currentContext)

			// Write indentation (tabs)
			_, err = writeString(writer, "", indentationLevel)
			check(err)
		}

		_, err = writer.WriteRune(rune)			
	}

	writeString(writer, "\r\n</svclog>", indentationLevel)

	writer.Flush()

	check(err)	
}

func writeString(writer *bufio.Writer, s string, indentationLevel int) (int, error) {
	writeIndentation(writer, indentationLevel)
	return writer.WriteString(s)
}

func writeIndentation(writer *bufio.Writer, indentationLevel int) {
	for i := 0; i < indentationLevel; i++ {
		writer.WriteString("\t")
	}	
}
