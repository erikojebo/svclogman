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

	writeString(writer, "<svclog>\r\n")

	indentationLevel := 0
	currentContext := ""
	previousContext := ""
	var isContextSwitch bool

	for true {

		rune, _, err := reader.ReadRune()

		if (err == io.EOF) {
			break
		} 

		common.Check(err)

		previousContext, currentContext, isContextSwitch, err = 
			determineContext(reader, rune, previousContext, currentContext)

		if isContextSwitch {
			indentationLevel += indentationLevelDelta(previousContext, currentContext)

			writeString(writer, "\r\n")
			writeIndentation(writer, indentationLevel)
		}

		writeRune(writer, rune)
	}

	writeString(writer, "\r\n</svclog>")

	writer.Flush()
}

func determineContext(reader *bufio.Reader, rune rune, previousContext string, currentContext string) (
	updatedPreviousContext string, updatedCurrentContext string, isContextSwitch bool, err error) {

	isContextSwitch = false
	updatedCurrentContext = currentContext
	updatedPreviousContext = previousContext

	if (rune == '/' || rune == '!') && currentContext == "startTag" {
		
		nextRune, _, err := reader.ReadRune()

		if (err == io.EOF) {
			return previousContext, currentContext, false, err			
		}

		common.Check(err)

		// Adjust when identifying that what looked like a start tag
		// was actually a self closing tag or a comment (which we can label as an end tag)
		if (rune == '/' && nextRune == '>') || (rune == '!' && nextRune == '-') {
			updatedCurrentContext = "endTag"
		}

		reader.UnreadRune()
		
	} else if rune == '>' {
		updatedPreviousContext = currentContext
		updatedCurrentContext = ""
	} else if rune == '<' {

		nextRune, _, err := reader.ReadRune()

		if (err == io.EOF) {
			return previousContext, currentContext, false, err
		}

		if currentContext != "" {
			updatedPreviousContext = currentContext				
		}


		if nextRune == '/' {
			updatedCurrentContext = "endTag"
		} else {
			updatedCurrentContext = "startTag"
		}

		isContextSwitch = true

		reader.UnreadRune()

	} else if currentContext == "" {
		updatedCurrentContext = "content"
		isContextSwitch = true
	}

	return updatedPreviousContext, updatedCurrentContext, isContextSwitch, err
}

func indentationLevelDelta(previousContext string, currentContext string) (indentationLevelDelta int) {
	// star+start => newline + indent++
	// start+alpha => newline + indent++
	// alpha + end => newline + indent--
	// end + end >= newline + indent--
	// end + start => newline

	if previousContext == "startTag" && currentContext == "startTag" {
		indentationLevelDelta = 1
	} else if previousContext == "startTag" && currentContext == "content" {
		indentationLevelDelta = 1
	} else if previousContext == "content" && currentContext == "endTag" {
		indentationLevelDelta = -1
	} else if previousContext == "endTag" && currentContext == "endTag" {
		indentationLevelDelta = -1
	} else if previousContext == "endTag" && currentContext == "startTag" {
		// No indent change
	}

	return indentationLevelDelta
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

