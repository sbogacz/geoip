package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/oschwald/geoip2-golang"
	"github.com/urfave/cli"
)

var (
	app                  = initApp()
	dblocation           string
	ipAddressColumnIndex int
	db                   *geoip2.Reader
	language             string
	quoteChar            string
	verbose              bool
)

func initApp() *cli.App {
	newApp := cli.NewApp()
	newApp.Name = "geoip"
	newApp.Version = "0.1.0"

	newApp.Usage = "geoip takes a line from a csv file with some column containing an IP address, and it appends columns with city, country, isAnonymousProxy, and isSatelliteProvider (in that order) from MaxMind's GeoIP2 database. If a field isn't found (such as city) then the column is left blank"
	newApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "dblocation,d",
			Usage:       "the relative location on disk of the GeoIP2 database",
			Destination: &dblocation,
			Value:       "/usr/local/var/GeoIP/GeoLite2-City.mmdb",
		},
		cli.IntFlag{
			Name:        "address-column-index,i",
			Usage:       "the index of the column containing the IP address information",
			Destination: &ipAddressColumnIndex,
		},
		cli.StringFlag{
			Name:        "language,l",
			Usage:       "the two letter code representing the language the data is desired in",
			Destination: &language,
			Value:       "en",
		},
		cli.StringFlag{
			Name:        "quote-char, c",
			Usage:       "the character used to quote the IP address, if any",
			Destination: &quoteChar,
		},
		cli.BoolFlag{
			Name:        "verbose",
			Usage:       "additional output logged",
			Destination: &verbose,
		},
	}
	newApp.Before = setup
	newApp.Action = enrich
	return newApp
}

func setup(c *cli.Context) error {
	var err error
	db, err = geoip2.Open(dblocation)
	return err
}

func enrich(c *cli.Context) error {
	scn := bufio.NewScanner(os.Stdin)
	for scn.Scan() {
		parts := strings.Split(scn.Text(), ",")
		if ipAddressColumnIndex >= len(parts) {
			return fmt.Errorf("configured columnd index %d exceeds number of columns %d", ipAddressColumnIndex, len(parts))
		}

		ipAddr := parts[ipAddressColumnIndex]
		if quoteChar != "" {
			ipAddr = strings.Trim(ipAddr, quoteChar)
		}

		geoData, err := query(ipAddr)
		if err != nil {
			return fmt.Errorf("failed to query GeoIP DB: %w", err)
		}

		data := append(parts, geoData...)
		fmt.Println(strings.Join(data, ","))
	}
	return nil
}

func query(ipAddr string) ([]string, error) {

	// If you are using strings that may be invalid, check that ip is not nil
	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return nil, fmt.Errorf("couldn't parse IP %s", ipAddr)
	}

	record, err := db.City(ip)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup ip %s: %w", ipAddr, err)
	}

	city, ok := record.City.Names[language]
	if !ok && verbose {
		fmt.Printf("no language %s entry found for City\n", language)
	}

	country, ok := record.Country.Names[language]
	if !ok && verbose {
		fmt.Printf("no language %s entry found for Country\n", language)
	}

	isAnonymousProxy := fmt.Sprintf("%t", record.Traits.IsAnonymousProxy)
	isSatelliteProvider := fmt.Sprintf("%t", record.Traits.IsSatelliteProvider)
	return []string{city, country, isAnonymousProxy, isSatelliteProvider}, nil
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("err: %v", err)
	}
	defer db.Close()
}
