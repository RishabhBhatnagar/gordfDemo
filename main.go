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

type test_struct struct {
	name string
}

func execute(w http.ResponseWriter, r *http.Request) {
	filename := "temp.rdf"
	fmt.Fprint(w, "hold on, your process is being sent to the server. The response will be redirected soon")
	if ioutil.WriteFile(filename, []byte(r.FormValue("data")), 777) != nil {
		fmt.Fprint(w, "Internal server error: writing to file err")
		return
	}

	err := r.ParseMultipartForm(2 * 1024 * 1024)
	if err != nil {
		fmt.Fprint(w, "failed to get input data. network issue.")
		return
	}
	fmt.Println(r.FormValue("data"))
	rdfParser := parser.New()
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
