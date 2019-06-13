#!/bin/sh
#ffmpeg -r 60 -i frame%08d.png -c:v libvpx-vp9 -b:v 40M out.webm
#ffmpeg -r 60 -i frame%08d.png -c:v libvpx-vp9 -b:v 2M out.webm
ffmpeg -r 60 -i frame%08d.png out.webm
