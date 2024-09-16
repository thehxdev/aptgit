package gpkg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/thehxdev/aptgit/config"
	"github.com/thehxdev/aptgit/gpath"
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

func (gp *Gpkg) ParseTagRegexp(tag string) string {
	if gp.TagRegexp != "" {
		tagRegexp := regexp.MustCompile(gp.TagRegexp)
		return tagRegexp.FindString(tag)
	}
	return tag
}

func (gp *Gpkg) GetLatestStableTag() (string, error) {
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
		// TODO: Check each type assertion to avoid panic
		allTags = append(allTags, t.(map[string]interface{})["tag_name"].(string))
	}

	return allTags, nil
}

func (gp *Gpkg) DownloadRelease(vars map[string]string) (string, error) {
	tag := vars["TAGNAME"]
	fileName := gvars.ResolveAll(gp.Template, vars)

	dlurl, err := url.JoinPath("https://github.com/", gp.Repository, "releases/download", tag, fileName)
	if err != nil {
		return "", err
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

	savePath := filepath.Join(config.G.DownloadPath, fileName)
	fp, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer fp.Close()
	writer := bufio.NewWriter(fp)

	jobChan := make(chan struct{}, 2)
	done := make(chan struct{}, 1)
	log.Inf.Println("Saving to", savePath)

	go func() {
		_, err = io.Copy(writer, resp.Body)
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

	return savePath, nil
}

func (gp *Gpkg) RunInstallSteps(vars map[string]string) error {
	err := os.MkdirAll(vars["INSTALL_PATH"], 0775)
	if err != nil {
		return err
	}

	normalizedCmds := make([][]string, 0)
	for _, step := range gp.InstallSteps {
		cmd := make([]string, 0)
		cmdWords := strings.Split(step, " ")
		for _, word := range cmdWords {
			resolved := gvars.ResolveAll(word, vars)
			if strings.Contains(resolved, " ") {
				resolved = gpath.Qoute(resolved)
			}
			cmd = append(cmd, resolved)
		}
		normalizedCmds = append(normalizedCmds, cmd)
	}

	for _, cmd := range normalizedCmds {
		cmd := exec.Command(cmd[0], cmd[1:]...)
		log.Inf.Printf("%+v\n", cmd)
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func (gp *Gpkg) SymlinkBinaryFiles(vars map[string]string) error {
	var err error = nil
	installPath := vars["INSTALL_PATH"]
	for _, bin := range gp.Bins {
		srcPath := filepath.Join(installPath, bin)
		_, binFile := filepath.Split(bin)
		destPath := filepath.Join(config.G.BinPath, binFile)
		log.Inf.Printf("%s -> %s", srcPath, destPath)
		err = os.Symlink(srcPath, destPath)
		if err != nil {
			goto ret
		}
	}
ret:
	return err
}
