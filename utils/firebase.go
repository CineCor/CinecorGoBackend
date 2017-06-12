package utils

import (
	"github.com/CineCor/CinecorGoBackend/models"
	"os"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gopkg.in/zabawaba99/firego.v1"
)

type FirebaseManager struct {
	fb *firego.Firebase
}

var firebasemutex sync.Mutex
var firebasemanager *FirebaseManager

func NewFirebaseManager() *FirebaseManager {
	manager := new(FirebaseManager)
	manager.fb = nil
	conf, err := google.JWTConfigFromJSON([]byte(os.Getenv("FIREBASE_KEY")), "https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/firebase.database")
	if err != nil {
		return manager
	}
	manager.fb = firego.New(os.Getenv("FIREBASE_DB"), conf.Client(oauth2.NoContext))
	return manager
}

func GetFirebaseManagerInstance() *FirebaseManager {
	if firebasemanager == nil {
		firebasemutex.Lock()
		defer firebasemutex.Unlock()
		if firebasemanager == nil {
			firebasemanager = NewFirebaseManager()
		}
	}
	return firebasemanager
}

func UploadCinemas(cinemas []models.Cinema) bool {
	var data models.UploadFormat
	data.Cinemas = cinemas
	data.LastUpdate = time.Now().UTC().Format(time.RFC3339)
	fbm := GetFirebaseManagerInstance()
	if err := fbm.fb.Set(data.ToMap()); err != nil {
		return false
	}
	return true
}
