package main

import (
	"biblioteca/books/repository"
	"biblioteca/entities"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
)

// var bookRepo = repository.NewRelationRepository()
var bookRepo = repository.NewInmemRepository()

// GET /books
func booksIndexHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	fmt.Println("-Executing booksIndex Handler...")
	defer fmt.Println("-Completed booksIndex Handler...")

	books, err := bookRepo.FindAll()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func bookIDMiddleware(next http.Handler) http.Handler {
	bookHandler := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// *w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		// *w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		fmt.Println("--Request intercepted: bookIDMiddleware...")
		defer fmt.Println("--Request completed: bookIDMiddleware...")

		// vars := mux.Vars(req)
		// bookID, err := strconv.Atoi(vars["id"])

		bookID, err := strconv.Atoi(chi.URLParam(req, "id"))

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, err)
			return
		}

		log.Println("Getting bookID: ", bookID)
		book, err := bookRepo.FindBy(bookID)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, err)
			return
		}

		ctx := req.Context()

		ctxWithBook := context.WithValue(ctx, "currentBook", book)
		ctxWithBook = context.WithValue(ctxWithBook, "HasBook", true)

		childReq := req.WithContext(ctxWithBook)

		next.ServeHTTP(w, childReq)
	}

	return http.HandlerFunc(bookHandler)
}

// GET /books/{id}
func bookShowHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	fmt.Println("-Executing bookShow Handler...")
	defer fmt.Println("-Completed bookShow Handler...")

	book := req.Context().Value("currentBook")

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// POST /books - TODO: Check if a book with the same ISBNCode already exists!
func bookCreateHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	fmt.Println("-Executing bookCreate Handler...")
	defer fmt.Println("-Completed bookCreate Handler...")

	var newBook entities.Book
	err := json.NewDecoder(req.Body).Decode(&newBook)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprintln(w, err)                 // Unable to read the json req body
		return
	}

	createdBook, err := bookRepo.Create(newBook)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdBook)
}

// PUT /books/{id}
func bookUpdateHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	fmt.Println("-Executing bookCreate Handler...")
	defer fmt.Println("-Completed bookCreate Handler...")

	existingBook := req.Context().Value("currentBook").(entities.Book)

	err := json.NewDecoder(req.Body).Decode(&existingBook)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err)
		return
	}

	err = bookRepo.Update(&existingBook)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingBook)
}

// DELETE /books/{id}
func bookDeleteHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	fmt.Println("-Executing bookDelete Handler...")
	defer fmt.Println("-Completed bookDelete Handler...")

	existingBook := req.Context().Value("currentBook").(entities.Book)

	err := bookRepo.Destroy(&existingBook)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	// w.Header().Add("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(existingBook)
}

// Middleware
func loggingHandler(h http.Handler) http.Handler {
	handlerFn := func(w http.ResponseWriter, req *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		fmt.Println("--Request intercepted: LoggingMiddleware...")
		defer fmt.Println("--Request completed: LoggingMiddleware...")
		beginTime := time.Now()

		h.ServeHTTP(w, req)

		fmt.Println("Request:", req.Method, req.URL, "completed in", time.Since(beginTime))

	}

	return http.HandlerFunc(handlerFn)
}

func main() {
	// var port = flag.Int("port", 8000, "port number to start the server on")
	port, ok := os.LookupEnv("PORT")

	if !ok {
		port = "9000"
	}

	flag.Parse()

	// port := "9000"
	// r := http.NewServeMux()
	// r := mux.NewRouter()
	r := chi.NewRouter()
	r.Use(cors.Default().Handler)
	r.Route("/books", func(ch chi.Router) {
		ch.Get("/", booksIndexHandler)
		ch.Post("/", bookCreateHandler)

		ch.Route("/{id}", func(ch chi.Router) {
			ch.Use(bookIDMiddleware)

			ch.Get("/", bookShowHandler)
			ch.Put("/", bookUpdateHandler)
			ch.Patch("/", bookUpdateHandler)
			ch.Delete("/", bookDeleteHandler)
		})
	})

	// r.HandleFunc("/books", booksIndexHandler).Methods("GET")

	// r.Handle("/books/{id}", bookIDMiddleware(http.HandlerFunc(bookShowHandler))).Methods("GET")
	// // r.HandleFunc("/books/{id}", bookIDMiddleware(http.HandlerFunc(bookShowHandler)).ServeHTTP).Methods("GET")

	// r.HandleFunc("/books", bookCreateHandler).Methods("POST")

	// r.Handle("/books/{id}", bookIDMiddleware(http.HandlerFunc(bookUpdateHandler))).Methods("PUT", "PATCH")
	// // r.HandleFunc("/books/{id}", bookIDMiddleware(http.HandlerFunc(bookUpdateHandler))).Methods("PU.ServeHTTPT", "PATCH")

	// r.Handle("/books/{id}", bookIDMiddleware(http.HandlerFunc(bookDeleteHandler))).Methods("DELETE")
	// // r.HandleFunc("/books/{id}", bookIDMiddleware(http.HandlerFunc(bookDeleteHandler)).ServeHTTP).Methods("DELETE")

	// log.Println("Server is running on port: ", *port)
	log.Println("Server is running on port: ", port)

	// handler := cors.Default().Handler(loggingHandler(r))

	// http.ListenAndServe(":"+port, handler)

	// http.ListenAndServe(":"+strconv.Itoa(*port), loggingHandler(r))
	http.ListenAndServe(":"+port, loggingHandler(r))

	// http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, r))
}
