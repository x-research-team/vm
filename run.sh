#!/bin/sh
export GOPATH=$(pwd)

# format each go file
#echo "Formatting go file..."
#for file in `find ./src/kernel -name "*.go"`; do
#	echo "    `basename $file`"
#	go fmt $file > /dev/null
#done

echo ""

# run: ./kernel demo.ks or ./kernel
echo "Building REPL...(kernel)"
go build -o kernel main.go

echo "Building mdoc...(mdoc)"
go build -o mdoc mdoc.go

# run: ./fmt demo.ks
echo "Building Formatter...(fmt)"
go build -o fmt fmt.go

# run:    ./highlight demo.ks               (generate: demo.ks.html)
#     or  ./fmt demo.ks | ./highlight   (generate: output.html)
echo "Building Highlighter...(highlight)"
go build -o highlight highlight.go
