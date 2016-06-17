package format

import (
	"os"
	"io"
	"bufio"
	"github.com/erikojebo/svclogman/common"
)

func FormatXml(sourcePath *string, destinationPath *string) {

	// Example formatting transformation:
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

	sourceFile, err := os.Open(*sourcePath)

	common.Check(err)

	destinationFile, err := os.Create(*destinationPath)
	
	reader := bufio.NewReader(sourceFile)
	writer := bufio.NewWriter(destinationFile)

	writer.WriteString("<svclog>\r\n")

	indentationLevel := 0
	currentContext := ""
	previousContext := ""

	var nextRune rune
	isContextSwitch := false

	for true {
		isContextSwitch = false

		rune, _, err := reader.ReadRune()

		if (err == io.EOF) {
			break
		} 

		common.Check(err)

		if (rune == '/' || rune == '!') && currentContext == "startTag" {
			
			nextRune, _, err = reader.ReadRune()

			if (err == io.EOF) {
				break
			}

			common.Check(err)

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
			common.Check(err)

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

			writeIndentation(writer, indentationLevel)
		}

		writeRune(writer, rune)
	}

	writeString(writer, "\r\n</svclog>")

	writer.Flush()
}

func writeString(writer *bufio.Writer, s string) {
	_, err := writer.WriteString(s)
	common.Check(err)
}

func writeRune(writer *bufio.Writer, rune rune) {
	_, err := writer.WriteRune(rune)
	common.Check(err)
}

func writeIndentation(writer *bufio.Writer, indentationLevel int) {

	var err error

	for i := 0; i < indentationLevel; i++ {
		_, err = writer.WriteString("\t")
		common.Check(err)
	}	
}

