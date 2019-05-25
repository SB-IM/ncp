#!/bin/sh


./ncpgo &


#export PATH=/var/lib/dockctl/.rbenv/shims/:$PATH
PATH="$HOME/.ruby/bin:$PATH"
#export PATH=/var/lib/dockctl/.rbenv/bin:$PATH
#eval "$(rbenv init -)"

sleep 5

ruby -v


cd /var/lib/dockctl/ncp

bundle exec ruby main.rb

