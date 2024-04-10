package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

// Album represents a music album
type Album struct {
	ID            string   `json:"id"`
	Name          string   `json:"name" binding:"required,min=5"`
	DateOfRelease string   `json:"date_of_release" binding:"required"`
	Genre         string   `json:"genre"`
	Price         float64  `json:"price" binding:"required,min=100,max=1000"`
	Description   string   `json:"description"`
	Musicians     []string `json:"musicians"`
}

// Musician represents a musician
type Musician struct {
	ID           string `json:"id"`
	Name         string `json:"name" binding:"required,min=3"`
	MusicianType string `json:"musician_type"`
}

var (
	albums    []Album
	musicians []Musician
)

func main() {
	router := gin.Default()

	if err := loadJSONData("music_album.json", &albums); err != nil {
		fmt.Println("Error loading music album data:", err)
		return
	}

	if err := loadJSONData("musician.json", &musicians); err != nil {
		fmt.Println("Error loading musician data:", err)
		return
	}

	router.POST("/albums", createOrUpdateAlbum)
	router.PUT("/albums/:id", createOrUpdateAlbum)
	router.POST("/musicians", createOrUpdateMusician)
	router.PUT("/musicians/:id", createOrUpdateMusician)
	router.GET("/albums", getAlbums)
	router.GET("/musicians", getMusicians)
	router.GET("/musicians/:id/albums", getAlbumsByMusician)
	router.GET("/albums/:id/musicians", getMusiciansForAlbum)

	router.Run(":8080")
}

func loadJSONData(filename string, v interface{}) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	return nil
}

func createOrUpdateAlbum(c *gin.Context) {
	var album Album
	if err := c.ShouldBindJSON(&album); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate an ID for the album
	album.ID = fmt.Sprintf("%d", len(albums)+1)

	// Check if album already exists, update it if yes, else create new
	for i, a := range albums {
		if a.ID == album.ID {
			albums[i] = album
			c.JSON(http.StatusOK, gin.H{"message": "Album updated successfully"})
			return
		}
	}

	albums = append(albums, album)
	c.JSON(http.StatusCreated, gin.H{"message": "Album created successfully"})
}

func createOrUpdateMusician(c *gin.Context) {
	var musician Musician
	if err := c.ShouldBindJSON(&musician); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate an ID for the musician
	musician.ID = fmt.Sprintf("%d", len(musicians)+1)

	// Check if musician already exists, update it if yes, else create new
	for i, m := range musicians {
		if m.ID == musician.ID {
			musicians[i] = musician
			c.JSON(http.StatusOK, gin.H{"message": "Musician updated successfully"})
			return
		}
	}

	musicians = append(musicians, musician)
	c.JSON(http.StatusCreated, gin.H{"message": "Musician created successfully"})
}

func getAlbums(c *gin.Context) {
	sort.Slice(albums, func(i, j int) bool {
		return albums[i].DateOfRelease < albums[j].DateOfRelease
	})

	c.JSON(http.StatusOK, albums)
}

func getMusicians(c *gin.Context) {
	// Read musician data from file
	musiciansData, err := ioutil.ReadFile("musician.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read musician data"})
		return
	}

	// Unmarshal musician data into slice of Musician structs
	var musicians []Musician
	if err := json.Unmarshal(musiciansData, &musicians); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse musician data"})
		return
	}

	// Return list of musicians as JSON response
	c.JSON(http.StatusOK, musicians)
}

func getAlbumsByMusician(c *gin.Context) {
	musicianID := c.Param("id")
	var musicianAlbums []Album

	for _, album := range albums {
		for _, musician := range album.Musicians {
			if musician == musicianID {
				musicianAlbums = append(musicianAlbums, album)
				break
			}
		}
	}

	sort.Slice(musicianAlbums, func(i, j int) bool {
		return musicianAlbums[i].Price < musicianAlbums[j].Price
	})

	c.JSON(http.StatusOK, musicianAlbums)
}

func getMusiciansForAlbum(c *gin.Context) {
	albumID := c.Param("id")
	var albumMusicians []Musician

	for _, album := range albums {
		if album.ID == albumID {
			for _, musicianID := range album.Musicians {
				for _, musician := range musicians {
					if musician.ID == musicianID {
						albumMusicians = append(albumMusicians, musician)
						break
					}
				}
			}
			break
		}
	}

	sort.Slice(albumMusicians, func(i, j int) bool {
		return albumMusicians[i].Name < albumMusicians[j].Name
	})

	c.JSON(http.StatusOK, albumMusicians)
}
