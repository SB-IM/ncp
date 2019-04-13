require './lib/http'

class NCP
  def initialize config, server
    @config = config
    @server = RestHttp.new server["id"], server["secret_key"]

    if config['shell']
      Dir.glob("#{config['shell']['path']}*").each do |item|
        system "#{config['shell']['prefix']}#{item}" if item =~ /^#{config['shell']['path']}_init_.*#{config['shell']['suffix']}$/
      end
    end
  end

  def method_missing(method, *args)
    "The method #{method} with you call not exists on NCP, params: #{args.join(' ')}"
  end

  def download file, source
    File.open(@config[__method__.to_s][file], 'wb') do |f|
      f.write(@server.request(url: source, method: :get) {|f1| f1.read})
    end
  end

  def upload file, target
    @server.request(
      url: target,
      method: :patch,
      payload: {
        file => File.open(@config[__method__.to_s][file], 'r')
      }
    )
  end

  def status
    config = @config[__method__.to_s]
    pp config
  end

  def shell cmd
    config = @config[__method__.to_s]
    config ? system(config['prefix'] + config['path'] + cmd + config["suffix"]) : "Disable shell"
  end
end
