package main

import (
	"fmt"
	Utils "github.com/CineCor/CinecorGoBackend/utils"
)

func main() {
	fmt.Println("Initializing Firebase...")
	Utils.GetFirebaseManagerInstance()

	fmt.Println("Parsing Data...")
	err, cinemas := Utils.GetCinemas()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Filling Movie Data...")
	Utils.FillMovieData(&cinemas)

	fmt.Println("Writing to Firebase...")
	Utils.UploadCinemas(cinemas)
}
