# NaClSecretBox

NaClSecretBox is a utility to encrypt and decrypt multiple files using
the same key but a different nonce.

1. To use NaClSecretBox, first install the Go language. Instructions for installation on Ubuntu 16.04 are [here](https://www.digitalocean.com/community/tutorials/how-to-install-go-1-6-on-ubuntu-16-04). (If you're installing locally, you can skip the first step of that tutorial.) Be sure to correctly set the GOPATH environment variable and test that your installation works by creating and building the Hello World program in the tutorial. 

2. To install NaClSecretBox, first download the program:
```
go get github.com/marcopeereboom/naclsbox
```
3. Next, move the NaClSecretBox download into the /src directory you created when you installed Go and tested your installation with the Hello World program. 

4. To build and run the NaClSecretBox program, follow the instructions [here](https://golang.org/cmd/go/#hdr-Compile_packages_and_dependencies).

5. Now you're ready to encrypt and decrypt files. Once you encrypt a file, you will need to send the file blob and the encryption key to the recipient. It's advisable to share these two items through different, encrypted channels. 

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
