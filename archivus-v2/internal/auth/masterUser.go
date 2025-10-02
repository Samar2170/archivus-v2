package auth

import (
	"archivus-v2/internal/db"
	"archivus-v2/internal/models"
	"errors"
	"os"
	"os/exec"
)

func SudoCheck() error {
	cmd := exec.Command("sudo", "-v")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func CheckMasterUser() error {
	var count int64
	db.StorageDB.Model(&models.User{}).Where("is_master = ?", true).Count(&count)
	if count == 0 {
		return errors.New(MasterUserNotFoundErr)
	}
	return nil
}
