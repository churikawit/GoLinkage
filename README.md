
# GoLinkage
Connect to Linkage Center system for Thai ID card by wrapping ami32.dll and scapi_ope.dll
and provide web service at client side

## Requirement
These required files are distributed separately by DOPA (Department Of Provincial Administration)
- ami32.dll
- scapi_ope.dll
- scapi_ope.dli
- lm.exe

## Build
$ go build -ldflags "-H=windowsgui" -o GoLinkage.exe

## Install
Make program startup automatically 

- Put program into "C:\ProgramData\Microsoft\Windows\Start Menu\Programs\StartUp"

- Install as a Windows Service (not support VerifyPin)
  - $ sc.exe create "GoLinkage" binpath= "c:\GoLinkage.exe" DisplayName= "GoLinkage" start= auto
  - $ sc.exe description "GoLinkage" "Web service for smartcard connection and LinkageCenter connection"

