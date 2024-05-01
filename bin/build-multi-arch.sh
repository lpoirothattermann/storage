#!/usr/bin/env bash

platforms=("darwin/amd64" "darwin/arm64" "linux/arm64" "linux/amd64")
basedir="build"

for platform in "${platforms[@]}"
do
	platform_split=(${platform//\// })
	GOOS=${platform_split[0]}
	GOARCH=${platform_split[1]}
	output_dir=${basedir}/${GOOS}-${GOARCH}
	output_name=${package_name}
	if [ $GOOS = "windows" ]; then
		output_name+='.exe'
	fi

	env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_dir/$output_name main.go
	if [ $? -ne 0 ]; then
		echo 'An error has occurred! Aborting the script execution...'
		exit 1
	fi
  echo "${platform} has been built"
done

echo -e "\nEverything has been built in ${basedir}/"
