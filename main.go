package main

// fitz "github.com/gen2brain/go-fitz"
import (
	"fmt"
	"log"

	rod "github.com/go-rod/rod"
	rodLauncher "github.com/go-rod/rod/lib/launcher"

	// html2text "github.com/jaytaylor/html2text"
	md "github.com/JohannesKaufmann/html-to-markdown"

	"os"
	"strings"

	"github.com/iancoleman/strcase"
)

var filename string = "file:///Users/stephencheng/Downloads/32146-Data-Visualisation-and-Visual-Analytics---Spring-2023-2023-Aug-10_00-25-59-452/index.html"
var broswerPath string = `/Applications/Brave Browser.app/Contents/MacOS/Brave Browser`

// var filename string = "http://www.google.com/"

func renderContent(broswer *rod.Browser, pageLink string, linkTitle string, link *string, exportDir string) {
	fmt.Println("link-title: ", linkTitle)
	fmt.Println("href:", *link)

	newContentPage := broswer.MustPage(pageLink+`#!`, *link)
	htmlTitle := newContentPage.MustElement(".body .content__title").MustText()

	htmlstr := newContentPage.MustElement(".body .content__content").MustHTML()
	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(htmlstr)
	if err != nil {
		log.Fatal(err)
	}

	// open or create file
	f, err := os.OpenFile(exportDir+"/"+strcase.ToSnake(strings.ReplaceAll(htmlTitle, "/", "-"))+".md", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	// write to file

	fullConent := `---
title: ` + htmlTitle + `
---
` + `# ` + htmlTitle + `
` + markdown

	if _, err := f.WriteString(fullConent); err != nil {
		log.Fatal(err)
	}

	// relocate to img resource

}

func moduleFolder(modName string, modIndexPage string, broswer *rod.Browser) {
	// create a sub folder to ./output
	outputDir := "./output/" + strcase.ToSnake(modName)
	os.MkdirAll(outputDir, os.ModePerm)

	// create a index.md file

	page := broswer.MustPage(modIndexPage)
	htmlstr := page.MustElement("body").MustHTML()

	println(htmlstr)

	links := page.MustElements(".module-item__wrapper a")
	for _, link := range links {
		renderContent(broswer, modIndexPage, link.MustText(), link.MustAttribute("href"), outputDir)
	}
}

var moduleList [][]string = [][]string{
	{"32146-Data-Visualisation-and-Visual-Analytics", "file:///Users/stephencheng/git_src/canvas-note-to-obsidian/input_module/32146-Data-Visualisation-and-Visual-Analytics---Spring-2023-2023-Aug-10_00-25-59-452/index.html"},
	{"32541-Project-Management", "file:///Users/stephencheng/git_src/canvas-note-to-obsidian/input_module/32541-Project-Management---Spring-2023-2023-Aug-10_01-10-02-635/index.html"},
	{"32557-Enabling-Enterprise-Information-Systems", "file:///Users/stephencheng/git_src/canvas-note-to-obsidian/input_module/32557-Enabling-Enterprise-Information-Systems---Autumn-2023-2023-Aug-10_01-12-39-605/index.html"},
	{"32144-Technology-Research-Preparation", "file:///Users/stephencheng/git_src/canvas-note-to-obsidian/input_module/32144-Technology-Research-Preparation---Spring-2023-2023-Aug-10_10-12-21-665/index.html"},
}

func main() {
	lookpath, _ := rodLauncher.LookPath()
	println(lookpath)
	launcherSetup := rodLauncher.
		New().
		Bin(broswerPath).
		Headless(true).
		Devtools(true).
		MustLaunch()

	broswer := rod.New().ControlURL(launcherSetup).MustConnect()

	for _, module := range moduleList {
		moduleFolder(module[0], module[1], broswer)
	}
}
