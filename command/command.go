package command

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/codegangsta/cli"
	"github.com/mkboudreau/sitemap/builder"
	"github.com/mkboudreau/sitemap/domain"
	"github.com/mkboudreau/sitemap/format"
)

const (
	FormatConsole string = "tab"
	FormatJson           = "json"
	FormatGraph          = "graph"
	FormatDigraph        = "digraph"
)

var (
	AllowableFormats []string = []string{FormatConsole, FormatJson, FormatDigraph, FormatGraph}
)

func AppAction() func(c *cli.Context) {
	return runApp
}

func AppCliFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   "format, f",
			Value:  "tab",
			Usage:  fmt.Sprintf("Valid values are %v", AllowableFormats),
			EnvVar: "FORMAT",
		},
		cli.IntFlag{
			Name:   "workers, w",
			Value:  10,
			Usage:  "Dictates how many workers are submitting requests to the server",
			EnvVar: "WORKERS",
		},
		cli.StringFlag{
			Name:   "rate, r",
			Value:  "2s",
			Usage:  "Rate of delay between each request for each worker. Format is time.Duration format (i.e. 1s, 2m, 4h, etc.).",
			EnvVar: "RATE",
		},
		cli.StringFlag{
			Name:   "timeout, t",
			Value:  "2500ms",
			Usage:  "Timeout for each client request. Format is time.Duration format (i.e. 1s, 2m, 4h, etc.).",
			EnvVar: "TIMEOUT",
		},
		cli.StringFlag{
			Name:   "outfile, o",
			Usage:  "Saves output in JSON file for future reformatting. (Use -i option and -f to reformat)",
			EnvVar: "OUTFILE",
		},
		cli.StringFlag{
			Name:   "infile, i",
			Usage:  "Reads from a previously saved JSON file instead of going out to a site and reprocessing. (See -o for saving results to file)",
			EnvVar: "INFILE",
		},
		cli.BoolFlag{
			Name:   "debug, d",
			Usage:  "Turn on debugging output",
			EnvVar: "DEBUG",
		},
	}
}

func runApp(c *cli.Context) {
	debug := c.Bool("debug")

	if err := validateCLI(c); err != nil {
		log.Println("CLI Validation Error:", err)
		cli.ShowAppHelp(c)
		os.Exit(1)
	}

	if !debug {
		TurnOffLogging()
	}

	var sitemap *domain.Sitemap
	if useSavedResults(c) {
		sitemap = runSitemapLoader(c)
	} else {
		sitemap = runSitemapBuilder(c)
	}

	if sitemap == nil {
		fmt.Println("Could not build Sitemap")
		return
	}

	outputSitemapToConsole(c, sitemap)

	if shouldSaveResults(c) {
		outputSitemapToFile(c, sitemap)
	}
}

func useSavedResults(c *cli.Context) bool {
	return c.String("infile") != ""
}
func shouldSaveResults(c *cli.Context) bool {
	return c.String("outfile") != ""
}

func outputSitemapToConsole(c *cli.Context, sitemap *domain.Sitemap) {
	formatter := getSitemapFormatter(c)
	formatter.Format(os.Stdout, sitemap.Top)
}

func outputSitemapToFile(c *cli.Context, sitemap *domain.Sitemap) {
	formatter := new(format.JsonSiteFormatter)
	file, err := os.Create(c.String("outfile"))
	if err != nil {
		log.Fatalf("Could not write to outfile [%v]: %v", c.String("outfile"), err)
		return
	}

	formatter.Format(file, sitemap.Top)
}

func getSitemapFormatter(c *cli.Context) format.SiteFormatter {
	f := c.String("format")

	var formatter format.SiteFormatter
	if f == FormatConsole {
		formatter = new(format.TabbedSiteFormatter)
	} else if f == FormatJson {
		formatter = new(format.JsonSiteFormatter)
	} else if f == FormatDigraph {
		formatter = new(format.DigraphSiteFormatter)
	} else if f == FormatGraph {
		formatter = new(format.GraphSiteFormatter)
	}

	return formatter
}

func runSitemapBuilder(c *cli.Context) *domain.Sitemap {
	workers := c.Int("workers")
	timeout := getTimeDurationConfig(c, "timeout")
	rate := getTimeDurationConfig(c, "rate")
	startingUrl := c.Args()[0]

	builder := builder.NewSitemapBuilder(rate, timeout, workers)

	osSignalShutdown(builder.Interrupt, 5)
	return builder.Build(startingUrl)
}

func runSitemapLoader(c *cli.Context) *domain.Sitemap {
	infile := c.String("infile")
	file, err := os.Open(infile)
	if err != nil {
		log.Fatalf("Could not open infile [%v] for reading: %v", infile, err)
		return nil
	}

	sitemap, err := loadSitemapFromFile(file)
	if err != nil {
		log.Fatalf("Could not decode json file [%v]: %v", infile, err)
		return nil
	}

	return sitemap
}

func loadSitemapFromFile(r io.Reader) (*domain.Sitemap, error) {
	var site domain.Site
	err := json.NewDecoder(r).Decode(&site)
	if err != nil {
		return nil, err
	}
	sitemap := domain.NewSitemap(&site)
	return sitemap, nil
}

func validateCLI(c *cli.Context) error {
	format := c.String("format")
	infile := c.String("infile")
	outfile := c.String("outfile")
	found := false
	for _, allowed := range AllowableFormats {
		if allowed == format {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("Did not find an allowable format to output. Found %v", format)
	}

	if infile == "" && len(c.Args()) != 1 {
		return fmt.Errorf("Must specify either an input file (-i) or a URL command-line argument")
	} else if infile != "" && len(c.Args()) == 1 {
		return fmt.Errorf("Cannot specify both an input file (-i) and a URL command-line argument")
	}

	if infile != "" && infile == outfile {
		return fmt.Errorf("infile and outfile options cannot be equal")
	}

	if infile != "" {
		if _, err := os.Open(infile); err != nil {
			return fmt.Errorf("Could not open infile [%v] for reading: %v", infile, err)
		}
	}
	if outfile != "" {
		if _, err := os.Create(outfile); err != nil {
			return fmt.Errorf("Could not create outfile [%v]: %v", outfile, err)
		}
	}

	return nil
}

func getTimeDurationConfig(c *cli.Context, key string) time.Duration {
	v := c.String(key)
	d, err := time.ParseDuration(v)
	if err != nil {
		d = 0 * time.Second
		log.Printf("Could not parse duration %v, defaulting to %v", v, d)
	}
	return d
}
