package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/kwuyoucloud/maxbook/pkg/file"
)

// PicType is the type of pictures
var PicType = "png"

var hrefRegglobal *regexp.Regexp
var innerRegglobal *regexp.Regexp

func init() {
	if !exists("download/") {
		createFold("download/")
	}
}

func main() {
	files, _ := ioutil.ReadDir("./")
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".html") {
			htmlfilename := f.Name()
			filename := strings.Replace(htmlfilename, ".html", "", -1)
			fmt.Println("start processing file: ", htmlfilename)
			dealAHtml(htmlfilename, filename)
		}
	}

	fmt.Println("Finish. You can find the pdf files in ./download/")
}

func dealAHtml(htmlfilename string, filename string) error {
	body, err := file.ReadFile(htmlfilename)
	if err != nil {
		return err
	}

	p := make(chan struct{})
	// type 1
	hrefRegglobal = regexp.MustCompile(`<div class="webpreview-item" data-id="(\w+)" style="height:(.*)<img(.*)(src="(.*).` + PicType + `")(.*)>`)
	innerRegglobal = regexp.MustCompile(`data-id="(\w+)".*src="((.*).` + PicType + `)"`)
	fmt.Println("try type1")
	go transfertopdf(filename, 1, hrefRegglobal, innerRegglobal, "<div class=\"webpreview-item\"", "\n<div class=\"webpreview-item\"", body, p)

	<-p
	// type 2
	hrefRegglobal = regexp.MustCompile(`<div id="p(\w+)" data-page=(.*)>`)
	innerRegglobal = regexp.MustCompile(`<div id="p(\w+)".*src="(.*)" style=`)
	fmt.Println("try type2")
	go transfertopdf(filename, 2, hrefRegglobal, innerRegglobal, "<div id=\"p", "\n<div id=\"p", body, p)

	<-p

	return nil
}

func transfertopdf(filename string, modetype uint, hrefReg *regexp.Regexp, innerReg *regexp.Regexp, replaceold string, replacenew string, body []byte, p chan struct{}) error {
	defer func() {
		p <- struct{}{}
	}()

	var pdfArr []string
	bodystr := string(body)
	pagebody := strings.ReplaceAll(bodystr, replaceold, replacenew)

	hrefs := hrefReg.FindAllString(pagebody, -1)
	fmt.Println("the number of pictures : ", len(hrefs))
	if len(hrefs) <= 1 {
		return errors.New("not this type")
	}

	for _, v := range hrefs {
		// fmt.Println(v)
		temp := innerReg.FindAllStringSubmatch(v, -1)
		if len(temp) > 0 {
			pagenum := temp[0][1]
			pagelink := temp[0][2]
			newpagelink := pagelink
			// newpagelink := strings.Replace(pagelink, "./", "./", -1)
			fmt.Println("pagenum: ", pagenum, "  pagelink: ", newpagelink)

			// run copy command
			err := cpCommand(modetype, newpagelink, pagenum)
			if err != nil {
				fmt.Println(err)
			}

			// change each picture to pdf file
			picturename := ""
			newpdfName := pagenum + ".pdf"
			switch modetype {
			case 1:
				picturename = "./download/" + pagenum + "." + PicType
			case 2:
				picturename = "./download/" + pagenum
			default:
				picturename = "./download/" + pagenum + "." + PicType
			}
			onePictoPDF(picturename, "./download/"+newpdfName)
			// delete picture, store pdf file
			deleteFile(picturename)

			if exists("./download/" + newpdfName) {
				pdfArr = append(pdfArr, newpdfName)
			}
		}

		fmt.Println("********************************")
	}

	// createPDF here.
	pos, _ := os.Getwd()
	os.Chdir(pos + "/download/")
	pdffilename := filename + ".pdf"
	err := unionPDF(pdfArr, pdffilename)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Create pdf file successfully, the name is: " + pdffilename)
	}

	fmt.Println("the number of pictures is :", len(hrefs))

	return nil
}

func cpCommand(modetype uint, pagelink string, pagenum string) error {
	pagelink = strings.ReplaceAll(pagelink, "(", "\\(")
	pagelink = strings.ReplaceAll(pagelink, ")", "\\)")
	command := "cp " + pagelink + " " + "./download/" + pagenum
	switch modetype {
	case 1:
		command = "cp " + pagelink + " " + "./download/" + pagenum + "." + PicType
	default:
		command = "cp " + pagelink + " " + "./download/" + pagenum
	}
	cmd := exec.Command("bash", "-c", command)
	err := cmd.Run()
	if err != nil {
		// fmt.Println("Run this command "+command+" with an error: ", err)
		return err
	}

	return nil
}

func cpfile(file1 string, newpath string) error {
	command := "cp " + file1 + " " + newpath
	cmd := exec.Command("bash", "-c", command)
	err := cmd.Run()
	if err != nil {
		fmt.Println("cpfile with an error: ", err)
		return err
	}

	return nil
}

func onePictoPDF(sourcepicName string, destipdfName string) error {
	command := "convert " + sourcepicName + " " + destipdfName
	cmd := exec.Command("bash", "-c", command)
	err := cmd.Run()

	return err
}

func deleteFile(filename string) error {
	command := "rm " + filename
	cmd := exec.Command("bash", "-c", command)
	err := cmd.Run()

	return err
}

func unionPDF(pdfArr []string, pdfName string) error {
	pdflist := " "
	for _, v := range pdfArr {
		pdflist = pdflist + v + " "
	}

	command := "pdfunite " + pdflist + " " + pdfName
	// fmt.Println("unionpdf with command: ", command)
	cmd := exec.Command("bash", "-c", command)
	err := cmd.Run()
	if err == nil {
		// delete all files
		for _, v := range pdfArr {
			deleteFile(v)
		}

	}
	return err
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func createFold(path string) bool {
	err := os.Mkdir(path, 0755)

	if err != nil {
		return false
	}
	return true
}
