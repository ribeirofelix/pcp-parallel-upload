package main

import (
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
	"time"
	"runtime"
)

//Compile templates on start
var templates = template.Must(template.ParseFiles("upload.html"))

//Display the named template
func display(w http.ResponseWriter, tmpl string, data interface{}) {
	templates.ExecuteTemplate(w, tmpl+".html", data)
}

//This is where the action happens.
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	//GET displays the upload form.
	case "GET":
		display(w, "upload", nil)

	//POST takes the uploaded file(s) and saves it to disk.
	case "POST":
		
		// println("normal")
		// normal(w,r)

		println("parallel")
		parallel(w,r)
		
		
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func normal(w http.ResponseWriter,r *http.Request) {
	//parse the multipart form in the request
	err := r.ParseMultipartForm(100000)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get a ref to the parsed multipart form
	m := r.MultipartForm

	//get the *fileheaders
	files := m.File["myfiles"]
	bfr := time.Now()
	for i, _ := range files {
		//for each fileheader, get a handle to the actual file
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//create destination file making sure the path is writeable.
		dst, err := os.Create("test/" + files[i].Filename)
		defer dst.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//copy the uploaded file to the destination file
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
	fmt.Printf("%.6f\n", time.Now().Sub(bfr).Seconds())
	//display success message.
	display(w, "upload", "Upload successful.")
}

func parallel(w http.ResponseWriter,r *http.Request) {

	//parse the multipart form in the request
		err := r.ParseMultipartForm(100000)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
 
		//get a ref to the parsed multipart form
		m := r.MultipartForm

		//get the *fileheaders
		files := m.File["myfiles"]
		var wg sync.WaitGroup

		bfr := time.Now()
		//copy each part to destination.
		for _, f := range files {

			wg.Add(1)
			go func(fileHeader *multipart.FileHeader, nameFile string) {				
				//for each fileheader, get a handle to the actual file
				file, err := fileHeader.Open()
				defer file.Close()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				//create destination file making sure the path is writeable.
				dst, err := os.Create("test/" + nameFile)
				defer dst.Close()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				//copy the uploaded file to the destination file
				if _, err := io.Copy(dst, file); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				wg.Done()

			}(f,f.Filename)

		}

		wg.Wait()
		fmt.Printf("%.6f\n", time.Now().Sub(bfr).Seconds())
		//display success message.
		display(w, "upload", "Upload successful.")
	
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	http.HandleFunc("/upload", uploadHandler)

	//static file handler.
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	//Listen on port 8080
	http.ListenAndServe(":8080", nil)
}
