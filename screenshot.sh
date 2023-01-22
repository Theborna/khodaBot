#!/bin/sh
./wkhtmltoimage.exe --load-error-handling ignore --enable-javascript --javascript-delay 3000 $1 $2