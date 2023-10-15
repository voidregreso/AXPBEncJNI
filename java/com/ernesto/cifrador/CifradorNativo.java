package com.ernesto.cifrador;

public class CifradorNativo {
	
	static {
		System.loadLibrary("AXPBEncJNI");
	}
	
	public static native String encryptFile(String path, String newPath, String key);
	public static native String decryptFile(String path, String newPath, String key);
}
