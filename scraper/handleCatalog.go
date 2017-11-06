package scraper

import (
    "strings"

    "github.com/anaskhan96/soup"
)

type CatalogData struct {
    Course string       //CSC1052
    Credits string      //3.00
    Description string  //Great class
}

//Extract Description and Credits
func handleCatBody(root soup.Root) []*CatalogData {
    list := make([]*CatalogData, 0)

    // keep chugging along until there's no more to parse
    for root.Error == nil {

        courseId := handleCatHead(root)
        root = root.FindNextElementSibling()
        description, credits := handleCatDetails(root)
        desc := &CatalogData{
            Course: courseId,
            Credits: credits,
            Description: description,
        }
        list = append(list, desc)
        if root.Error == nil {
            root = root.FindNextElementSibling()
        }
    }
    return list
}

//Handle catalog metadata, return course ID
func handleCatHead(root soup.Root) string {
    a := root.Find("a")
    split := strings.Split(a.Text(), " ")
    return split[0] + split[1]
}

//Handle catalog credits and descriptions - body area
func handleCatDetails(root soup.Root) (string, string) {
    if root.Error != nil {
      return "", ""
    }
    ptr := root.Find("td", "class", "ntdefault")
    ptr = ptr.Find("br")
    description := ptr.FindPrevSibling()
    ptr = ptr.FindNextElementSibling()
    credits := ptr.FindNextSibling()
    return clean(description.Pointer.Data), clean(strings.Split(credits.Pointer.Data, " ")[0])
}

func clean(str string) string {
    cleaned := strings.Replace(strings.Replace(str, "\n", " ", -1), "&nbsp;", " ", -1)
    return strings.Trim(cleaned, " ")
}
