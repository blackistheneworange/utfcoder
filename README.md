# UTFCoder
An UTF transcoder utility built using golang

## Usage

Requires golang to run (or) build

```
go run . 
 -s "source file path" 
 -t "optional target file path"
 -from "one of utf-8/utf-16/utf-32" 
 -to "one of utf-8/utf-16/utf-16le/utf-16be/utf-32"
 -bom "boolean" (used to specify if output should have byte order mark added. false by default)
 ```
