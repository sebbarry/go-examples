/**
This is the same file that we used to build out the link parser for our previous package.
**/


package sitemap

import (
	"io"
    "strings"

	"golang.org/x/net/html"
)

//this is what we return as a single link
//link represents a link in an HTML document. (<a href=".."> </a>)
type Link struct {
    Href string
    Text string
}


//parse will take in an html document and will return a sliice
// of links parsed from it.
/*
we will use a dfs to recurse over the nodes in the file
*/
func Parse(r io.Reader) ([]Link, error) {
    doc, err := html.Parse(r)
    if err != nil {
        return nil, err
    }

    //1. find <a> nodes in document
    //2. for each link node..
    //  2a. make a link struct.
    //  2b. make the text inside the link.
    nodes := linkNodes(doc)

    var links []Link
    for _, node := range nodes {
        //build Links that we know are link nodes.
        links = append(links, buildLink(node))
    }
    return links, nil
}


func buildLink(n *html.Node) Link {
    var ret Link
    for _, attr := range n.Attr {
        if attr.Key == "href" {
            ret.Href = attr.Val
            break
        }
    }
    ret.Text = checkText(n)
    return ret
}

//check fortext inside a node.
func checkText(n *html.Node) string {
    if n.Type == html.TextNode {        //return teh data for the node -> it is a text node.
        return n.Data
    }
    if n.Type != html.ElementNode {
        return ""
    }
    var ret string
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        ret += checkText(c) + " " //TODO lookinto a byte buffer to build a string in a more optimized way.
    }
    return strings.Join(strings.Fields(ret), " ") //TODO check out the strings package.
}

//recursive function to recurse down the DOM.
func linkNodes(n *html.Node) []*html.Node { //return all the html nodes that are links.
    //what is the base case. if the node is a link.
    if n.Type == html.ElementNode && n.Data == "a" {
        return []*html.Node{n}
    }

    var ret []*html.Node
    //recursive case.
    for c:= n.FirstChild; c != nil; c = c.NextSibling {
        ret = append(ret, linkNodes(c)...)
    }

    return ret
}



