require './lib/ncp'
require 'yaml'
require 'socket'


config = YAML.load_file('./config.yml')

puts config

socket = TCPSocket.new config['ctl']['hostname'], config['ctl']['port']
sleep 1    # 这里延时连接的确认信息

# 冲掉连接确认信息缓冲区
puts socket.recvmsg


ncpc = Ncp.new config['api_host'], config['id']

thr = Thread.new do
  while true do
    @status = ncpc.heartbeat
    #puts @status
    sleep @status['delay']
  end
end

# 这个值是考虑网络延时。。
sleep 3

pp @status


#run = []

while true do
#
#  status = ncpc.heartbeat
  if @status['has_msg?']
    response = ncpc.get_mission
    #pp response
    #if run
    #  run << response
      #socket.sendmsg response[0]['name']
      socket.puts response[0]['name']
    #end
    puts "+++++++++++++"
    #puts socket.recvmsg
    puts socket.gets
    #ncpc.finish_mission response[0]['id']

    # 这个延时没什么意义，为了调试方便
    sleep 3

  end
end

