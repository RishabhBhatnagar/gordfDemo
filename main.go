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

//func triplesStringHelper(triples *[]*parser.Triple, lo, hi int) string {
//	if lo == hi {
//		op := ""
//		triple := (*triples)[lo]
//		op += fmt.Sprintf("Triple %v\n", lo)
//		op  += fmt.Sprintf("\tSubject:   %v\n", triple.Subject)
//		op  += fmt.Sprintf("\tPredicate: %v\n", triple.Predicate)
//		op  += fmt.Sprintf("\tObject:    %v\n", triple.Object)
//		op += fmt.Sprintf("\n")
//		return op
//	}
//	mid := lo + hi >> 1
//	return triplesStringHelper(triples, lo, mid) + triplesStringHelper(triples, mid + 1, hi)
//}
//
//func triplesString(triples []*parser.Triple) string {
//	return triplesStringHelper(&triples, 0, len(triples))
//}

func triplesString(triples []*parser.Triple) string {
	i := 0
	sb := strings.Builder{}
	for tripleHash := range triples {
		i++
		op := ""
		triple := triples[tripleHash]
		op += fmt.Sprintf("Triple %v\n", i)
		op  += fmt.Sprintf("\tSubject:   %v\n", triple.Subject)
		op  += fmt.Sprintf("\tPredicate: %v\n", triple.Predicate)
		op  += fmt.Sprintf("\tObject:    %v\n", triple.Object)
		op += fmt.Sprintf("\n")
		sb.WriteString(op)
	}
	return sb.String()
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
		outputFromForm := r.FormValue("tabchars")
		tab, err := strconv.Unquote(outputFromForm)
		err = rdfwriter.WriteToFile(w, rdfParser.Triples, rdfParser.SchemaDefinition, tab)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
	}
}

func wrapper(counter *uint64, function func (w http.ResponseWriter, r *http.Request)) func (w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {
		*counter++
		function(w, r)
		fmt.Println(*counter)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	var calledCounter uint64
	fmt.Println(port)
	http.HandleFunc("/", handler)
	http.HandleFunc("/get", wrapper(&calledCounter, execute))
	http.HandleFunc("/get1", wrapper(&calledCounter, execute1))
	log.Fatal(http.ListenAndServe(":" + port, nil))
}
