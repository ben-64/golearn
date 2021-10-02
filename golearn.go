package main

import (
    "log"
	"strconv"
	"bufio"
	"time"
    "net/http"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"github.com/gorilla/mux"
	"html/template"
)

type Error struct {
	Good string
	Bad string
}

type Configuration struct {
	// Current Word in use
	CurrentWord string

	// Initial number of Words asked
	NbWords uint8

	// List of errors
	Errors []Error

	// List of words
	Words []string

	// List of words that already have been proposed
	AlreadyAskedWords []string

	// Good answser, usefull at the end
	Good uint8

	// State of the exercice
	State uint8
}

var conf Configuration

// difference returns the elements in `a` that aren't in `b`.
// https://stackoverflow.com/questions/19374219/how-to-find-the-difference-between-two-slices-of-strings
func difference(a, b []string) []string {
    mb := make(map[string]struct{}, len(b))
    for _, x := range b {
        mb[x] = struct{}{}
    }
    var diff []string
    for _, x := range a {
        if _, found := mb[x]; !found {
            diff = append(diff, x)
        }
    }
    return diff
}

// Return a simple file
func getfile(w http.ResponseWriter, r *http.Request) {
	filename := "www/" + r.URL.Path
	body, e := ioutil.ReadFile(filename)
	if e != nil {
		log.Printf("Error with file %+v : %+v",r.URL.Path,e)
	}

	// Set good content-type
	ext := filepath.Ext(filename)
	if ext == ".js" {
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	} else if ext == ".css" {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	}
	w.Write(body)
}

// Main function
func getrequest(w http.ResponseWriter, r *http.Request) {
	if(r.Method == "POST") {
		r.ParseForm()

		// Check if it is a restart
		_, need_restart := r.Form["restart"]
		if need_restart && len(conf.Words) > 0 {
			conf.AlreadyAskedWords = nil
			conf.State = 1
			conf.Errors = nil
		} else {
			file, _, has_word_file := r.FormFile("wordFile")
			if has_word_file == nil {
				// We have just receive the words list
				conf = Configuration{CurrentWord:"", NbWords:0, Errors:nil, Words:nil, AlreadyAskedWords:nil, State:1, Good:0}
				val, err := strconv.ParseUint(r.FormValue("nb_word"),10,8)
				if err != nil {
					conf.NbWords = 3
				} else {
					conf.NbWords = uint8(val)
				}
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					conf.Words = append(conf.Words, scanner.Text())
				}
			} else {
				word := r.FormValue("word")
				if len(word) != 0 {
					if word != conf.CurrentWord {
						// Mispelled word
						var e Error = Error{Good: conf.CurrentWord, Bad: word}
						conf.Errors = append(conf.Errors,e)
					}
				}
			}
		}
		if uint8(len(conf.AlreadyAskedWords)) < conf.NbWords {
			// Choose a random word, and extract it from the list to avoid asking the same word
			var remaining_words = difference(conf.Words,conf.AlreadyAskedWords)
			i := rand.Intn(len(remaining_words))
			conf.CurrentWord = remaining_words[i]
			conf.AlreadyAskedWords = append(conf.AlreadyAskedWords,conf.CurrentWord)
		} else {
			// End, display of results
			conf.State = 2
			conf.Good = conf.NbWords - uint8(len(conf.Errors))
		}
	} else {
			conf.State = 0
	}
	log.Printf("conf: %+v", conf)
	t, _ := template.ParseFiles("www/golearn.html")
	t.Execute(w, conf)
}


func main() {
	rand.Seed(time.Now().Unix())

	// Use of mux, because it handles regexp
	r := mux.NewRouter()
    r.HandleFunc("/", getrequest)
	r.HandleFunc("/{js|css}/{.*}",getfile)
	http.Handle("/",r)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
