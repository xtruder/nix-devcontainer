set -xe

# check if hello is in path
hello

# check if XDG_DATA_DIRS is updated, so autocompletions will work
echo $XDG_DATA_DIRS | grep -i docker

# make sure extension preloader is running
ps aux | grep -i ext-preloader
