#!/bin/sh

cd ${0%/*}

if [ $# != 4 ]; then
  echo 'Usage: ./gen-site.sh <title> <description> <URL> <output-dir>'
  exit 1
fi

TITLE=${1}
DESCRIPTION=${2}
URL=${3}
OUTPUT_DIR=${4}

echo "Title: ${TITLE}"
echo "Description: ${DESCRIPTION}"
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

curl $IMG_URL -sLo $OUTPUT_DIR/site/pic.jpg

echo Done!

