package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/deploji/deploji-server/dto"
	"github.com/deploji/deploji-server/models"
	"github.com/deploji/deploji-server/services"
	"github.com/deploji/deploji-server/services/auth"
	"github.com/deploji/deploji-server/utils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

var Encrypt = func(w http.ResponseWriter, r *http.Request) {
	var encrypt dto.Encrypt
	err := json.NewDecoder(r.Body).Decode(&encrypt)
	log.Println(err)
	if nil != err {
		utils.Error(w, "Cannot decode encrypt request", err, http.StatusInternalServerError)
		return
	}
	jwt := services.GetJWTClaims(r)
	if !auth.Enforce(jwt, dto.ObjectTypeSshKey, encrypt.KeyID, dto.ActionTypeUse) {
		utils.Error(w, fmt.Sprintf("Access denied to vault key: %d", encrypt.KeyID), err, http.StatusForbidden)
		return
	}
	key := models.GetSshKey(encrypt.KeyID)
	vaultKeyPath, err := WriteKey(encrypt.KeyID, key.Key)
	if err != nil {
		utils.Error(w, "", err, http.StatusInternalServerError)
		return
	}
	cmd := exec.Command("ansible-vault", "encrypt_string", "--name", encrypt.Name, "--vault-id", vaultKeyPath, fmt.Sprintf("%s", encrypt.Content))
	stdout, err := cmd.CombinedOutput()

	if err != nil {
		utils.Error(w, "", err, http.StatusInternalServerError)
		log.Println(string(stdout))
		return
	}
	encrypt.Content = string(stdout)
	json.NewEncoder(w).Encode(encrypt)
}

func WriteKey(id uint, content models.Key) (string, error) {
	if err := os.MkdirAll("storage/keys", 0700); err != nil {
		log.Printf("Error creating directory: %s", err)
		return "", err
	}
	if err := ioutil.WriteFile(fmt.Sprintf("storage/keys/%d", id), []byte(content), 0600); err != nil {
		log.Printf("Error saving key file: %s", err)
		return "", err
	}
	os.Chmod(fmt.Sprintf("storage/keys/%d", id), 0600)
	return fmt.Sprintf("storage/keys/%d", id), nil
}
