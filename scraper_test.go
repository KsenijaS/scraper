package scraper

import (
	"context"
	"path/filepath"
	"testing"
)

func TestGatherNodeInfos(t *testing.T) {
	filePath, err := filepath.Abs("./styles.html")
	if err != nil {
		t.Fatal(err)
	}
	url := "file://" + filePath

	ctxt := context.Background()

	infoh1, err := gatherNodeInfos(ctxt, "//h1", url)
	if err != nil {
		t.Fatal(err)
	}

	//Number of infos
	if len(infoh1) != 1 {
		t.Errorf("Number of infos is not correct %d\n", len(infoh1))
	}

	//Check text
	if *(infoh1[0].text) != "This is a heading" {
		t.Errorf("text for h1 is not correct %s\n", *(infoh1[0].text))
	}

	//Check cssPropertie
	for _, attr := range *(infoh1[0].cssProperties) {
		if attr.Name == "color" {
			if attr.Value != "rgb(0, 0, 255)" {
				t.Errorf("Wrong color")
			}
		}
	}

	infop, err := gatherNodeInfos(ctxt, "//p", url)
	if err != nil {
		t.Fatal(err)
	}

	//Number of infos
	if len(infop) != 1 {
		t.Errorf("Number of infos is not correct %d\n", len(infop))
	}

	//Check text
	if *(infop[0].text) != "cmok" {
		t.Errorf("text for p is not correct %s\n", *(infop[0].text))
	}

	//Check cssPropertie
	for _, attr := range *(infop[0].cssProperties) {
		if attr.Name == "color" {
			if attr.Value != "rgb(255, 0, 0)" {
				t.Errorf("Wrong color")
			}
		}
	}
}
