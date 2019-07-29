package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/sotomskir/mastermind-server/models"
	"github.com/sotomskir/mastermind-server/services"
	"github.com/sotomskir/mastermind-server/utils"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"log"
	"net/http"
	"os"
	"strconv"
)

var GetProjects = func(w http.ResponseWriter, r *http.Request) {
	projects := models.GetProjects()
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

var GetProject = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	project := models.GetProject(uint(id))
	if project == nil {
		utils.Error(w, "Cannot load project", errors.New("not found"), http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

var GetProjectFiles = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	files, err := services.GetProjectFiles(uint(id))
	if err != nil {
		utils.Error(w, "Cannot load project files", err, http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

var SaveProjects = func(w http.ResponseWriter, r *http.Request) {
	var project models.Project
	err := json.NewDecoder(r.Body).Decode(&project)
	log.Println(err)
	if nil != err {
		utils.Error(w, "Cannot decode project", err, http.StatusInternalServerError)
		return
	}
	err = models.SaveProject(&project)
	if nil != err {
		utils.Error(w, "Cannot save project", err, http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

var DeleteProject = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	project := models.GetProject(uint(id))
	if project == nil {
		utils.Error(w, "Cannot load project", errors.New("not found"), http.StatusNotFound)
		return
	}
	err := models.DeleteProject(project)
	if err != nil {
		utils.Error(w, "Cannot delete project", err, http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
}

var SynchronizeProject = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	if err := SynchronizeProjectRepo(uint(id)); err != nil {
		utils.Error(w, "Cannot synchronize project", err, http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
}

func SynchronizeProjectRepo(id uint) error {
	project := models.GetProject(id)
	path := fmt.Sprintf("./storage/repositories/%s", project.Name)
	var (
		repo *git.Repository
		err error
	)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		keys, err := ssh.NewPublicKeys(project.RepoUser, []byte(project.SshKey.Key), "")
		if err != nil {
			log.Printf("NewPublicKeys: %s", err)
			return err
		}
		logrus.Infoln("git clone")
		repo, err = git.PlainClone(fmt.Sprintf("./storage/repositories/%s", project.Name), false, &git.CloneOptions{
			URL:      project.RepoUrl,
			Progress: os.Stdout,
			Auth: keys,
		})
		if err != nil {
			log.Printf("git clone: %s", err)
			return err
		}
	} else {
		logrus.Infoln("git open")
		repo, err = git.PlainOpen(path)
		if err != nil {
			log.Printf("git open: %s", err)
			return err
		}
		logrus.Infoln("git fetch")
		err = repo.Fetch(&git.FetchOptions{RemoteName: "origin"})
		if err != nil && err.Error() != "already up-to-date" {
			log.Printf("git fetch: %s", err)
			return err
		}
	}
	logrus.Infoln("git tree")
	wTree, err := repo.Worktree()
	if err != nil {
		log.Printf("git tree: %s", err)
		return err
	}
	logrus.Infof("git rev-parse origin/%s", project.RepoBranch)
	hash, err := repo.ResolveRevision(plumbing.Revision(fmt.Sprintf("origin/%s", project.RepoBranch)))
	if err != nil {
		log.Printf("git rev-parse: %s", err)
		return err
	}
	logrus.Infof("git reset --hard %s", hash.String())
	err = wTree.Reset(&git.ResetOptions{Mode:git.HardReset, Commit:*hash})
	if err != nil {
		log.Printf("git reset: %s", err)
		return err
	}
	return nil
}
