package main

import (
	"errors"

	"github.com/asaskevich/govalidator"
)

type User struct {
	Name     string `json:"Name"`
	Email    string `json:"Email"`
	GpgKeyID string `json:"GpgKeyId"`
}

func UsersToString(users []User) []string {
	us := make([]string, len(users))
	for i, user := range users {
		us[i] = user.Name + " <" + user.Email + ">"
	}

	return us
}

func ListUser() ([]User, error) {
	configPath, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	if !IsExist(configPath) {
		return nil, errors.New("No user")
	}

	config, err := ReadConfig()
	if err != nil {
		return nil, err
	}

	return config.Users, nil
}

func CreateUser(name, email, gpgKeyID string) error {
	configPath, err := ConfigPath()
	if err != nil {
		return err
	}

	user := User{Name: name, Email: email, GpgKeyID: gpgKeyID}

	if !IsExist(configPath) {
		if err := CreateConfig(Config{Users: []User{user}}); err != nil {
			return err
		}
		return nil
	}

	config, err := ReadConfig()
	if err != nil {
		return err
	}

	if config.Users == nil {
		config.Users = []User{user}
	} else {
		config.Users = append(config.Users, user)
	}

	if err := CreateConfig(config); err != nil {
		return err
	}

	return nil
}

func RemoveUser(idx int, users []User) error {
	newUsers := make([]User, len(users)-1)
	if idx+1 == len(users) {
		newUsers = users[:idx]
	} else {
		newUsers = append(users[:idx], users[idx+1:]...)
	}

	if err := CreateConfig(Config{Users: newUsers}); err != nil {
		return err
	}

	return nil
}

func ValidateEmail(email string) error {
	if !govalidator.IsExistingEmail(email) {
		return errors.New("Invalid email address")
	}

	return nil
}
