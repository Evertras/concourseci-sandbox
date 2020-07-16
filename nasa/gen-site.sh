#!/bin/sh -e

cd ${0%/*}

if [ $# != 5 ]; then
  echo 'Usage: ./gen-site.sh <title> <description> <pubDate> <URL> <output-dir>'
  exit 1
fi

TITLE=${1}
DESCRIPTION=${2}
PUBDATE=${3}
URL=${4}
OUTPUT_DIR=${5}

# This doesn't work on mac because BSD date is sadness
SHORTDATE=$(date -d "${PUBDATE}" '+%Y-%m-%d')

echo "Title: ${TITLE}"
echo "Description: ${DESCRIPTION}"
echo "PubDate: ${PUBDATE}"
echo "URL: ${URL}"
echo "Output directory: ${OUTPUT_DIR}"
echo "Short Date: ${SHORTDATE}"

cp Dockerfile ${OUTPUT_DIR}/

sed "s/{{title}}/${SHORTDATE} - ${TITLE}/" ./site/index.html \
  | sed "s/{{description}}/${DESCRIPTION}/" \
  > ${OUTPUT_DIR}/site/index.html

IMG_URL=$(curl -sL ${URL} \
  | grep 'meta property="og:image" content="' \
  | sed -n 's/.*content="\(.*\)".*/\1/p')

echo Targeting ${IMG_URL}

curl ${IMG_URL} -sLo ${OUTPUT_DIR}/site/pic.jpg

echo ${SHORTDATE} > ${OUTPUT_DIR}/tag

echo Done!

