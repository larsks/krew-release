package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	flag "github.com/spf13/pflag"
)

type (
	Asset struct {
		Name   string
		Url    string
		Sha256 string
		Arch   string
		Os     string
	}

	Options struct {
		Release  string
		Template string
		Version  bool
		Help     bool
	}
)

var options Options

func must(err error, msg string, v ...any) {
	if err != nil {
		newmsg := fmt.Sprintf(msg, v...)
		log.Fatalf("%s: %s", newmsg, err)
	}
}

func (a *Asset) ExtractArchOs() {
	stage1 := strings.Split(a.Name, ".")
	stage2 := strings.Split(stage1[0], "-")
	a.Arch = stage2[len(stage2)-1]
	a.Os = stage2[len(stage2)-2]
}

func (a *Asset) DownloadAndChecksum() error {
	log.Printf("fetching %s", a.Url)
	resp, err := http.Get(a.Url)
	must(err, "failed to download asset %s", a.Name)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	hash := sha256.New()
	if _, err := io.Copy(hash, resp.Body); err != nil {
		return fmt.Errorf("failed to calculate checksum for %s: %v", a.Name, err)
	}

	cksum := hash.Sum(nil)
	a.Sha256 = hex.EncodeToString(cksum)
	return nil
}

func getAssets(release string) ([]*Asset, error) {
	cmd := exec.Command("gh", "release", "view", release, "--json", "assets", "--jq", "[.assets[]|{name: .name, url: .url}]")
	var out strings.Builder
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	var assets []*Asset
	must(json.Unmarshal([]byte(out.String()), &assets), "failed to unmarshal")

	for _, asset := range assets {
		asset.ExtractArchOs()
	}
	return assets, nil
}

func init() {
	flag.BoolVarP(&options.Version, "version", "v", false, "")
	flag.BoolVarP(&options.Help, "help", "h", false, "")
	flag.Usage = usage
}

func usage() {
	prg := os.Args[0]
	fmt.Fprintf(flag.CommandLine.Output(), "%s: usage: %s [options] release template\n\nOptions:\n\n", prg, prg)
	flag.CommandLine.PrintDefaults()
}

func parseArgs() {
	flag.Parse()

	if options.Help {
		flag.CommandLine.SetOutput(os.Stdout)
		usage()
		os.Exit(0)
	}

	if flag.NArg() < 2 {
		usage()
		os.Exit(2)
	}

	options.Release = flag.Arg(0)
	options.Template = flag.Arg(1)
}

func main() {
	parseArgs()

	content, err := os.ReadFile(options.Template)
	must(err, "failed to read template")
	tmpl, err := template.New("release").Parse(string(content))
	must(err, "failed to process template")

	assets, err := getAssets(options.Release)
	must(err, "failed to get list of assets for release %s", options.Release)

	for _, asset := range assets {
		must(asset.DownloadAndChecksum(), "failed to calculate checksum for %s", asset.Name)
	}

	must(tmpl.Execute(os.Stdout, map[string][]*Asset{"assets": assets}), "failed to render template")
}
