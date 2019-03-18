require 'open-uri'
require 'yaml'

class NCP
  def initialize(config)
    @config=config
  end

  def download file, source
    File.open(@config[__method__.to_s][file], 'wb') {|f| f.write(open(source) {|f1| f1.read})}
  end

  def upload file, target
    File.open(target, 'wb') {|f| f.write(open(@config[__method__.to_s][file]) {|f1| f1.read})}
  end

  def status
    config = @config[__method__.to_s]
    pp config
  end
end
