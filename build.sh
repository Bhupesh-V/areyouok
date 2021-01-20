#!/usr/bin/env bash

package_name="areyouok"
platforms=("windows/amd64" "linux/amd64" "linux/arm64" "linux/386" "darwin/amd64")

# remove any old binaries
echo -e "Removing previous builds"
rm "$package_name"-*

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    printf "%s " "Building for $GOOS-$GOARCH"
    if export GOOS="$GOOS" GOARCH="$GOARCH"; then
        output_name="$package_name-$GOOS-$GOARCH$(go env GOEXE)"
        go build -ldflags="-X 'main.aroVersion=$1' -X 'main.aroDate=$(date '+(%d %b %Y)')'" -o "$output_name"
        echo -e "✅"
    else
        echo -e "An error has occurred! Aborting ..."
        exit 1
    fi
done
