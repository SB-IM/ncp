# About the environment

#### Compile ncpgo
go >= 1.12.4

```sh
go build
```

#### ruby
Compile and install ruby （Not rbenv and rvm）

Need ruby 2.5.0 +

Suggest ruby 2.6.0 +

```sh
wget https://cache.ruby-lang.org/pub/ruby/2.6/ruby-2.6.1.tar.xz

# SHA256: 47b629808e9fd44ce1f760cdf3ed14875fc9b19d4f334e82e2cf25cb2898f2f2

./configure --prefix=$HOME/.ruby

make

# Not sudo
make install

echo 'export PATH="$HOME/.ruby/bin:$PATH"' >> $HOME/.bashrc
```

ruby 2.6.x Default bundler 1.17.x

If Need bundler 2.0.x
```sh
gem install bundler

exit
```

#### Run
```sh
ncpgo

bundler exec ruby ncp.rb
```

