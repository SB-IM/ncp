
require './lib/ncp'
require './lib/http'
require './lib/mqtt'
require './lib/help'
require './lib/chain'
require 'yaml'
require 'socket'
require 'logger'

include Help

config = YAML.load_file('./config.yml')

$ncp = config['ncp']
#puts config

#log = Logger.new(STDOUT, level: :info)
log = Logger.new(config['env'] == "development" ? STDOUT : "log/#{config['env']}.log", level: config['log_level'])

$log = log

socket = TCPSocket.new config['ctl']['host'], config['ctl']['port']
sleep 1    # 这里延时连接的确认信息

# 冲掉连接确认信息缓冲区
#puts socket.recvmsg

mqtt = Mqtt.new config['mqtt'], config['id']

ncpc = RestHttp.new config['api_host'], config['id']

@payload = {
  link_id: ( config["link_id"] || 0 ),
  gps: {
    lat: "226876808",
    lon: "1142248069"
  }
}

@status = {}

threads = []

incoming_chain = 'change_json', 'filter_ncp'

threads << Thread.new do
  loop do
    topic, message = mqtt.cloud_get
    log.info "Sub == #{topic} #{message}"

    #puts chain(message, incoming_chain)
    bool, msg = chain(message, incoming_chain)
    #p bool
    #ncpc.finish_mission JSON.parse(msg)['id'] if JSON.parse(msg)['method'] == 'ncp'
    if bool
      socket.puts msg
    else
      mqtt.cloud_put msg
    end

  end
end

threads << Thread.new do
  loop do
    begin
      message = socket.gets.chomp

      if is_json_rpc? message
        #puts "#{message} is json"
        #log.info "Pub == #{message}"
        mqtt.cloud_put message
      else
        #puts "#{message} not json"
        #log.info "Pub == #{message}"
        mqtt.send_message message
      end

    rescue
      log.error socket.gets
      sleep 10
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

