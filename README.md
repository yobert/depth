# depth
Adventures in depth of field

# credit
https://inconvergent.net/2019/depth-of-field/

# render

    go build
    ./depth

(wait a while)

    ffmpeg -r 60 -i frame%08d.png -c:v libvpx-vp9 -b:v 40M out.webm
