package cmd

import (
		"flag"
		"fmt"
		"log"
		"os"
		"time"

		"github.com/ahermida/Goberon/scraper"
)

type CLI struct{}

// gets goberon commands
func (c *CLI) printHelp() {
    fmt.Println("goberon usage:")
    fmt.Println("  fetch - get all course data [2-4 mins]")
    fmt.Println("  index - shoves all that data into your RethinkDB instance")
		fmt.Println("  drop  - clears course data from the database")
}

// if one writes goberon alone, spits out help
func (cli *CLI) validateArgs() {
		if len(os.Args) < 2 {
				cli.printHelp()
				os.Exit(1)
		}
}


// Run parses commands and executes them
func (cli *CLI) Run() {
		cli.validateArgs()

		fetchCmd := flag.NewFlagSet("fetch", flag.ExitOnError)
		indexCmd := flag.NewFlagSet("index", flag.ExitOnError)
		dropCmd := flag.NewFlagSet("drop", flag.ExitOnError)
		loopCmd := flag.NewFlagSet("loop", flag.ExitOnError)

		switch os.Args[1] {

		case "fetch":
				err := fetchCmd.Parse(os.Args[2:])
				if err != nil {
						log.Panic(err)
				}
		case "index":
				err := indexCmd.Parse(os.Args[2:])
				if err != nil {
						log.Panic(err)
				}
		case "drop":
				err := dropCmd.Parse(os.Args[2:])
				if err != nil {
						log.Panic(err)
				}
		case "loop":
				err := loopCmd.Parse(os.Args[2:])
				if err != nil {
						log.Panic(err)
				}

		// people will type anything into a command line these days amirite
		default:
				cli.printHelp()
				os.Exit(1)
		}

		// execute fetch
		if fetchCmd.Parsed() {
				cli.fetch()
		}

		// run index
		if indexCmd.Parsed() {
				cli.index()
		}

		// run drop function
		if dropCmd.Parsed() {
				cli.drop()
		}

		// run loop function
		if loopCmd.Parsed() {
				cli.loop()
		}

}

// Gets the course schedule
func (c *CLI) fetch() {
		err := Fetch();
		if err != nil {
				fmt.Println("\n\nFor fetch to work you must set up config.Secret to have proper request data.")
				fmt.Println("You also might want to double check that your internet connection is working.")
				fmt.Println("\nERROR: ")
				log.Panic(err)
		}
		err = FetchCat();
		if err != nil {
				log.Panic(err)
		}
}

// Index courses into database and build index
func (c *CLI) index() {
		_, err := scraper.Index()
		if err != nil {
				fmt.Println("\n\nFor index to work you must first run fetch.")
				fmt.Println("You also might want to double check that your RethinDB instance is running.")
				fmt.Println("\nERROR: ")
				log.Panic(err)
		}
}

// Index courses into database and build index
func (c *CLI) drop() {
		fmt.Println("Dropping tables...")
	  err := scraper.Drop()
		if err != nil {
				fmt.Println("\n\nFor drop to work you must have a database connection.")
				fmt.Println("You might want to double check that your RethinDB instance is running.")
				fmt.Println("\nERROR: ")
				log.Panic(err)
		}
		fmt.Println("Drop complete.")
}

// Index courses into database and build index UPDATES EVERY 30 MINUTES
func (c *CLI) loop() {
		for {

				//get courses and index every hour
				c.fetch()
				c.index()
	      time.Sleep(30 * time.Minute)
	  }
}
