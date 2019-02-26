### About the environment

Compile and install ruby （Not rbenv and rvm）

Need ruby 2.5.0 +

Suggest ruby 2.6.0 +

```
wget https://cache.ruby-lang.org/pub/ruby/2.6/ruby-2.6.1.tar.xz

./configure --prefix=$HOME/.ruby

make

# Not sudo
make install

echo 'export PATH="$HOME/.ruby/bin:$PATH"' >> $HOME/.bashrc
```

ruby 2.6.x Default bundler 1.17.x

If Need bundler 2.0.x
```
gem install bundler

exit
```

#### Run
```
bundler exec ruby main.rb
```

