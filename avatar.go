package main

import (
	"errors"
	"crypto/md5"
	"io"
	"strings"
	"fmt"
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
	if email, ok := c.userData["email"]; ok {
		if emailStr, ok := email.(string); ok {
			m := md5.New()//ハッシュ化
			io.WriteString(m, strings.ToLower(emailStr))
			return fmt.Sprintf("//www.gravatar.com/avatar/%x", m.Sum(nil)), nil
		}
	}

	return "", ErrNoAvatarURL
}