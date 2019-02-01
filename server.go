package main

import (
	"net/http"
    // "fmt"
	"github.com/labstack/echo"
    "encoding/json"
    "io/ioutil"
)

func main() {
	e := echo.New()

    type Song struct{
        Title string `json:"title"`
    }

    type SongsList struct{
        Songs []Song
    }

    // Display the list of song titles
	e.GET("/", func(c echo.Context) error {
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

        // fmt.Println(responseSongs)
        // can't print this. its a struct, containing an array.

        // Printing the title to see correctness.
        // for _, songname := range responseSongs.Songs {
        //     fmt.Println(songname.Title)
        // }

        return c.JSON(http.StatusOK, responseSongs.Songs)
		// return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
