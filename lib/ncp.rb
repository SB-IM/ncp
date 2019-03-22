require './lib/http'

class NCP
  def initialize config, server
    @config = config
    @server = RestHttp.new server["id"], server["secret_key"]
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
end
