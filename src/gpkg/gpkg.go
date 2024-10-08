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
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/thehxdev/aptgit/genv"
	"github.com/thehxdev/aptgit/gpath"
	"github.com/thehxdev/aptgit/gvars"
	"github.com/thehxdev/aptgit/log"
)

const GH_API_URL = "https://api.github.com/repos"

type PkgInfo struct {
	Repository   string            `json:"repository"`
	PlatformMap  map[string]string `json:"platform,omitempty"`
	ArchMap      map[string]string `json:"arch,omitempty"`
	TagRegexp    string            `json:"tagRegexp,omitempty"`
	Template     string            `json:"template"`
	InstallSteps []string          `json:"install"`
	Bins         []string          `json:"bins"`
}

type Gpkg struct {
	Info    *PkgInfo
	TagName string
	Vars    map[string]string
	MainTag string
}

func Init(pdef string) (*Gpkg, error) {
	info, err := readDefinitionFile(pdef)
	if err != nil {
		return nil, err
	}

	gp := &Gpkg{
		Info: info,
	}

	return gp, nil
}

func ReadMdFile() (map[string]string, error) {
	allMds := make(map[string]string)

	fp, err := os.Open(genv.G.LockFile)
	if err != nil {
		return allMds, err
	}
	defer fp.Close()
	reader := bufio.NewReader(fp)

	err = json.NewDecoder(reader).Decode(&allMds)
	if err != nil {
		return allMds, err
	}

	return allMds, nil
}

func WriteMdFile(mds map[string]string) error {
	fp, err := os.Create(genv.G.LockFile)
	if err != nil {
		return err
	}
	defer fp.Close()

	return json.NewEncoder(fp).Encode(mds)
}

func readDefinitionFile(p string) (*PkgInfo, error) {
	defContent, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}

	pdef := &PkgInfo{}

	err = json.Unmarshal(defContent, pdef)
	if err != nil {
		return nil, err
	}

	return pdef, nil
}

func (gp *Gpkg) GetArch() string {
	a := runtime.GOARCH
	if garch, ok := gp.Info.ArchMap[a]; ok {
		a = garch
	}
	return a
}

func (gp *Gpkg) GetPlatform() string {
	p := runtime.GOOS
	if gplat, ok := gp.Info.PlatformMap[p]; ok {
		p = gplat
	}
	return p
}

func (gp *Gpkg) ParseTagRegexp(tag string) string {
	if gp.Info.TagRegexp != "" {
		tagRegexp := regexp.MustCompile(gp.Info.TagRegexp)
		return tagRegexp.FindString(tag)
	}
	return tag
}

func (gp *Gpkg) GetLatestStableTag() (string, error) {
	req_url, err := url.JoinPath(GH_API_URL, gp.Info.Repository, "releases/latest")
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
	req_url, err := url.JoinPath(GH_API_URL, gp.Info.Repository, "releases")
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

	jdata := make([]interface{}, 0)

	err = json.NewDecoder(resp.Body).Decode(&jdata)
	if err != nil {
		return nil, err
	}

	allTags := make([]string, 0)
	for _, t := range jdata {
		if t, ok := t.(map[string]interface{}); ok {
			if tname, ok := t["tag_name"].(string); ok {
				allTags = append(allTags, tname)
			}
		}
	}

	return allTags, nil
}

func (gp *Gpkg) DownloadRelease(vars map[string]string) (string, error) {
	tag := vars["TAGNAME"]
	fileName := gvars.ResolveAll(gp.Info.Template, vars)

	dlpath := genv.G.DownloadPath
	err := gpath.MkdirIfNotExist(dlpath)
	if err != nil {
		log.Err.Printf("Failed to create %s path", dlpath)
		return "", err
	}

	dlurl, err := url.JoinPath("https://github.com/", gp.Info.Repository, "releases/download", tag, fileName)
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

	savePath := filepath.Join(dlpath, fileName)
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
				log.Wrn.Println(err)
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

	return savePath, nil
}

