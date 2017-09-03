#!/bin/bash
./protoc --js_out=output_dir=../bin,library=myproto_libs,binary:. ./*.proto
read -p "Press any key to continue." var
