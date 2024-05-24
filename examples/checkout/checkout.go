package main

import (
	"fmt"
	"github.com/MR5356/go-workflow"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/sirupsen/logrus"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

type CheckoutTask struct {
	repository string
	branch     string
	submodules bool

	workflow.UnimplementedITask
}

func (t *CheckoutTask) SetParams(params *workflow.TaskParams) error {
	return nil
}

func (t *CheckoutTask) Run() error {
	srcAuth, repoUrl, err := getAuth(t.repository, "", "")
	if err != nil {
		return err
	}

	dirName := fmt.Sprintf("/tmp/%s", strings.ReplaceAll(filepath.Base(repoUrl), ".git", ""))
	if _, err := os.Stat(dirName); !os.IsNotExist(err) {
		_ = os.RemoveAll(dirName)
	}

	workflow.Logger.Debug("clone %s to %s", repoUrl, dirName)

	co := &git.CloneOptions{
		URL:               repoUrl,
		Auth:              srcAuth,
		ReferenceName:     plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", t.branch)),
		RecurseSubmodules: git.NoRecurseSubmodules,
		InsecureSkipTLS:   true,
	}

	repo, err := git.PlainClone(dirName, false, co)
	if err != nil {
		return err
	}

	if t.submodules {
		worktree, err := repo.Worktree()
		if err != nil {
			return err
		}
		submodules, err := worktree.Submodules()
		if err != nil {
			return err
		}

		for _, sm := range submodules {
			logrus.Debugf("clone submodule: %s", sm.Config().URL)
			smAuth, _, _ := getAuth(sm.Config().URL, "", "")
			err = sm.Update(&git.SubmoduleUpdateOptions{
				Init: true,
				Auth: smAuth,
			})

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func getAuth(repo, privateKeyFile, privateKeyPassword string) (auth transport.AuthMethod, repoUrl string, err error) {
	/**
	支持以下形式：
		1. https://github.com/MR5356/syncer.git
		2. git@github.com:MR5356/syncer.git
		3. https://username:password@github.com/MR5356/syncer.git
		4. https://<token>@github.com/MR5356/syncer.git
		5. https://oauth2:access_token@github.com/MR5356/syncer.git
	*/

	repoUrl = repo
	switch getUrlType(repo) {
	case gitUrlType:
		if privateKeyFile == "" {
			u, err := user.Current()
			if err == nil {
				privateKeyFile = fmt.Sprintf("%s/.ssh/id_rsa", u.HomeDir)
			}
		}
		_, err = os.Stat(privateKeyFile)
		if err != nil {
			return auth, repoUrl, err
		}
		workflow.Logger.Debug("privateKeyFile: %s, privateKeyPassword: %s", privateKeyFile, privateKeyPassword)
		auth, err = ssh.NewPublicKeysFromFile("git", privateKeyFile, privateKeyPassword)
		return auth, repoUrl, err
	case httpUrlType:
		auth = nil
		return
	case tokenizedHttpUrlType:
		token := strings.ReplaceAll(strings.ReplaceAll(strings.Split(repo, "@")[0], "https://", ""), "http://", "")
		workflow.Logger.Info(token)
		auth = &http.TokenAuth{
			Token: token,
		}
		repoUrl = strings.ReplaceAll(repo, token+"@", "")
		return
	case basicHttpUrlType:
		basicInfo := strings.ReplaceAll(strings.ReplaceAll(strings.Split(repo, "@")[0], "https://", ""), "http://", "")
		fields := strings.Split(basicInfo, ":")
		auth = &http.BasicAuth{
			Username: fields[0],
			Password: fields[1],
		}
		repoUrl = strings.ReplaceAll(repo, basicInfo+"@", "")
		return
	default:
		return nil, "", fmt.Errorf("unsupported repo url: %s", repo)
	}
}

type urlType int

const (
	unknownUrlType urlType = iota
	gitUrlType
	httpUrlType
	tokenizedHttpUrlType
	basicHttpUrlType
)

var (
	isGitUrl           = regexp.MustCompile(`^git@[-\w.:]+:[-\/\w.]+\.git$`)
	isHttpUrl          = regexp.MustCompile(`^(https|http)://[-\w.:]+/[-\/\w.]+\.git$`)
	isTokenizedHttpUrl = regexp.MustCompile(`^(https|http)://[a-zA-Z0-9_]+@[-\w.:]+/[-\/\w.]+\.git$`)
	isBasicHttpUrl     = regexp.MustCompile(`^(https|http)://[a-zA-Z0-9]+:[\w]+@[-\w.:]+/[-\/\w.]+\.git$`)
)

func getUrlType(url string) (t urlType) {
	if isGitUrl.MatchString(url) {
		t = gitUrlType
	} else if isHttpUrl.MatchString(url) {
		t = httpUrlType
	} else if isTokenizedHttpUrl.MatchString(url) {
		t = tokenizedHttpUrlType
	} else if isBasicHttpUrl.MatchString(url) {
		t = basicHttpUrlType
	} else {
		t = unknownUrlType
	}
	workflow.Logger.Debug("getUrlType: %v", t)
	return t
}

func main() {
	workflow.Serve(&CheckoutTask{})
}
