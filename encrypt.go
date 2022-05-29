package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
)

func EncryptPassword(password string) string {
	const pubKey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDl/aCgRl9f/4ON9MewoVnV58OL
OU2ALBi2FKc5yIsfSpivKxe7A6FitJjHva3WpM7gvVOinMehp6if2UNIkbaN+plW
f5IwqEVxsNZpeixc4GsbY9dXEk3WtRjwGSyDLySzEESH/kpJVoxO7ijRYqU+2oSR
wTBNePOk1H+LRQokgQIDAQAB
-----END PUBLIC KEY-----`
	encryptedPassword, err := encryptByPublicKey(password, pubKey)
	if err != nil {
		panic(err)
	}
	return encryptedPassword
}

func encryptByPublicKey(data, pubKey string) (string, error) {
	block, _ := pem.Decode([]byte(pubKey))
	if block == nil {
		return "", errors.New("decode public key error: PEM structure is invalid")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	output := &bytes.Buffer{}
	err = pubKeyIO(pub.(*rsa.PublicKey), bytes.NewReader([]byte(data)), output)
	if err != nil {
		return "", err
	}

	all, err := ioutil.ReadAll(output)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(all), nil
}

func pubKeyIO(pub *rsa.PublicKey, in io.Reader, out io.Writer) (err error) {
	k := (pub.N.BitLen() + 7) / 8
	k -= 11
	buf := make([]byte, k)
	var b []byte
	size := 0
	for {
		size, err = in.Read(buf)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if size < k {
			b = buf[:size]
		} else {
			b = buf
		}
		b, err = rsa.EncryptPKCS1v15(rand.Reader, pub, b)
		if err != nil {
			return err
		}
		if _, err = out.Write(b); err != nil {
			return err
		}
	}
}
