# Node control protocol

#### Depend
- libgstreamer1.0-dev
- libgstreamer-plugins-base1.0-dev
- gstreamer1.0-plugins-good

For debian
```sh
apt-get install -y libgstreamer1.0-dev libgstreamer-plugins-base1.0-dev gstreamer1.0-plugins-good
```

#### Compile && Run
go >= 1.12.4

```sh
git submodule update --init --recursive
make dep
make build

# run
cp conf/config-dist.yml config.yml
ncp -c config.yml
```

