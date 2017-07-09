package main

import (
	"errors"
	"io/ioutil"
	"path/filepath"
)

//ErrNoAvatarはAvatarインスタンスがアバターのURLを返すことができない場合に発生するエラーです
var ErrNoAvatarURL = errors.New("chat: アバターのURLを取得できません。")

//ユーザのプロフィール画像を表す
type Avatar interface {

	//指定されたclientのアバターのURLを返す　取得できなかった場合はErrNoAvatarを返す　問題が発錆した場合はエラーを返す
	GetAvatarURL(c *client) (string, error)
}

type  AuthAvatar struct { }
var UseAuthAvatar AuthAvatar
func (_ AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}

	return "", ErrNoAvatarURL
}

type GravatarAvatar struct { }
var UseGravatar GravatarAvatar
func (_ GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			return "//www.gravatar.com/avatar/" + useridStr, nil
		}
	}

	return "", ErrNoAvatarURL
}

type FileSystemAvatar struct {}
var UsefileSystemAvatar FileSystemAvatar
func (_ FileSystemAvatar) GetAvatarURL(c * client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			if files, err := ioutil.ReadDir("avatars"); err == nil {
				for _, file := range files {
					if file.IsDir() {
						continue
					}
					if match, _ := filepath.Match(useridStr+"*", file.Name()); match {
						return "/avatars/" + file.Name(), nil
					}
				}
			}
		}
	}
	return "", ErrNoAvatarURL
}