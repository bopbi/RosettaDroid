package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/tealeg/xlsx"
)

// xml structure
type String struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",innerxml"`
}

type Resources struct {
	XMLName xml.Name `xml:"resources"`
	Strings []String `xml:"string"`
}

func main() {
	toXML(os.Args)
}

func toXML(args []string) {
	pathSeparator := string(os.PathSeparator)
	var inputFilename = ""
	var outputDir = ""
	if len(args) > 1 {
		fmt.Println("Generating output on current directory")
		inputFilename = args[1]
		if len(args) == 3 {
			outputDir = args[2]
			fmt.Println("Generating output on " + outputDir)
		}
	} else {
		fmt.Println("Please input input-file and optionally output-path")
		os.Exit(1)
	}

	excelFilePath, err := filepath.Abs(inputFilename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// fmt.Println("Opening From" + excelFilePath)

	xlFile, error := xlsx.OpenFile(excelFilePath)
	if error != nil {
		fmt.Println(error)
		os.Exit(1)
	}

	if outputDir != "" {
		os.Mkdir(outputDir, 0777)
	}

	/*
		this app work by generating string xml by each language
	*/

	sheet := xlFile.Sheets[0] // only process first sheet

	var languages []string
	// get all available language on the first row
	languagesRow := sheet.Rows[0]
	for cellNumber, cell := range languagesRow.Cells {
		// skip the first cell
		if cellNumber > 0 {
			// create directory first
			var path = ""
			var cellContent, _ = cell.String()
			if outputDir == "" {
				if cellContent != "" {
					path = "values-" + cellContent
				} else {
					path = "values"
				}
			} else {
				if cellContent != "" {
					path = strings.Join([]string{outputDir, "values-" + cellContent}, pathSeparator)
				} else {
					path = strings.Join([]string{outputDir, "values"}, pathSeparator)
				}
			}

			os.Mkdir(path, 0777)

			// insert the language code into the array
			languages = append(languages, cellContent)
		} else {
			continue
		}

	}

	var stringKey []string

	// save the string key on an array
	for rowNumber, row := range sheet.Rows {
		for cellNumber, cell := range row.Cells {
			// first colomn is for available languages
			var cellContent, _ = cell.String()
			if rowNumber > 0 {
				if (cellNumber == 0) && (cellContent != "") {
					stringKey = append(stringKey, cellContent)
				} else {
					continue
				}
			}
		}
	}

	// now write the xml one by one based on the languages
	for languageIndex, language := range languages {
		fmt.Printf("Working for language [%s] ", language)
		fmt.Println("")
		xmlContent := &Resources{}
		for rowNumber, row := range sheet.Rows {
			for cellNumber, cell := range row.Cells {
				if rowNumber > 0 {
					var cellContent, _ = cell.String()
					if (cellNumber == languageIndex+1) && (cellContent != "") {
						name := stringKey[rowNumber-1]
						xmlContent.Strings = append(xmlContent.Strings, String{Name: name, Value: cellContent})
					} else {
						continue
					}
				} else {
					continue
				}

			}
		}

		outputFilename := "strings.xml"
		var langDirectory string
		if language != "" {
			langDirectory = "values-" + language
		} else {
			langDirectory = "values"
		}

		var generatedPath string
		if outputDir == "" {
			generatedPath = strings.Join([]string{langDirectory, outputFilename}, pathSeparator)
		} else {
			generatedPath = strings.Join([]string{outputDir, langDirectory, outputFilename}, pathSeparator)
		}

		file, _ := os.Create(generatedPath)

		xmlWriter := io.Writer(file)
		if _, errWrite := xmlWriter.Write([]byte(xml.Header)); errWrite != nil {
			fmt.Printf("error: %v\n", errWrite)
			break
		}
		enc := xml.NewEncoder(xmlWriter)
		enc.Indent("  ", "    ")
		if err := enc.Encode(xmlContent); err != nil {
			fmt.Printf("error: %v\n", err)
			break
		}

		fmt.Printf("the xml working for language [%s] is generated", language)
		fmt.Println("")
	}
}
