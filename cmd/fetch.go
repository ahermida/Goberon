package cmd

import (
    "fmt"
    "io"
    "net/http"
    "strings"
    "os"
    "time"

    "gopkg.in/cheggaaa/pb.v1"
    "github.com/ahermida/Goberon/config"
)

// Prints out a progress bar so you're not left hanging
func keepMeSane() {

  //seconds to complete request on avg
  i := 127
  bar := pb.New(i)
  bar.SetMaxWidth(80)
  bar.Start()
  for {
      time.Sleep(1 * time.Second)
      bar.Increment()
  }
  bar.Finish()
}

// Fetches data and writes it to file
func Fetch() error {

    client := &http.Client{}
    data := strings.NewReader(config.Secret.Data)

    // create request
    req, errCreating := http.NewRequest("POST", config.Secret.URL, data)
    if errCreating != nil {
        return errCreating
    }

    // configure it
    req.Header = config.Secret.Header
    req.AddCookie(config.Secret.Cookie)


    // Let ourselves know that we started
    fmt.Println("Fetching data")
    go keepMeSane()

    // Send request
    resp, errSending := client.Do(req);
    if errSending != nil {
        return errSending
    }

    // Get file object if there's none
    f, errFile := os.Create(config.Local.DefaultFN)
    if errFile != nil {
        return errFile
    }

    defer f.Close()

    // Write data to file
    io.Copy(f, resp.Body)

    // Give ourselves some filestats
    fi, _ := f.Stat()
    fmt.Printf("\n\nDownload Completed - %vMB long\n", fi.Size() / 1000000)
    return nil
}
