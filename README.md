# NaClSecretBox

NaClSecretBox is a utility to encrypt and decrypt multiple files using
the same key but a different nonce.

Installation:
```
go get github.com/marcopeereboom/naclsbox
```

Encryption example:
```
$ naclsbox -e myfile 
encryption key: 0c4d1b5840347b698b92e7c33a87f52e721616fede4c82ce22a0783001cc6c92
```
This generated a file called `myfile.sbox` which contains the encrypted data.

**Take note of the encryption key.  It can not be recovered.**

In order to reverse the process must provide the key and encrypted blob filename.

Decryption example:
```
$ naclsbox -d -k 0c4d1b5840347b698b92e7c33a87f52e721616fede4c82ce22a0783001cc6c92 myfile.sbox
```
This generated a file called `myfile.sbox.decrypted` which should be an identical copy of `myfile`.

Verification:
```
$ shasum myfile
f951b101989b2c3b7471710b4e78fc4dbdfa0ca6  myfile
$ shasum myfile.sbox.decrypted 
f951b101989b2c3b7471710b4e78fc4dbdfa0ca6  myfile.sbox.decrypted
```
