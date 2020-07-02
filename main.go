package main

import (
	"fmt"
	"github.com/RishabhBhatnagar/gordf/rdfloader/parser"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)


func handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, nil)
}


func triplesString(triples map[string]*parser.Triple) string {
	op := ""
	i := 0
	for tripleHash := range triples {
		i++
		triple := triples[tripleHash]
		if triple.Subject.ID == "https://www.person.com/BOB" && triple.Predicate.ID == "N0" && strings.Contains(triple.Object.ID, "Alice"){
			triple = &parser.Triple{
				Subject:   &parser.Node{NodeType: parser.IRI, ID: "https://www.person.com/BOB"},
				Predicate: &parser.Node{NodeType: parser.BLANK, ID: "https://www.sample.com/namespace#likes"},
				Object:    &parser.Node{NodeType: parser.LITERAL, ID: "Alice "},
			}
		}
		op += fmt.Sprintf("Triple %v\n", i)
		fmt.Println(triple.Subject.ID, "https://www.person.com/BOB", triple.Predicate.ID == "N0", strings.Contains(triple.Object.ID, "Alice"))
		op  += fmt.Sprintf("\tSubject:   %v\n", triple.Subject)
		op  += fmt.Sprintf("\tPredicate: %v\n", triple.Predicate)
		op  += fmt.Sprintf("\tObject:    %v\n", triple.Object)

		op += fmt.Sprintf("\n")
	}
	return op
}


func execute(w http.ResponseWriter, r *http.Request) {
	filename := "temp.rdf"
	err := ioutil.WriteFile(filename, []byte(r.FormValue("data")), 777)
	defer os.Remove(filename)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	err = r.ParseMultipartForm(2 * 1024 * 1024)
	if err != nil {
		fmt.Fprint(w, "failed to get input data. network issue.")
		return
	}
	fmt.Println(r.FormValue("data"))
	rdfParser := parser.New()
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	err = rdfParser.Parse(filename)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	} else {
		err = os.Remove(filename)
		if err != nil {
			fmt.Fprint(w, "error removing temporary file")
		}
		fmt.Println(fmt.Fprintf(w, triplesString(rdfParser.Triples)))
		return
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
	log.Fatal(http.ListenAndServe(":" + port, nil))
}
