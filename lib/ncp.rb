require 'open-uri'

module NCP
  def self.download source, target
    File.open(target, 'wb') {|f| f.write(open(source) {|f1| f1.read})}
  end
end
