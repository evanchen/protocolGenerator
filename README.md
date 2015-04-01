# protocolGenerator

This is a simple tool,which aims to generate encoding/decoding APIs for user defined protocols to/from byte stream, mainly used for socket stream encoding/decoding.

1> The prototype is written in golang, so you have to install go to run the test.

2> For Windows, you can cd into protocolGenerator/ and run install.bat. It builds main.exe beneath /bin.

3> run main.exe -src [path of protocol files(test.proto,etc.)] -tar [path of generating go files(testProto.go,etc.)].

Any advices are appreciated! :)
