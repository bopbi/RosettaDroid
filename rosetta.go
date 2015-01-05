package rosettaDroid

import (
    "fmt"
    "os"
    "path/filepath"
    "github.com/tealeg/xlsx"
)

func main() {
    
    if len(os.Args) > 3 {
        fmt.Println("Please input input-file and output-path")
        os.Exit(1)
    }
    
    filename := os.Args[1] // get command line first parameter
    excelFilePath, err := filepath.Abs(filename)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    fmt.Println("Opening From" + excelFilePath)
    
    xlFile, error := xlsx.OpenFile(excelFilePath)
    if error != nil {
        fmt.Println(error)
        os.Exit(1)
    }
    for _, sheet := range xlFile.Sheets {
        for _, row := range sheet.Rows {
            for _, cell := range row.Cells {
                fmt.Printf("%s\n", cell.String())
            }
        }
    }
    
}