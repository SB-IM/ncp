require './lib/ncp'
require './lib/http'
require './lib/help'
require 'yaml'
require 'socket'
require 'logger'


include Help

config = YAML.load_file('./config.yml')


#log = Logger.new(STDOUT, level: :info)
log = Logger.new(config['env'] == "development" ? STDOUT : "log/#{config['env']}.log", level: config['log_level'])

ncp = NCP.new config['ncp'], config['server']

socket = TCPSocket.new *config['server']['tcps'].split(":")
sleep 1    # 这里延时连接的确认信息

# 冲掉连接确认信息缓冲区
#puts socket.recvmsg

loop do
  message = socket.gets.chomp

  msg = change_json(message)

  if JSON.parse(msg)['method'] == 'ncp'
    begin
      socket.puts JSON.generate({
        jsonrpc: "2.0",
        result: ncp.public_send(*JSON.parse(msg)['params']),
        id: JSON.parse(msg)['id'] })

    rescue Exception => e
      puts e
      socket.puts JSON.generate({
        jsonrpc: "2.0",
        error: e,
        id: JSON.parse(msg)['id'] })
    end
  else
    socket.puts msg
  end

end


