require 'open-uri'
require 'yaml'


module NCP
  #def self.download source, target
  #  File.open(target, 'wb') {|f| f.write(open(source) {|f1| f1.read})}
  #end

  def self.download file, source
    #config = YAML.load_file('./config.yml')['ncp']
    #pp $ncp
    File.open($ncp['file'][file], 'wb') {|f| f.write(open(source) {|f1| f1.read})}
  end
end
