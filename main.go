package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/RishabhBhatnagar/gordf/rdfloader/parser"
	rdfloader "github.com/RishabhBhatnagar/gordf/rdfloader/xmlreader"
	"github.com/RishabhBhatnagar/gordf/rdfwriter"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)


func handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, nil)
}


func triplesString(triples []*parser.Triple) string {
	op := ""
	i := 0
	for tripleHash := range triples {
		i++
		triple := triples[tripleHash]
		op += fmt.Sprintf("Triple %v\n", i)
		op  += fmt.Sprintf("\tSubject:   %v\n", triple.Subject)
		op  += fmt.Sprintf("\tPredicate: %v\n", triple.Predicate)
		op  += fmt.Sprintf("\tObject:    %v\n", triple.Object)
		op += fmt.Sprintf("\n")
	}
	return op
}

func xmlreaderFromString(fileContent string) rdfloader.XMLReader {
	return rdfloader.XMLReaderFromFileObject(bufio.NewReader(io.Reader(bytes.NewReader([]byte(fileContent)))))
}

func execute(w http.ResponseWriter, r *http.Request) {
	data := strings.ReplaceAll(r.FormValue("data"), "”", "\"")
	err := r.ParseMultipartForm(2 * 1024 * 1024)
	if err != nil {
		fmt.Fprint(w, "failed to get input data. network issue.")
		return
	}
	xmlReader := xmlreaderFromString(data)
	rootBlock, err := xmlReader.Read()
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	rdfParser := parser.New()
	err = rdfParser.Parse(rootBlock)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	} else {
		fmt.Println(fmt.Fprintf(w, triplesString(rdfParser.Triples)))
		return
	}
}

func execute1(w http.ResponseWriter, r *http.Request) {
	data := strings.ReplaceAll(r.FormValue("data"), "”", "\"")
	err := r.ParseMultipartForm(2 * 1024 * 1024)
	fmt.Println(err)
	if err != nil {
		fmt.Fprint(w, "failed to get input data. network issue.")
		return
	}
	xmlReader := xmlreaderFromString(data)
	rootBlock, err := xmlReader.Read()
	fmt.Println(err)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	rdfParser := parser.New()
	err = rdfParser.Parse(rootBlock)
	fmt.Println(err)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	} else {
		outputFromForm := r.FormValue("tabchars")
		tab, err := strconv.Unquote(outputFromForm)
		err = rdfwriter.WriteToFile(w, rdfParser.Triples, rdfParser.SchemaDefinition, tab)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	fmt.Println(port)
	http.HandleFunc("/", handler)
	http.HandleFunc("/get", execute)
	http.HandleFunc("/get1", execute1)
	log.Fatal(http.ListenAndServe(":" + port, nil))
}
