package main

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// ErrNoAvatarURL インスタンスがアバターのURLを返すことができない場合に発生するエラー
var ErrNoAvatarURL = errors.New("chat: アバターのURLを取得できません。")

// Avatar ユーザーのプロフィール画像を表す型
type Avatar interface {
	// GetAvatarURL 指定されたクライアントのアバターのURLを返す
	// *問題が発生した場合にはエラーを返す
	// *URLを取得できなかった場合にはErrNoAvatarURLを返す
	GetAvatarURL(ChatUser) (string, error)
}

// TryAvatars 3つのアバター機能を格納
type TryAvatars []Avatar

// GetAvatarURL 3つのアバター機能の振り分け. 下記の順番で実装される
// FileSystemAvatar → AuthAvatar → Gravatar
func (a TryAvatars) GetAvatarURL(u ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.GetAvatarURL(u); err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}

// AuthAvatar 認証サービスを使用したアバター
type AuthAvatar struct{}

// UseAuthAvatar AuthAvatarを使うことを明示的にするため変数としている
var UseAuthAvatar AuthAvatar

// GetAvatarURL Receiver:AuthAvatar
func (AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()
	if url != "" {
		return url, nil
	}
	return "", ErrNoAvatarURL
}

// GravatarAvatar Gravatarを使用したアバター
type GravatarAvatar struct{}

// UseGravatar GravatarAvatarを使うことを明示的にするため変数としている
var UseGravatar GravatarAvatar

// GetAvatarURL Receiver:GravatarAvatar
func (GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}

// FileSystemAvatar FileSystemを使用したアバター
type FileSystemAvatar struct{}

// UseFileSystemAvatar FileSystemAvatarを使うことを明示的にするため変数としている
var UseFileSystemAvatar FileSystemAvatar

// GetAvatarURL Receiver:FileSystemAvatar
func (FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	files, err := ioutil.ReadDir("avatars")
	if err != nil {
		return "", ErrNoAvatarURL
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filename := file.Name()
		if u.UniqueID() == strings.TrimSuffix(filename, filepath.Ext(filename)) {
			return "/avatars/" + filename, nil
		}
	}
	return "", ErrNoAvatarURL
}
