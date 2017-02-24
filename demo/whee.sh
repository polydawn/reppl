mkdir -p wares
mkdir -p debug
reppl init
reppl  put hash base  aLMH4qK1EdlPDavdhErOs0BPxqO0i6lUaeRE4DuUmnNMxhHtF56gkoeSulvwWNqT  --warehouse=http+ca://repeatr.s3.amazonaws.com/assets/
#reppl put hash go    sZOo52xMYaezehGNWm5c7W9bLNTIxGtybW_TAUzMeoKA--o2dtxusu5dYsTct2cV  --warehouse=https://storage.googleapis.com/golang/go1.7.1.linux-amd64.tar.gz
reppl  put hash go    UoY1amg4W8_JVQJ6tg6I4BQm1Mlw3ngT_kutZNr6XfFvvWAZfGrwDxDcQD2TzOVz  --warehouse=https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz
reppl eval step-A.frm
reppl eval step-B.frm
reppl eval step-G.frm
reppl unpack hellogopher debug/hellogopher
