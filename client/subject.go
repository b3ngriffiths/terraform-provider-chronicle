package client

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type Subject struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Roles []Role `json:"roles,omitempty"`
}

type Role struct {
	Name        string       `json:"name"`
	Title       string       `json:"title,omitempty"`
	Description string       `json:"description,omitempty"`
	CreateTime  string       `json:"createTime,omitempty"`
	IsDefault   string       `json:"isDefault,omitempty"`
	Permissions []Permission `json:"permissions,omitempty"`
}

type Permission struct {
	Name        string `json:"name"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	CreateTime  string `json:"createTime,omitempty"`
}

func (cli *Client) GetSubject(name string) (*Subject, error) {
	url := fmt.Sprintf("%s/%s", cli.SubjectsBasePath, name)

	err := cli.rateLimiters.RBACGetSubject.Wait(cli.context)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error waiting for rateLimiter getting subject %s", name))
	}

	res, err := sendRequest(cli, cli.backstoryAPIClient, "GET", cli.userAgent, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed getting subject")
	}

	var subject Subject
	err = json.Unmarshal(res, &subject)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal subject response")
	}

	return &subject, nil
}

func (cli *Client) CreateSubject(subject Subject) error {
	url := cli.SubjectsBasePath

	err := cli.rateLimiters.RBACCreateSubject.Wait(cli.context)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error waiting for rateLimiter creating subject %v", subject))
	}

	_, err = sendRequest(cli, cli.backstoryAPIClient, "POST", cli.userAgent, url, subject)
	if err != nil {
		return errors.Wrap(err, "failed creating subject")
	}

	return nil
}

func (cli *Client) UpdateSubject(subject Subject) error {
	url := fmt.Sprintf("%s/%s", cli.SubjectsBasePath, subject.Name)

	err := cli.rateLimiters.RBACUpdateSubject.Wait(cli.context)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error waiting for rateLimiter updating subject %v", subject))
	}

	_, err = sendRequest(cli, cli.backstoryAPIClient, "PATCH", cli.userAgent, url, subject)
	if err != nil {
		return errors.Wrap(err, "failed updating subject")
	}

	return nil
}

func (cli *Client) DeleteSubject(name string) error {
	url := fmt.Sprintf("%s/%s", cli.SubjectsBasePath, name)

	err := cli.rateLimiters.RBACDeleteSubject.Wait(cli.context)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error waiting for rateLimiter deleting subject %s", name))
	}

	_, err = sendRequest(cli, cli.backstoryAPIClient, "DELETE", cli.userAgent, url, nil)
	if err != nil {
		return errors.Wrap(err, "failed deleting subject")
	}

	return nil
}
