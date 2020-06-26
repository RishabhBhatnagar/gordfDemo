package main

import (
	"fmt"
	"github.com/RishabhBhatnagar/gordf/rdfloader/parser"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
		op += fmt.Sprintf("Triple %v\n", i)
		op  += fmt.Sprintf("\tSubject:   %v\n", triple.Subject)
		op  += fmt.Sprintf("\tPredicate: %v\n", triple.Predicate)
		op  += fmt.Sprintf("\tObject:    %v\n", triple.Object)
	}
	return op
}

func execute(w http.ResponseWriter, r *http.Request) {
	filename := "temp.rdf"
	r.ParseForm()
	fmt.Println(ioutil.WriteFile(filename, []byte(r.FormValue("data")), 777))
	rdfParser := parser.New()
	err := rdfParser.Parse(filename)
	os.Remove(filename)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	} else {
		fmt.Println(fmt.Fprintf(w, triplesString(rdfParser.Triples)))
	}
}

func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/", handler)
	http.HandleFunc("/get", execute)
	log.Fatal(http.ListenAndServe(":" + port, nil))
}
