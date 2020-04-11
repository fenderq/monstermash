# monstermash
A program to generate a fixed password list from a file and master password.

## Technical Details:
- 100bits of [entropy](https://en.wikipedia.org/wiki/Password_strength) per password generated
- 128bit salt produced and used from SHA256 of file
- 200,000 rounds of [PBKDF2](https://en.wikipedia.org/wiki/PBKDF2) to derive Key and IV for cipher
- [AES-256](https://en.wikipedia.org/wiki/Advanced_Encryption_Standard) cipher in [CTR mode](https://en.wikipedia.org/wiki/Block_cipher_mode_of_operation#Counter_(CTR)) used as a [CSPRNG](https://en.wikipedia.org/wiki/Cryptographically_secure_pseudorandom_number_generator)
- Debug mode command line argument _-d_ to inspect data
- Passwords consist of the [Base32](https://en.wikipedia.org/wiki/Base32) character set (A-Z 2-7)
- Regression test functionality with _go test_
- Ten passwords are generated for easy [storage on paper](https://www.schneier.com/news/archives/2010/11/bruce_schneier_write.html) by default

## Example Usage:

```
$ monstermash myfile
enter password:
confirm password:
01: YLNF5 GHLMK 53PTF VLH2D
02: 7YWOF QUYSJ BGTBW RWZQP
03: 6HR6Y LSIGF 46SZH K3Z7W
04: BEAQ6 Z3J4I LXIE2 TKVFN
05: 7GMHI FQ6ID N564P FB3GD
06: Y5HSO 3L2S7 ZCOR7 REWJ2
07: QV3OS Y5AAJ BVBV5 FWOCB
08: S2PV6 LHPBY AIWSA 5NE6O
09: RCFPD VEYNN C4ZDB 5RVW3
10: HO6SU H4S3I KG466 WPWJT
```
## More Information

You can find more information on the _monstermash_ [wiki page](https://github.com/fenderq/monstermash/wiki).
