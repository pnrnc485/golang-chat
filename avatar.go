package main

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	_ "net/url"
)

//ErrNoAvatarはAvatarインスタンスがアバターのURLを返すことができない場合に発生するエラーです
var ErrNoAvatarURL = errors.New("chat: アバターのURLを取得できません。")

//ユーザのプロフィール画像を表す
type Avatar interface {

	//指定されたChatUserのアバターのURLを返す　取得できなかった場合はErrNoAvatarを返す　問題が発錆した場合はエラーを返す
	GetAvatarURL(ChatUser) (string, error)
}

type  AuthAvatar struct { }
var UseAuthAvatar AuthAvatar
func (_ AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()
	if url != "" {
		return url, nil
	}

	return "", ErrNoAvatarURL
}

type GravatarAvatar struct { }
var UseGravatar GravatarAvatar
func (_ GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}

type FileSystemAvatar struct {}
var UsefileSystemAvatar FileSystemAvatar
func (_ FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	if files, err := ioutil.ReadDir("avatars"); err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if match, _ := filepath.Match(u.UniqueID()+"*", file.Name()); match {
				return "/avatars/" + file.Name(), nil
			}
		}
	}
	return "", ErrNoAvatarURL
}