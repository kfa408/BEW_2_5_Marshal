package main

//stretch challenge from Dani: querying the API for all songs that has a bass chord.

import (
	"net/http"
    "fmt"
	"github.com/labstack/echo"
    "encoding/json"
    "io/ioutil"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	e := echo.New()

    type Song struct{
		gorm.Model
        Title string `json:"title"`
    }

    type SongsList struct{
        Songs []Song
    }

	db.AutoMigrate(&Song{})

    // Display the list of song titles
	e.GET("/", func(c echo.Context) error {
        context_ip := c.RealIP()
        fmt.Println(context_ip)
        context_req := c.Request()
        fmt.Println(context_req)
        pattern := c.QueryParam("pattern")
        if len(pattern) == 0{
            pattern = "Marley"
        }

        var url string = "https://www.songsterr.com/a/ra/songs.json?pattern=" + pattern

        resp, err := http.Get(url)
        if err != nil {
            panic(err)
        }
        // fmt.Println(resp.Body)
        // Currently resp.Body is in io.ReadCloser form
        // Thus using ioutil to make it into a byte array.
        // This is needed for json.Unmarshal

        body, _ := ioutil.ReadAll(resp.Body)

        // fmt.Println(body)
        // This is now a byte array.

        responseSongs := SongsList{}

        err1 := json.Unmarshal(body, &responseSongs.Songs)
        if err1 != nil {
            panic(err1)
        }

		for _, songname := range responseSongs.Songs {
			db.Create(&songname)
		}

        // fmt.Println(responseSongs)
        // can't print this. its a struct, containing an array.

        // Printing the title to see correctness.
        // for _, songname := range responseSongs.Songs {
        //     fmt.Println(songname.Title)
        // }

        return c.JSON(http.StatusOK, responseSongs.Songs)
		// return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/get-stored", func(c echo.Context) error {
		// var songsList SongsList
		songsList := SongsList{}
		db.Find(&songsList.Songs)
		return c.JSON(http.StatusOK, songsList.Songs)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
