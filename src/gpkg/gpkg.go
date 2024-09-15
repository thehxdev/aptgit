package gpkg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"time"

	"github.com/thehxdev/aptgit/gvars"
	"github.com/thehxdev/aptgit/log"
)

const GH_API_URL = "https://api.github.com/repos"

type Gpkg struct {
	Repository   string            `json:"repository"`
	PlatformMap  map[string]string `json:"platform,omitempty"`
	ArchMap      map[string]string `json:"arch,omitempty"`
	TagRegexp    string            `json:"tagRegexp,omitempty"`
	Template     string            `json:"template"`
	InstallSteps []string          `json:"install"`
	Bins         []string          `json:"bins"`
}

func ReadDefinitionFile(p string) (*Gpkg, error) {
	defContent, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}

	pdef := &Gpkg{}

	err = json.Unmarshal(defContent, pdef)
	if err != nil {
		return nil, err
	}

	return pdef, nil
}

func (gp *Gpkg) GetArch(a string) string {
	if garch, ok := gp.ArchMap[a]; ok {
		return garch
	}
	return a
}

func (gp *Gpkg) GetPlatform(p string) string {
	if gplat, ok := gp.PlatformMap[p]; ok {
		return gplat
	}
	return p
}

func (gp *Gpkg) ParseTagRegexp(v string) string {
	if gp.TagRegexp != "" {
		verRegexp := regexp.MustCompile(gp.TagRegexp)
		return verRegexp.FindString(v)
	}
	return v
}

func (gp *Gpkg) GetLatestTag() (string, error) {
	req_url, err := url.JoinPath(GH_API_URL, gp.Repository, "releases/latest")
	if err != nil {
		return "", err
	}

	resp, err := http.Get(req_url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("got %d status code", resp.StatusCode)
	}

	jtag := struct {
		TagName string `json:"tag_name"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&jtag)
	if err != nil {
		return "", err
	}

	return jtag.TagName, nil
}

func (gp *Gpkg) GetAllTags() ([]string, error) {
	req_url, err := url.JoinPath(GH_API_URL, gp.Repository, "releases")
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(req_url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got %d status code", resp.StatusCode)
	}

	jtags := make([]interface{}, 0)

	err = json.NewDecoder(resp.Body).Decode(&jtags)
	if err != nil {
		return nil, err
	}

	allTags := make([]string, 0)
	for _, t := range jtags {
		allTags = append(allTags, t.(map[string]interface{})["tag_name"].(string))
	}

	return allTags, nil
}

func (gp *Gpkg) DownloadLatest(outDir string) error {
	latestTag, err := gp.GetLatestTag()
	if err != nil {
		return err
	}

	fileName := gvars.ResolveAll(gp.Template, map[string]string{
		"TAGNAME":  latestTag,
		"VERSION":  gp.ParseTagRegexp(latestTag),
		"PLATFORM": gp.GetPlatform(runtime.GOOS),
		"ARCH":     gp.GetArch(runtime.GOARCH),
	})

	dlurl, err := url.JoinPath("https://github.com/", gp.Repository, "releases/download", latestTag, fileName)
	if err != nil {
		return err
	}

	log.Inf.Println("Downloading", dlurl)
	resp, err := http.Get(dlurl)
	if err != nil {
		log.Err.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Err.Fatalf("Got %d status code", resp.StatusCode)
	}

	var fileSize float64
	cl, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil || cl == 0 {
		log.Err.Println("Could not get Content-Length")
	} else {
		fileSize = float64(cl) / 1024 / 1024
		log.Inf.Printf("File size: %.3f MiB", fileSize)
	}

	savePath := filepath.Join(outDir, fileName)
	fp, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer fp.Close()

	jobChan := make(chan struct{}, 2)
	done := make(chan struct{}, 1)
	log.Inf.Println("Saving to", savePath)

	go func() {
		_, err = io.Copy(fp, resp.Body)
		if err != nil {
			log.Err.Fatal(err)
		}
		done <- struct{}{}
		jobChan <- struct{}{}
	}()

	go func() {
	showProgress:
		for {
			time.Sleep(time.Second * 1)
			stat, err := fp.Stat()
			if err != nil {
				log.Err.Println(err)
				continue
			}
			downloaded := float64(stat.Size()) / 1024 / 1024
			fmt.Printf("\r%.3f MiB of %.3f MiB Downloaded...", downloaded, fileSize)
			select {
			case <-done:
				fmt.Print("\n")
				break showProgress
			default:
			}
		}
		jobChan <- struct{}{}
	}()

	for i := 0; i < 2; i++ {
		<-jobChan
	}
	log.Inf.Println("Success!")

	return nil
}
