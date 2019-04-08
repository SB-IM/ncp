
require './lib/ncp'
require './lib/http'
require './lib/mqtt'
require './lib/help'
require 'yaml'
require 'socket'
require 'logger'

include Help

config = YAML.load_file('./config.yml')


#log = Logger.new(STDOUT, level: :info)
log = Logger.new(config['env'] == "development" ? STDOUT : "log/#{config['env']}.log", level: config['log_level'])

ncp = NCP.new config['ncp'], config['server']

socket = TCPSocket.new config['ctl']['host'], config['ctl']['port']
sleep 1    # 这里延时连接的确认信息

# 冲掉连接确认信息缓冲区
#puts socket.recvmsg

#http = RestHttp.new

mqtt = Mqtt.new config['mqtt'], config['server']['id']

threads = []

threads << Thread.new do
  loop do
    topic, message = mqtt.cloud_get
    log.info "Sub == #{topic} #{message}"

    msg = change_json(message)

    if JSON.parse(msg)['method'] == 'ncp'
      begin
        mqtt.cloud_put JSON.generate({
          jsonrpc: "2.0",
          result: ncp.public_send(*JSON.parse(msg)['params']),
          id: JSON.parse(msg)['id'] })

      rescue Exception => e
        puts e
        mqtt.cloud_put JSON.generate({
          jsonrpc: "2.0",
          error: e,
          id: JSON.parse(msg)['id'] })
      end
    else
      socket.puts msg
    end

  end
end

threads << Thread.new do
  loop do
    begin
      message = socket.gets.chomp
      log.info "Pub == #{message}"

      if is_json_rpc? message
        #puts "#{message} is json"
        mqtt.cloud_put message
      else
        #puts "#{message} not json"
        mqtt.send_message message
      end

    rescue Exception => e
      #log.error socket.gets
      log.error e.message
      log.error e.backtrace.inspect
      begin
        log.warn "Try Reconnect socket"
        socket = TCPSocket.new config['ctl']['host'], config['ctl']['port']
      rescue
        log.error "Reconnect failed for 10s"
        sleep 10
      end
    end
  end
end

log.warn "===== started ====="

[:INT, :QUIT, :TERM].each do |sig|
#[:QUIT].each do |sig|
  trap(sig) do
    # clear pid file
    puts "#{sig} signal received, exit!"

    threads.each { |thr| thr.exit }
    socket.close
    puts socket.inspect
  end
end

threads.each { |thr| thr.join }
mqtt.offline

log.warn "===== stoped  ====="

