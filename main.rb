
require './lib/ncp'
require './lib/http'
require './lib/mqtt'
require 'yaml'
require 'socket'


config = YAML.load_file('./config.yml')
#puts config

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

####
# type : air
####

  #@payload = {
  #  link_id: 1,
  #  gps: {
  #    type: 4,
  #    satellites: 8,
  #    lat: "226876808",
  #    lon: "1142248069",
  #    height: "2.879999876022339"
  #  },
  #  battery: {
  #    remain: "99",
  #    voltage: "48369"
  #  },
  #  flight: {
  #    speed: "0.3416789770126343",
  #    time: "111",
  #    status: "3",
  #    mode: "0"
  #  }
  #}


@status = {}

threads = []

threads << Thread.new do
  loop do
    topic, message = mqtt.get_mission

    socket.puts JSON.generate({ method: message })
    puts socket.recvmsg
  end
end

threads << Thread.new do
  loop do
    @status = ncpc.heartbeat(@payload)
    #puts @status
    sleep @status['delay']
  end
end

threads << Thread.new do
  loop do
    response = ncpc.get_mission

    # 注： 双重判断是为了消除时间差而产生的误差
    if @status['has_msg?'] && response.length != 0

      if response[0]['name'] =~ /^ncp.*/
        puts "ncp cmd ==================="
        puts response[0]['name']
        #pp config['file']

        NCP.download response[0]['name'].split[3], config['file'][response[0]['name'].split[2]]
      else
        puts "send socket #{response[0]['name']}"

        socket.puts JSON.generate({ method: response[0]['name'] })

        puts "+++++++++++++"

        # not \n
        puts socket.recvmsg

        # have \n
        #puts socket.gets.chomp
      end

      ncpc.finish_mission response[0]['id']

      # 这个延时没什么意义，为了调试方便
      sleep 3

    end
  end
end

sleep 3

#pp @status

#thr.join

puts "===== started ====="

#thr.exit
#Thread.kill(thr)


#socket.puts "2333333333"
#socket.close

loop do end


#[:INT, :QUIT, :TERM].each do |sig|
#[:QUIT].each do |sig|
#  trap(sig) do
#    # clear pid file
#    puts "#{sig} signal received, exit!"
#  end
#end

