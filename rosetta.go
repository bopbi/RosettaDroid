package main

import (
    "fmt"
    "os"
    "path/filepath"
    "github.com/tealeg/xlsx"
)

func main() {
    
    if len(os.Args) == 2 {
        fmt.Println("Generating output on current directory")
    } else if len(os.Args) == 3 {
        fmt.Println("Generating output on "+ os.Args[2])
    } else {
        fmt.Println("Please input input-file and optionally output-path")
        os.Exit(1)
    }
    
    filename := os.Args[1] // get command line first parameter
    excelFilePath, err := filepath.Abs(filename)
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
    for _, sheet := range xlFile.Sheets {
        for rowNumber, row := range sheet.Rows {
            for cellNumber, cell := range row.Cells {
                // first row is for available languages
                if (rowNumber == 0) {
                    
                    // skip empty cell
                    if (cellNumber > 0) {
                        fmt.Printf("%s\n", cell.String())
                    }
                    
                } else {
                    // first column is for the key
                    if (cellNumber == 0) {
                        
                    } else { // this is for the language value
                        fmt.Printf("%s\n", cell.String())
                    }
                }
                
            }
        }
    }
    
}