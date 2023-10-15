// Este archivo se genera dinamicamente, asi que no lo modifique.

#include <stdlib.h>
#include <stdio.h>
#include <stddef.h>
#include <stdint.h>
#include <jni.h>

#include "_cgo_export.h"

extern void jniOnLoad(GoUintptr vm);
extern void jniOnUnload(GoUintptr vm);

extern GoUintptr jni_com_ernesto_cifrador_encryptFile1(GoUintptr env, GoUintptr clazz, GoUintptr path, GoUintptr newPath, GoUintptr key);
extern GoUintptr jni_com_ernesto_cifrador_decryptFile2(GoUintptr env, GoUintptr clazz, GoUintptr path, GoUintptr newPath, GoUintptr key);

jint JNI_OnLoad(JavaVM *vm, void *reserved) {
    JNIEnv *env = NULL;
    if ((*vm)->GetEnv(vm, (void **) &env, JNI_VERSION_1_6) != JNI_OK) {
        fprintf(stderr, "[%s:%d] GetEnv() return error\n", __FILE__, __LINE__);
        abort();
    }

    jclass clazz;
    JNINativeMethod methods[255];
    jint size;
    char *name;

    name = "com/ernesto/cifrador/CifradorNativo";
    clazz = (*env)->FindClass(env, name);
    size = 0;

    methods[size].fnPtr = jni_com_ernesto_cifrador_encryptFile1;
    methods[size].name = "encryptFile";
    methods[size].signature = "(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)Ljava/lang/String;";
    size++;

    methods[size].fnPtr = jni_com_ernesto_cifrador_decryptFile2;
    methods[size].name = "decryptFile";
    methods[size].signature = "(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)Ljava/lang/String;";
    size++;

    if ((*env)->RegisterNatives(env, clazz, methods, size) != 0) {
        fprintf(stderr, "[%s:%d] %s RegisterNatives() return error\n", __FILE__, __LINE__, name);
        abort();
    }

    jniOnLoad((GoUintptr) vm);
    return JNI_VERSION_1_6;
}

void JNI_OnUnload(JavaVM *vm, void *reserved) {
    jniOnUnload((GoUintptr) vm);
}