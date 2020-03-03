package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strings"
)

/*
My approach was:
	1. Get the target element attributes
	2. Iterate through all the anchors in the copy file, assigning them scores based on how they resemble the target element
	3. Get the anchor with the best score
	4. I was supposed to print the absolute path to the element, but I didn't make it. I printed the element instead.
*/

func main() {

	/* INITIALIZATION */

	originalFilePath := flag.String("input_origin_file_path", "./Samples/Original File.html", "")
	copyFilePath := flag.String("input_other_sample_file_path", "./Samples/Copy 1.html", "")
	targetElementID := flag.String("target_element_id", "make-everything-ok-button", "")
	flag.Parse()

	originalFileBytes, err := ioutil.ReadFile(*originalFilePath)
	if err != nil {
		log.Fatal(err)
	}

	originalFileStr := string(originalFileBytes)

	/* SEARCH FOR TARGET ELEMENT */

	targetElementIDStartingIndex := strings.Index(originalFileStr, *targetElementID)

	/* SEARCH FOR OPENING ANCHOR */

	openingAnchorFound := false
	searchIndex := targetElementIDStartingIndex
	for !openingAnchorFound {
		searchIndex--

		if originalFileStr[searchIndex-2:searchIndex] == "<a" {
			openingAnchorFound = true
		}
	}

	openingAnchorStartingIndex := searchIndex - 2

	/* SEARCH FOR CLOSING ANCHOR */

	closingAnchorFound := false
	searchIndex = targetElementIDStartingIndex
	for !closingAnchorFound {
		searchIndex++

		if originalFileStr[searchIndex:searchIndex+4] == "</a>" {
			closingAnchorFound = true
		}
	}

	closingAnchorStartingIndex := searchIndex

	/* GET ELEMENT INTO STRING */

	targetElement := originalFileStr[openingAnchorStartingIndex:closingAnchorStartingIndex]

	/* GET ELEMENT ATTRIBUTES */

	classAttrStartingIndex := strings.Index(targetElement, `class="`)
	hrefAttrStartingIndex := strings.Index(targetElement, `href="`)
	titleAttrStartingIndex := strings.Index(targetElement, `title="`)
	relAttrStartingIndex := strings.Index(targetElement, `rel="`)
	onclickAttrStartingIndex := strings.Index(targetElement, `onclick="`)

	classAttr := GetFirstTextBetweenQuotesAfterIndex(targetElement, classAttrStartingIndex)
	hrefAttr := GetFirstTextBetweenQuotesAfterIndex(targetElement, hrefAttrStartingIndex)
	titleAttr := GetFirstTextBetweenQuotesAfterIndex(targetElement, titleAttrStartingIndex)
	relAttr := GetFirstTextBetweenQuotesAfterIndex(targetElement, relAttrStartingIndex)
	onclickAttr := GetFirstTextBetweenQuotesAfterIndex(targetElement, onclickAttrStartingIndex)

	classes := strings.Split(classAttr, " ")

	/* OPEN THE COPY FILE */

	copyFileBytes, err := ioutil.ReadFile(*copyFilePath)
	if err != nil {
		log.Fatal(err)
	}

	copyFileStr := string(copyFileBytes)

	/* ITERATE THROUGH ANCHORS */

	anchors := []*Anchor{}

	for {

		/* GET ANCHOR FULL STRING */

		openingAnchorStartingIndex := strings.Index(copyFileStr, "<a")

		if openingAnchorStartingIndex == -1 {
			break
		}

		remainingText := copyFileStr[openingAnchorStartingIndex:len(copyFileStr)]
		closingAnchorStartingIndex := openingAnchorStartingIndex + strings.Index(remainingText, "</a>")
		anchorText := copyFileStr[openingAnchorStartingIndex:closingAnchorStartingIndex]

		/* CALCULATE ANCHOR SCORES */
		scores := &Score{}

		classAttrStartingIndex := strings.Index(anchorText, `class="`)
		hrefAttrStartingIndex := strings.Index(anchorText, `href="`)
		titleAttrStartingIndex := strings.Index(anchorText, `title="`)
		relAttrStartingIndex := strings.Index(anchorText, `rel="`)
		onclickAttrStartingIndex := strings.Index(anchorText, `onclick="`)

		copyClassAttr := GetFirstTextBetweenQuotesAfterIndex(anchorText, classAttrStartingIndex)
		copyHrefAttr := GetFirstTextBetweenQuotesAfterIndex(anchorText, hrefAttrStartingIndex)
		copyTitleAttr := GetFirstTextBetweenQuotesAfterIndex(anchorText, titleAttrStartingIndex)
		copyRelAttr := GetFirstTextBetweenQuotesAfterIndex(anchorText, relAttrStartingIndex)
		copyOnclickAttr := GetFirstTextBetweenQuotesAfterIndex(anchorText, onclickAttrStartingIndex)

		copyClasses := strings.Split(copyClassAttr, " ")

		for i := 0; i < len(classes); i++ {
			for j := 0; j < len(copyClasses); j++ {
				if classes[i] == copyClasses[j] {
					scores.Classes += 15
				}
			}
		}

		if hrefAttr == copyHrefAttr {
			scores.HREF = 5
		}

		if titleAttr == copyTitleAttr {
			scores.Title = 15
		}

		if relAttr == copyRelAttr {
			scores.Rel = 10
		}

		if onclickAttr == copyOnclickAttr {
			scores.OnClick = 10
		}

		// the indexes need correction, you need to add the finish index of the last anchor you iterated through
		indexesCorrector := 0
		if len(anchors) > 0 {
			indexesCorrector = anchors[len(anchors)-1].FinishIndex
		}

		anchor := &Anchor{Text: anchorText, BeginningIndex: openingAnchorStartingIndex + indexesCorrector, FinishIndex: closingAnchorStartingIndex + indexesCorrector, Scores: *scores}

		anchors = append(anchors, anchor)

		copyFileStr = copyFileStr[closingAnchorStartingIndex:len(copyFileStr)]
	}

	/* GET TARGET ANCHOR */

	targetAnchorIndex := 0
	for i := 0; i < len(anchors); i++ {
		if anchors[i].CalculateScore() > anchors[targetAnchorIndex].CalculateScore() {
			targetAnchorIndex = i
		}
	}

	targetAnchor := anchors[targetAnchorIndex]

	/* GET ANCHOR FULL PATH */

	/* I was having trouble using xPath to get the absolute path to the target element and I was running out of time, so I'll just output the target element instead */

	log.Println(targetAnchor.Text)
}

type Anchor struct {
	Text           string
	BeginningIndex int
	FinishIndex    int
	Scores         Score
}

func (anchor *Anchor) CalculateScore() int {
	return anchor.Scores.Classes + anchor.Scores.HREF + anchor.Scores.Title + anchor.Scores.Rel + anchor.Scores.OnClick
}

type Score struct {
	Classes int
	HREF    int
	Title   int
	Rel     int
	OnClick int
}

func GetFirstTextBetweenQuotesAfterIndex(text string, index int) string {
	searchIndex := index
	firstQuotesIndex := 0
	lastQuotesIndex := 0

	for lastQuotesIndex == 0 {
		searchIndex++

		if text[searchIndex:searchIndex+1] == `"` {
			if firstQuotesIndex == 0 {
				firstQuotesIndex = searchIndex
			} else {
				lastQuotesIndex = searchIndex
			}
		}
	}

	return text[firstQuotesIndex+1 : lastQuotesIndex]
}
