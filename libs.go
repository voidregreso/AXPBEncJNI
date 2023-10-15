package main

//
// #include <stdlib.h>
// #include <stddef.h>
// #include <stdint.h>
import "C"
import (
	"fmt"
	"github.com/ClarkGuan/jni"
	"os"
)

//export jniOnLoad
func jniOnLoad(vm uintptr) {
	// TODO
	fmt.Println("JNI_OnLoad")
}

//export jniOnUnload
func jniOnUnload(vm uintptr) {
	// TODO
	fmt.Println("JNI_OnUnload")
}

//export jni_com_ernesto_cifrador_encryptFile1
func jni_com_ernesto_cifrador_encryptFile1(env uintptr, clazz uintptr, path uintptr, newPath uintptr, key uintptr) uintptr {
	// TODO
	ckey := jni.Env(env).GetStringUTF(key)
	cpath := jni.Env(env).GetStringUTF(path)
	cnPath := jni.Env(env).GetStringUTF(newPath)
	fmt.Printf("Encrypt Mode: {Input:%s, Output:%s, Key: %s}\n", string(cpath), string(cnPath), string(ckey))
	f, err := os.Open(string(cpath))
	strResult := "Success"
	if err != nil {
		fmt.Println("could not open file", string(cpath))
		strResult = err.Error()
	} else {
		err = encryptFile(ckey, f, string(cnPath))
		if err != nil {
			strResult = err.Error()
		}
	}
	return jni.Env(env).NewString(strResult)
}

//export jni_com_ernesto_cifrador_decryptFile2
func jni_com_ernesto_cifrador_decryptFile2(env uintptr, clazz uintptr, path uintptr, newPath uintptr, key uintptr) uintptr {
	// TODO
	cpath := jni.Env(env).GetStringUTF(path)
	cnPath := jni.Env(env).GetStringUTF(newPath)
	ckey := jni.Env(env).GetStringUTF(key)
	fmt.Printf("Decrypt Mode: {Input:%s, Output:%s, Key: %s}\n", string(cpath), string(cnPath), string(ckey))
	f, err := os.Open(string(cpath))
	strResult := "Success"
	if err != nil {
		fmt.Println("could not open file", string(cpath))
		strResult = err.Error()
	} else {
		err = decryptFile(ckey, f, string(cnPath))
		if err != nil {
			strResult = err.Error()
		}
	}
	return jni.Env(env).NewString(strResult)
}
