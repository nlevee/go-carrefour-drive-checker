package carrefour

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/goodsign/monday"
	"github.com/nlevee/go-auchan-drive-checker/pkg/drivestate"
	"github.com/nlevee/go-auchan-drive-checker/pkg/utils"
)

const (
	driveAPIURL = "https://www.carrefour.fr/api/firstslot?storeId="
)

type DriveConfig struct {
	DriveID string
	State   *drivestate.DriveState
}

type DriveStore struct {
	DriveID string
	Name    string
}

// NewConfig Create a new Drive config with driveId
func NewConfig(driveID string) DriveConfig {
	state := &drivestate.DriveState{
		IsActive: false,
		Dispo:    "",
	}
	return DriveConfig{
		DriveID: driveID,
		State:   state,
	}
}

type store struct {
	Data struct {
		Stores []struct {
			Ref  string
			Name string
		}
	}
}

// GetStoreByPostalCode fetch stores by postal code
func GetStoreByPostalCode(postalCode string) ([]DriveStore, error) {
	stores := []DriveStore{}

	cities, err := utils.GetCitiesByPostalCode(postalCode)
	if err != nil || len(cities) == 0 {
		return stores, err
	}

	url := "https://www.carrefour.fr/geoloc?modes[]=delivery&modes[]=picking&page=1&limit=5"

	city := cities[0]

	requrl := url + "&lat=" + fmt.Sprintf("%f", city.Lat) + "&lng=" + fmt.Sprintf("%f", city.Lon)
	log.Print(requrl)

	bodyContent, err := reqCarrefour(requrl)
	if err != nil {
		log.Print(err)
		return stores, err
	}

	storeFound := store{}
	json.Unmarshal(bodyContent, &storeFound)

	for _, v := range storeFound.Data.Stores {
		stores = append(stores, DriveStore{
			DriveID: v.Ref,
			Name:    v.Name,
		})
	}
	fmt.Println(stores)
	return stores, nil
}

// GetStoreIDByPostalCode fetch storeIDs by postal code
func GetStoreIDByPostalCode(postalCode string) ([]string, error) {
	storeIds := []string{}

	stores, err := GetStoreByPostalCode(postalCode)
	if err != nil {
		return storeIds, err
	}

	for _, v := range stores {
		storeIds = append(storeIds, v.DriveID)
	}

	return storeIds, nil
}

func reqCarrefour(url string) ([]byte, error) {
	bodyContent := []byte{}

	req, err := http.NewRequest("GET", url, bytes.NewReader([]byte{}))
	if err != nil {
		log.Print(err)
		return bodyContent, err
	}

	// header to fetch content as json
	req.Header.Add("x-requested-with", "XMLHttpRequest")
	req.Header.Add("accept", "application/json")

	// header to bypass datadome protextion
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36")
	req.Header.Add("accept-language", "fr-FR,fr;q=0.9")

	// exec request
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		dump, _ := httputil.DumpRequestOut(req, true)
		log.Println(err, resp.Status, string(dump))
		return bodyContent, err
	}

	defer resp.Body.Close()

	bodyContent, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return bodyContent, err
	}

	return bodyContent, nil
}

type carrefourState struct {
	Data struct {
		Attributes struct {
			BegDate string
			EndDate string
		}
	}
}

const (
	layoutCarrefour = "2006-01-02T15:04:05-0700"
	layoutDispo     = "dès Monday 2 January - 15:04"
)

func convertDate(inputDate string) string {
	t, _ := time.Parse(layoutCarrefour, inputDate)
	return monday.Format(t, layoutDispo, monday.LocaleFrFR)
}

// LoadDriveState charge le status du drive
func LoadDriveState(config DriveConfig) (hasChanged bool, err error) {
	driveURL := driveAPIURL + config.DriveID
	currentState := config.State

	log.Printf("Request uri : %v", driveURL)

	bodyContent, err := reqCarrefour(driveURL)
	if err != nil {
		log.Print(err)
		return false, err
	}

	foundState := carrefourState{}
	json.Unmarshal(bodyContent, &foundState)

	if foundState.Data.Attributes.BegDate != "" {
		newDispo := convertDate(foundState.Data.Attributes.BegDate)
		if (*currentState).Dispo != newDispo {
			(*currentState).IsActive = true
			(*currentState).Dispo = newDispo
			log.Printf("Nouveau créneau %v", currentState)
			return true, nil
		}
	} else if (*currentState).IsActive {
		log.Printf("Aucun créneau pour le moment")
		(*currentState).IsActive = false
		return true, nil
	} else {
		log.Printf("Aucun créneau pour le moment")
	}
	return false, nil
}

// LoadIntervalDriveState fetch each tick the drive state config
func LoadIntervalDriveState(config DriveConfig, tick *time.Ticker, done chan bool) {
	log.Printf("Démarrage du check de créneau Carrefour Drive %v", config.DriveID)

	// premier appel sans attendre le premier tick
	if _, err := LoadDriveState(config); err != nil {
		log.Print(err)
	}

	for {
		select {
		case <-tick.C:
			// a chaque tick du timer on lance une recherche de state
			if _, err := LoadDriveState(config); err != nil {
				log.Print(err)
			}
		case <-done:
			log.Printf("Ticker stopped")
			tick.Stop()
			return
		}
	}
}

// GetDriveState get the state of a drive
func GetDriveState(driveID string) *drivestate.DriveState {
	return drivestate.GetDriveState(driveID)
}

// NewDriveHandler add a new drive handler
func NewDriveHandler(driveID string) {
	config := NewConfig(driveID)
	drivestate.NewDriveState(driveID, config.State)

	tick := time.NewTicker(2 * time.Minute)
	done := make(chan bool)

	LoadIntervalDriveState(config, tick, done)
}
