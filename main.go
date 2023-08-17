package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type asset struct {
	// ID               string    `json:"id"`
	// DeviceAssetID    string    `json:"deviceAssetId"`
	// OwnerID          string    `json:"ownerId"`
	// DeviceID         string    `json:"deviceId"`
	// Type             string    `json:"type"`
	OriginalPath string `json:"originalPath"`
	// OriginalFileName string    `json:"originalFileName"`
	// Resized          bool      `json:"resized"`
	// Thumbhash        string    `json:"thumbhash"`
	// FileCreatedAt    time.Time `json:"fileCreatedAt"`
	// FileModifiedAt   time.Time `json:"fileModifiedAt"`
	// UpdatedAt        time.Time `json:"updatedAt"`
	// IsFavorite       bool      `json:"isFavorite"`
	// IsArchived       bool      `json:"isArchived"`
	// Duration         string    `json:"duration"`
	// ExifInfo         struct {
	// 	Make             string      `json:"make"`
	// 	Model            string      `json:"model"`
	// 	ExifImageWidth   int         `json:"exifImageWidth"`
	// 	ExifImageHeight  int         `json:"exifImageHeight"`
	// 	FileSizeInByte   int         `json:"fileSizeInByte"`
	// 	Orientation      string      `json:"orientation"`
	// 	DateTimeOriginal time.Time   `json:"dateTimeOriginal"`
	// 	ModifyDate       time.Time   `json:"modifyDate"`
	// 	TimeZone         string      `json:"timeZone"`
	// 	LensModel        string      `json:"lensModel"`
	// 	FNumber          int         `json:"fNumber"`
	// 	FocalLength      int         `json:"focalLength"`
	// 	Iso              int         `json:"iso"`
	// 	ExposureTime     string      `json:"exposureTime"`
	// 	Latitude         interface{} `json:"latitude"`
	// 	Longitude        interface{} `json:"longitude"`
	// 	City             interface{} `json:"city"`
	// 	State            interface{} `json:"state"`
	// 	Country          interface{} `json:"country"`
	// 	Description      string      `json:"description"`
	// } `json:"exifInfo,omitempty"`
	// LivePhotoVideoID interface{}   `json:"livePhotoVideoId"`
	// Tags             []interface{} `json:"tags"`
	// Checksum         string        `json:"checksum"`
}

type MyApp struct {
	Server    string
	ApiKey    string
	LocalPath string
}

func main() {
	s, err := settings()
	if err != nil {
		fmt.Println(err)
		return
	}

	server_assets, err := s.get_Server_Assets()
	if err != nil {
		log.Println(err)
	}
	//"/home/paul/Images/Atrier/1/"
	//"http://localhost:2283/api/asset"
	//"nNAFwdkBDGQyVbxl7QIa91aSDMY4UgUoOOjmr0NEtug"

	local_assets, err := s.get_Local_Assets()
	if err != nil {
		fmt.Println(err)
	}

	i := 0
	var list []string
	for local := range local_assets {
		if path.Ext(local) != ".CR3" {
			_, ok := server_assets[local]
			if !ok {
				i++
				//fmt.Println("Fichier", local, i)
				list = append(list, local)
			}
		}

	}
	var rep string

	for {
		fmt.Println("(Delete", len(list), "files ? (y or n)")
		fmt.Scanln(&rep)
		if rep == "n" {
			break
		} else if rep == "y" {
			remove_image(list)
			fmt.Println("Finnish.")
			return
		}
	}

	// for {
	// 	fmt.Println("(Delete", len(list), "files ? (y or n)")
	// 	fmt.Scanln(&rep)

	// 	if rep == "n" {
	// 		break
	// 	} else if rep == "y" {
	// 		for {
	// 			fmt.Println("Do you want move", len(list), "files to another location before deleting ? (y or n)")
	// 			fmt.Scanln(&rep)
	// 			if rep == "y" {
	// 				for i, img := range list {

	// 					final_path, err := s.move(img)
	// 					fmt.Printf("%s -> %s   %d/%d\n", img, final_path, i, len(list))
	// 					if err != nil {
	// 						fmt.Println(err)
	// 					}
	// 				}
	// 				remove_image(list)
	// 				fmt.Println("Finnish.")
	// 				return
	// 			} else if rep == "n" {
	// 				remove_image(list)
	// 				fmt.Println("Finnish.")
	// 				return
	// 			}
	// 		}
	// 	}
	// }

}

func (s *MyApp) move(img string) (string, error) {
	src, err := os.ReadFile(img)
	if err != nil {
		return "", err
	}

	dest_path := s.LocalPath + "moved/" + img[len(s.LocalPath):]
	createDirectories(dest_path)
	dst, err := os.Create(dest_path)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	os.WriteFile(dest_path, src, 0750)

	return dest_path, nil
}

func createDirectories(path string) error {
	tab_path := strings.Split(path, "/")
	//fileutil.Split
	dest := ""
	for _, directories := range tab_path[1 : len(tab_path)-1] {
		dest += "/" + directories
		_, err := os.Stat(dest)
		if err != nil {
			err := os.Mkdir(dest, 0750)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func settings() (*MyApp, error) {
	s := MyApp{}
	flag.StringVar(&s.Server, "server", "", "Server adress")
	flag.StringVar(&s.ApiKey, "key", "", "Api Key")
	flag.StringVar(&s.LocalPath, "s", "", "Local path")
	flag.Parse()
	var err error
	if s.Server == "" {
		err = errors.Join(err, errors.New("missing server"))
	}
	if s.ApiKey == "" {
		err = errors.Join(err, errors.New("missing Api key"))
	}
	if s.LocalPath == "" {
		err = errors.Join(err, errors.New("missing path to images"))
	}
	return &s, err
}

func remove_image(name []string) error {
	for _, img := range name {
		err := os.Remove(img)
		if err != nil {
			return err
		}
		fmt.Println(img, "deleted")
		name1 := strings.TrimSuffix(img, path.Ext(img)) + ".CR3"
		err = os.Remove(name1)
		if err == nil {
			fmt.Println(name1, "deleted")
		}

	}

	return nil
}

func (s *MyApp) get_Local_Assets() (map[string]any, error) {
	local_assets := map[string]any{}
	err := filepath.Walk(s.LocalPath, func(img string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if img == s.LocalPath+"moved" {
				return filepath.SkipDir
			}
			return nil
		}
		local_assets[img] = nil

		return nil
	})

	if err != nil {
		return nil, err
	}
	return local_assets, nil
}

func (s *MyApp) get_Server_Assets() (map[string]asset, error) {
	req, err := http.NewRequest("GET", s.Server, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("x-api-key", s.ApiKey)
	client := &http.Client{}
	resp, err := (client.Do(req))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return nil, errors.New(string(b))
	}

	decoder := json.NewDecoder(resp.Body)
	var r []asset
	err = decoder.Decode(&r)
	if err != nil {
		log.Println(err)
	}
	/*
		for _, path := range r {
			fmt.Println(path.OriginalPath)
		}
	*/
	server_assets := map[string]asset{}

	for _, asset := range r {
		server_assets[asset.OriginalPath] = asset
	}

	return server_assets, nil
}