func RunCommands(commands []string, vars map[string]string) error {
	normalizedCmds := make([][]string, 0)
	for _, cmd := range commands {
		newCmd := make([]string, 0)
		cmdWords := strings.Split(cmd, " ")
		for _, word := range cmdWords {
			resolved := gvars.ResolveAll(word, vars)
			if strings.ContainsAny(resolved, " \\") {
				resolved = gpath.Qoute(resolved)
			}
			newCmd = append(newCmd, resolved)
		}
		normalizedCmds = append(normalizedCmds, newCmd)
	}

	for _, cmd := range normalizedCmds {
		cmd := exec.Command(cmd[0], cmd[1:]...)
		log.Inf.Println(cmd.String())
		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func (gp *Gpkg) SymlinkBinaryFiles(vars map[string]string) error {
	var err error

	binsPath := genv.G.BinPath
	err = gpath.MkdirIfNotExist(binsPath)
	if err != nil {
		return fmt.Errorf("failed to create %s path", binsPath)
	}

	installPath := vars["INSTALL_PATH"]
	for _, bin := range gp.Info.Bins {
		srcPath := filepath.Join(installPath, bin)
		_, binFile := filepath.Split(bin)
		destPath := filepath.Join(binsPath, binFile)
		log.Inf.Printf("%s -> %s", srcPath, destPath)
		err = os.Symlink(srcPath, destPath)
		if err != nil {
			goto ret
		}
	}
ret:
	return err
}

func (gp *Gpkg) Install() error {
	var err error

	gInstallPath := genv.G.InstallPath
	err = gpath.MkdirIfNotExist(gInstallPath)
	if err != nil {
		return fmt.Errorf("failed to create %s path", gInstallPath)
	}

	gp.Vars["FILE"], err = gp.DownloadRelease(gp.Vars)
	if err != nil {
		return err
	}

	pkgInstallPath := gp.Vars["INSTALL_PATH"]
	if gpath.Exist(pkgInstallPath) {
		os.Remove(pkgInstallPath)
	}

	err = os.MkdirAll(pkgInstallPath, 0775)
	if err != nil {
		return err
	}

	err = RunCommands(gp.Info.InstallSteps, gp.Vars)
	if err != nil {
		return err
	}

	gp.RemoveExistingSymlinks()
	err = gp.SymlinkBinaryFiles(gp.Vars)
	if err != nil {
		log.Err.Fatal(err)
	}

	gp.MainTag = gp.TagName
	return nil
}

func (gp *Gpkg) RemoveExistingSymlinks() {
	for _, bin := range gp.Info.Bins {
		_, filename := filepath.Split(bin)
		path := filepath.Join(genv.G.BinPath, filename)
		if err := os.Remove(path); err != nil {
			log.Wrn.Println(err)
		}
	}
}

func (gp *Gpkg) SetTagNameAsMain() error {
	tagInstallPath := gp.Vars["INSTALL_PATH"]
	if !gpath.Exist(tagInstallPath) {
		return fmt.Errorf("tag name %s is not installed for package %s", gp.TagName, gp.Info.Repository)
	}

	gp.RemoveExistingSymlinks()
	gp.SymlinkBinaryFiles(gp.Vars)

	gp.MainTag = gp.TagName
	return nil
}

func (gp *Gpkg) Uninstall() error {
	var err error
	tag := gp.TagName

	pkgInstallPath := path.Join(genv.G.InstallPath, gp.Info.Repository)
	if !gpath.Exist(pkgInstallPath) {
		return fmt.Errorf("package %s is not installed", gp.Info.Repository)
	}

	if tag == "all" {
		if err = os.RemoveAll(path.Dir(pkgInstallPath)); err != nil {
			return err
		}
		goto finalize
	}

	err = os.RemoveAll(path.Join(pkgInstallPath, tag))
	if err != nil {
		return err
	}

finalize:
	gp.RemoveExistingSymlinks()
	return nil
}
