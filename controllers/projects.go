package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/sotomskir/mastermind-server/models"
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
	project := models.GetProject(id)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

var SaveProjects = func(w http.ResponseWriter, r *http.Request) {
	var project models.Project
	err := json.NewDecoder(r.Body).Decode(&project)
	log.Println(err)
	if nil != err {
		// Simplified
		log.Println(err)
		return
	}
	err2 := models.SaveProject(&project)
	if nil != err2 {
		// Simplified
		log.Println(err2)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

var SynchronizeProject = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	project := models.GetProject(id)
	path := fmt.Sprintf("./storage/repositories/%s", project.Name)
	var (
		repo *git.Repository
		err error
	)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		keys, err := ssh.NewPublicKeys(project.RepoUser, []byte(project.SshKey.Key), "")
		if err != nil {
			utils.Error(w, err, http.StatusInternalServerError)
			return
		}
		logrus.Infoln("git clone")
		repo, err = git.PlainClone(fmt.Sprintf("./storage/repositories/%s", project.Name), false, &git.CloneOptions{
			URL:      project.RepoUrl,
			Progress: os.Stdout,
			Auth: keys,
		})
		if err != nil {
			utils.Error(w, err, http.StatusInternalServerError)
			return
		}
	} else {
		logrus.Infoln("git open")
		repo, err = git.PlainOpen(path)
		if err != nil {
			utils.Error(w, err, http.StatusInternalServerError)
			return
		}
		logrus.Infoln("git fetch")
		err = repo.Fetch(&git.FetchOptions{RemoteName: "origin"})
		if err != nil && err.Error() != "already up-to-date" {

			utils.Error(w, err, http.StatusInternalServerError)
			return
		}
	}
	logrus.Infoln("git tree")
	wTree, err := repo.Worktree()
	if err != nil {
		utils.Error(w, err, http.StatusInternalServerError)
		return
	}
	logrus.Infoln("git reset")
	err = wTree.Reset(&git.ResetOptions{Mode:git.HardReset, Commit:plumbing.NewHash(fmt.Sprintf("origin/%s", project.RepoBranch))})
	if err != nil {
		utils.Error(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}
