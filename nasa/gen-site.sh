#!/bin/sh -e

cd ${0%/*}

if [ $# != 5 ]; then
  echo 'Usage: ./gen-site.sh <title> <description> <pubDate> <URL> <output-dir>'
  exit 1
fi

TITLE=${1}
DESCRIPTION=${2}
PUDDATE=${3}
URL=${4}
OUTPUT_DIR=${5}

echo "Title: ${TITLE}"
echo "Description: ${DESCRIPTION}"
echo "PubDate: ${PUBDATE}"
echo "URL: ${URL}"
echo "Output directory: ${OUTPUT_DIR}"

cp Dockerfile ${OUTPUT_DIR}/

sed "s/{{title}}/${TITLE}/" ./site/index.html \
  | sed "s/{{description}}/${DESCRIPTION}/" \
  > ${OUTPUT_DIR}/site/index.html

IMG_URL=$(curl -sL ${URL} \
  | grep 'meta property="og:image" content="' \
  | sed -n 's/.*content="\(.*\)".*/\1/p')

echo Targeting ${IMG_URL}

curl ${IMG_URL} -sLo ${OUTPUT_DIR}/site/pic.jpg

TAG=$(date -d "${DATE}" '+%Y-%m-%d')

echo ${TAG} > ${OUTPUT_DIR}/tag

echo Done!

