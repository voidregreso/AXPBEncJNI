package com.ernesto.cifrador;

public class Main {
    public static void main(String[] args) {
        if (args.length < 3) {
            System.out.println("Usage: java -jar Cifrador.jar [-d] <inputFilePath> <outputFilePath> <key>");
            System.exit(1);
        }

        boolean isDecrypt = args[0].equals("-d");
        String inputFilePath = args[isDecrypt ? 1 : 0];
        String outputFilePath = args[isDecrypt ? 2 : 1];
        String key = args[isDecrypt ? 3 : 2];

        if (isDecrypt) {
            String decrypted = CifradorNativo.decryptFile(inputFilePath, outputFilePath, key);
            System.out.println("Decryption result: " + decrypted);
        } else {
            String encrypted = CifradorNativo.encryptFile(inputFilePath, outputFilePath, key);
            System.out.println("Encryption result: " + encrypted);
        }
    }
}
