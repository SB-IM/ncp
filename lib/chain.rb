require './lib/ncp'

module Chain
  def self.change_json str
    begin
      !!JSON.parse(str)

      return true, str
    rescue
      return true, JSON.generate({ jsonrpc: "2.0", method: str.split.first, params: str.split[1..-1], id: "0" })
    end
  end

  def self.filter_ncp str
    if JSON.parse(str)['method'] == 'ncp'
      #p JSON.parse(str)['params']
      #pp *JSON.parse(str)['params'][1..-1]

      # NCP...
      NCP.public_send JSON.parse(str)['params'].first, *JSON.parse(str)['params'][1..-1]
      #NCP.public_send(JSON.parse(str)['params'].first) response[0]['name'].split[3], config['file'][response[0]['name'].split[2]]
      return false, str
    else
      return true, str
    end
  end
end
