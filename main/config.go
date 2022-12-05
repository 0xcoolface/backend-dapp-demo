package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
)

const ConfigFile = "config.json"

const (
	AccountSecretKey = iota
	AccountKeyStoreFile
)

type Config struct {
	AccountType  uint   `json:"accountType"`
	SecretHex    string `json:"secretHex,omitempty"`
	KeyStoreFile string `json:"keyStoreFile,omitempty"`
	Password     string `json:"password,omitempty"`
	RpcUrl       string `json:"rpcUrl"`

	isHttp bool
	secret *ecdsa.PrivateKey
}

func LoadConfig() *Config {
	bs, err := ioutil.ReadFile(ConfigFile)
	var conf Config
	if err == nil {
		err = json.Unmarshal(bs, &conf)
	}
	if err != nil {
		log.Panicf("load config failed. err=%s\n", err)
	}

	err = conf.loadSecret()
	if err != nil {
		log.Panicf("load secret failed. err=%s\n", err)
	}

	if strings.HasPrefix(strings.TrimSpace(conf.RpcUrl), "http") {
		conf.isHttp = true
	}

	return &conf
}

func (c *Config) loadSecret() error {
	switch c.AccountType {
	case AccountSecretKey:
		sk, err := crypto.HexToECDSA(c.SecretHex)
		if err != nil {
			return err
		}
		c.secret = sk
	case AccountKeyStoreFile:
		// Load the key from the keystore and decrypt its contents
		keyjson, err := ioutil.ReadFile(c.KeyStoreFile)
		if err != nil {
			return err
		}
		key, err := keystore.DecryptKey(keyjson, c.Password)
		if err != nil {
			return err
		}
		c.secret = key.PrivateKey

	default:
		return fmt.Errorf("unsupported AccountType %d", c.AccountType)
	}
	return nil
}

func (c *Config) PrivateKey() *ecdsa.PrivateKey {
	return c.secret
}
