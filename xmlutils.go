// MIT License Copyright (C) 2022 Hiroshi Shimamoto
package pptxpack

import (
	"fmt"
	"regexp"
	"strings"
)

func xmlfmt(xml string) string {
	tag := regexp.MustCompile("<([/]?)([^>]+?)(/?)>")
	indent := 0
	x := tag.ReplaceAllStringFunc(xml, func(m string) string {
		curr := indent
		if !strings.HasSuffix(m, "/>") {
			if strings.HasPrefix(m, "</") {
				indent--
				curr = indent
			} else {
				indent++
			}
		}
		return "\n" + strings.Repeat("  ", curr) + m + "\n"
	})
	x = strings.Replace(x, "\n\n", "\n", -1)
	return x
}

func unpackXML(xml string) (string, error) {
	// check Microsoft Office PowerPoint
	a := strings.Split(xml, "\r\n")
	if len(a) != 2 {
		// check LibreOffice on Linux
		a = strings.Split(xml, "\n")
		if len(a) < 2 {
			return "", fmt.Errorf("bad xml")
		}
		// remove LFs and spaces between tags
		a1 := strings.Join(a[1:], "")
		a1 = regexp.MustCompile("> +<").ReplaceAllString(a1, "><")
		a[1] = a1
	}
	s := xmlfmt(a[1])
	r := strings.Replace(
		regexp.MustCompile(" +<").ReplaceAllString(s, "<"),
		"\n", "", -1)
	if r != a[1] {
		return "", fmt.Errorf("unable to revert")
	}
	return a[0] + "\n" + s, nil
}

func packXML(xml string) (string, error) {
	a := strings.SplitN(xml, "\n", 2)
	if len(a) != 2 {
		return "", fmt.Errorf("bad xml")
	}
	r := strings.Replace(
		regexp.MustCompile(" +<").ReplaceAllString(a[1], "<"),
		"\n", "", -1)
	return a[0] + "\r\n" + r, nil
}
