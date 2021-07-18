package main

import (
	"fmt"
	"log"
	"os"

	"strconv"
	"time"

	"example.com/finder/models"
	"example.com/finder/services"
	"github.com/fatih/color"
	"github.com/rodaine/table"
)

var dirPath string

func main() {
	start := time.Now()
	arguments := os.Args

	if len(arguments) < 2 {
		fmt.Println("what's the source folder?")
		return
	}

	dirPath = arguments[1]

	entities, totalSize, err := services.GetEntitiesOrderedBySizeFromPath(dirPath)
	if err != nil {
		log.Println("handle error properly")
	}

	tb := initializeTableToPrintOutput("Name", "Size", "Last Modified At")
	for _, entity := range entities {
		totalSize += entity.Size
		tb.AddRow(entity.Name, entity.HumanReadableSize, entity.ModTime().Local().Format("Mon Jan 2 15:04:05 MST 2006"))
	}

	cl := color.New(color.FgRed, color.Bold)
	cl.Println("Total Size: " + models.GetHumanReadableSize(totalSize))
	cl.Println("Total files and folders: " + strconv.Itoa(len(entities)))
	tb.Print()
	elapsed := time.Since(start)
	log.Printf("Total time %s", elapsed)
}

func initializeTableToPrintOutput(columns ...interface{}) (tb table.Table) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tb = table.New(columns...)
	tb.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	return
}
