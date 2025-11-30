# UTFCoder
An UTF transcoder command line utility built using golang

[Latest Release](https://github.com/blackistheneworange/utfcoder/releases/tag/v0.1.0)

## Usage

```
utfcoder
 -s "source file path" 
 -t "optional target file path"
 -from "one of utf-8/utf-16/utf-32" 
 -to "one of utf-8/utf-16/utf-16le/utf-16be/utf-32"
 -bom "boolean" (used to specify if output should have byte order mark added. false by default.)
 -verbose "boolean" (used to print logs for debugging. false by default.)
 ```
