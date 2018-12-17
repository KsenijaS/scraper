package scraper

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/css"
	"github.com/chromedp/chromedp"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	textJS = `(function(a) {
		var s = '';
		for (var i = 0; i < a.length; i++) {
			if (a[i].offsetParent !== null) {
				s += a[i].textContent;
			}
		}
		return s;
	})($x('%s/node()'))`
	dollarSelector = "//*[contains(text(), '$')]"
)

type NodeInfo struct {
	node          *cdp.Node
	text          *string
	cssProperties *[]*css.ComputedProperty
}

func FindStyles(nodes *[]*cdp.Node, cssAttributes *[]*[]*css.ComputedProperty) chromedp.Action {
	return chromedp.ActionFunc(func(ctxt context.Context, h cdp.Executor) error {
		for _, node := range *nodes {
			attributes, err := css.GetComputedStyleForNode(node.NodeID).Do(ctxt, h)
			if err != nil {
				return err
			}
			*cssAttributes = append(*cssAttributes, &attributes)
		}
		return nil
	})
}

func FindTexts(nodes *[]*cdp.Node, strings *[]*string) chromedp.Action {
	return chromedp.ActionFunc(func(ctxt context.Context, h cdp.Executor) error {
		for _, node := range *nodes {
			var s *string
			err := chromedp.EvaluateAsDevTools(fmt.Sprintf(textJS, node.FullXPath()), &s).Do(ctxt, h)
			if err != nil {
				return err
			}
			*strings = append(*strings, s)
		}
		return nil
	})
}

func findMaxFont(infos *[]NodeInfo) string {
	var maxfont string
	for _, node := range *infos {
		attr := node.cssProperties
		for _, field := range *attr {
			if field.Name == "font-size" {
				if strings.Compare(field.Value, maxfont) > 0 {
					maxfont = field.Value
				}
			}
		}
	}
	return maxfont
}

func parseColor(color string) (int, int, int) {
	c := strings.TrimPrefix(color, "rgb(")
	c = strings.TrimSuffix(c, ")")

	colors := strings.Split(c, ",")
	colors[0] = strings.TrimSpace(colors[0])
	colors[1] = strings.TrimSpace(colors[1])
	colors[2] = strings.TrimSpace(colors[2])
	r, _ := strconv.Atoi(colors[0])
	g, _ := strconv.Atoi(colors[1])
	b, _ := strconv.Atoi(colors[2])

	return r, g, b
}

func red(r int, g int, b int) bool {
	if (g+100 < r) && (100+b < r) {
		return true
	}

	return false
}

func isRed(node NodeInfo) bool {
	attrs := *node.cssProperties
	for _, attr := range attrs {
		if attr.Name == "color" {
			r, g, b := parseColor(attr.Value)
			if red(r, g, b) {
				return true
			}
		}
	}

	return false
}

func isCrossed(node NodeInfo) bool {
	attrs := *node.cssProperties
	for _, attr := range attrs {
		if attr.Name == "text-decoration-line" {
			if attr.Value == "line-through" {
				return true
			}
		}
	}

	return false
}

func gatherNodeInfos(ctxt context.Context, selector string, url string) ([]NodeInfo, error) {
	ctxt, cancel := context.WithCancel(ctxt)
	defer cancel()

	// create chrome instance
	c, err := chromedp.New(ctxt) //, chromedp.WithLog(log.Printf))
	if err != nil {
		return nil, err
	}

	// collect node info
	var nodes []*cdp.Node
	var texts []*string
	var cssAttributes []*[]*css.ComputedProperty
	err = c.Run(ctxt, chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.Sleep(2 * time.Second),
		chromedp.Nodes(selector, &nodes),
		FindTexts(&nodes, &texts),
		FindStyles(&nodes, &cssAttributes),
	})
	if err != nil {
		return nil, fmt.Errorf("Error opening %s: %v", url, err)
	}

	// shutdown chrome
	err = c.Shutdown(ctxt)
	if err != nil {
		return nil, err
	}

	// package node info
	nodeInfos := make([]NodeInfo, len(nodes))
	for i := range nodes {
		nodeInfos[i].node = nodes[i]
		nodeInfos[i].text = texts[i]
		nodeInfos[i].cssProperties = cssAttributes[i]
	}

	// wait for chrome to finish
	err = c.Wait()
	if err != nil {
		return nil, err
	}

	return nodeInfos, nil
}

func findPrice(infos *[]NodeInfo) string {
	var nodes []NodeInfo
	var price string

	maxfont := findMaxFont(infos)
	for _, node := range *infos {
		attr := *(node.cssProperties)
		for j, _ := range attr {
			if attr[j].Name == "font-size" {
				if attr[j].Value == maxfont {
					nodes = append(nodes, node)
				}
			}
		}
	}

	for _, node := range nodes {
		if isCrossed(node) {
			log.Println("Crossed")
			continue
		}
		if isRed(node) {
			price = *node.text
			break
		}

		price = *node.text
	}

	return price
}

func ParseUrl(url string) error {
	var err error

	// create context
	ctxt := context.Background()

	infos, err := gatherNodeInfos(ctxt, dollarSelector, url)
	if err != nil {
		log.Print(err)
		return err
	}

	price := findPrice(&infos)
	log.Print(price)

	return nil
}
