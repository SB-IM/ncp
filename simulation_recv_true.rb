require './lib/mqtt'
require './lib/help'
require 'yaml'

include Help

config = YAML.load_file('./config.yml')

mqtt = Mqtt.new config['server']['mqtt'], config['server']['id']

threads = []

threads << Thread.new do
  loop do
    topic, message = mqtt.cloud_get
    puts "Sub == #{topic} #{message}"

    msg = change_json(message)

    begin
      mqtt.cloud_put JSON.generate({
        jsonrpc: "2.0",
        result: msg,
        id: JSON.parse(msg)['id'] })

    rescue Exception => e
      puts e
      mqtt.cloud_put JSON.generate({
        jsonrpc: "2.0",
        error: e,
        id: JSON.parse(msg)['id'] })
    end
  end
end


puts "===== started ====="

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


