# syntax=docker/dockerfile:experimental

## ENV var to set before/when running `docker build` command:
##   export DOCKER_BUILDKIT=1

########################################################################
###################   begin "full" image target   ######################
########################################################################
# start from the latest golang base image
FROM    golang:latest as full

# install build tools / package dependencies
RUN     apt update && apt install \
            -y \
            --no-install-recommends \
            --no-install-suggests \
            build-essential \
            ca-certificates \
            git \
            openssh-client \
            tree \
            vim

# download public SSH key for github.com
RUN     mkdir -p -m 0600 ~/.ssh/know_hosts \
            && ssh-keyscan github.com >> ~/.ssh/known_hosts
# ensure git commands use SSH instead of HTTPS
RUN     git config --global url."git@github.com:".insteadOf "https://github.com/"

# set environment vars
ENV     BASE_DIR    /app
ENV     CGO_ENABLED 0
ENV     GO111MODULE on
ENV     GOOS        linux
ENV     GOPATH      $BASE_DIR/go
ENV     GOPRIVATE   "github.com/Pii-Hole-Engineering/*"

# set the current working directory inside the container
WORKDIR $BASE_DIR

# copy the source fromt he current directory to the working directory
# inside the container
COPY    . .

# install build dependencies
RUN     --mount=type=ssh make vendor
## build the app / compile the no-phi-ai binary
RUN     --mount=type=ssh make build_only

# set the command to run the app
CMD     ["./build/no-phi-ai", "--help"]
########################################################################
####################   end "full" image target   #######################
########################################################################


########################################################################
###################   begin "mini" image target   ######################
########################################################################
# start a new stage from scratch, using the alpine:latest image as a
# base to build a very small "mini" (output) image containing only the
# compiled binary and any docs/
FROM    alpine:latest as mini

# install ca-certificates package
RUN     apk --no-cache add ca-certificates

# set the current working directory inside the container
WORKDIR /app

# copy the contents of the /app/docs directory from the "full" stage
COPY    --from=full /app/docs/ /app/docs/
# copy the compiled binary from the "full" image stage
COPY    --from=full /app/build/no-phi-ai /app/no-phi-ai

# command to run the app
CMD     ["./no-phi-ai"]
########################################################################
####################   end "mini" image target   #######################
########################################################################
