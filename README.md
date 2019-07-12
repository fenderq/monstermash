# monstermash
A program to generate a fixed password list from a file and password.

## Password Details:
- 100bits of [password entropy](https://en.wikipedia.org/wiki/Password_strength) per password
- Ten passwords are generated for easy [storage on paper](https://www.schneier.com/news/archives/2010/11/bruce_schneier_write.html)
- Uses the [Base32](https://en.wikipedia.org/wiki/Base32) character set, which contains A-Z and 2-7
- Uses 200,000 rounds of [PBKDF2](https://en.wikipedia.org/wiki/PBKDF2) with a 128bit salt to produce Key and IV for [AES-256](https://en.wikipedia.org/wiki/Advanced_Encryption_Standard) cipher in [CBC mode](https://en.wikipedia.org/wiki/Block_cipher_mode_of_operation#Cipher_Block_Chaining_(CBC))

## Example Output:

```
01: MLMLI 4E3OM ZGQB3 7N4RD
02: Y3DRZ 7WUHN SWX7O X5BGD
03: SW6PL GWUKB QL2H3 O7IH6
04: EKAD6 LNN4I DWICK OEDE4
05: 34254 DT554 WYMKA VTPOR
06: LI54M 266CG SGNQP UASCC
07: IHGRF BIQQ4 ZGDUY FXZ2I
08: 6BQ7Q 25OVY QNW72 GLKWB
09: L22HG YU25A V44NV WJM4W
10: AGFFL KFIAB B43TO DUME2
```
