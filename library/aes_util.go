package library

import (
    "bytes"
    "crypto/aes"
    "crypto/cipher"
    "encoding/base64"
    "encoding/json"
    "gitlab.ceibsmoment.com/c/mp/logger"
)

type _aesUtil struct {
}

var (
    AesUtil *_aesUtil
)

func (o *_aesUtil) _pKCS7Padding(cipherText []byte, blockSize int) []byte {
    padding := blockSize - len(cipherText)%blockSize
    padText := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(cipherText, padText...)
}

func (o *_aesUtil) _pKCS7UnPadding(origData []byte) []byte {
    length := len(origData)
    unPadding := int(origData[length-1])
    return origData[:(length - unPadding)]
}

func (o *_aesUtil) Encrypt(key []byte, origData []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        logger.Logger.Errorf("_aesUtil.Encrypt() aes.NewCipher(%s) error: %#v", key, err)
        return nil, err
    }
    blockSize := block.BlockSize()
    origData = o._pKCS7Padding(origData, blockSize)
    blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
    crypted := make([]byte, len(origData))
    blockMode.CryptBlocks(crypted, origData)
    return crypted, nil
}

func (o *_aesUtil) Decrypt(key []byte, crypted []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        logger.Logger.Errorf("_aesUtil.Decrypt() aes.NewCipher(%s) error: %#v", key, err)
        return nil, err
    }
    blockSize := block.BlockSize()
    blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
    origData := make([]byte, len(crypted))
    blockMode.CryptBlocks(origData, crypted)
    origData = o._pKCS7UnPadding(origData)
    return origData, nil
}

func (o *_aesUtil) EncryptString(key string, OriginData interface{}) (resp string, err error) {
    var rs, data []byte
    data, err = json.Marshal(OriginData)
    if err != nil {
        logger.Logger.Errorf("_aesUtil.EncryptString() json.Marshal(%v) key = %s error: %v", OriginData, key, err)
        return
    }
    rs, err = o.Encrypt([]byte(key), data)
    if err != nil {
        logger.Logger.Errorf("_aesUtil.EncryptString(%s, $s) error: %v", key, OriginData, err)
        return
    }
    resp = base64.RawStdEncoding.EncodeToString(rs)

    return
}

func (o *_aesUtil) DecryptString(key string, crypted string) (resp interface{}, err error) {
    var rs, bs []byte
    bs, err = base64.RawStdEncoding.DecodeString(crypted)
    rs, err = o.Decrypt([]byte(key), bs)
    if err != nil {
        logger.Logger.Errorf("_aesUtil.DecryptString() o.Decrypt error: %v", key, err)
        return
    }

    err = json.Unmarshal(rs, &resp)
    if err != nil {
        logger.Logger.Errorf("_aesUtil.DecryptString() json.Unmarshal(%s, &resp) error: %v", string(rs), err)
        return
    }

    return
}
