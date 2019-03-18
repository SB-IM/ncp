require './lib/ncp'

module Chain
  def self.change_json str
    begin
      !!JSON.parse(str)

      return true, str
    rescue
      #return true, JSON.generate({ jsonrpc: "2.0", method: str.split.first, params: str.split[1..-1], id: "0" })
      return true, JSON.generate({ jsonrpc: "2.0", method: str.split.first, params: str.split[1..-1], id: Time.now.to_i.to_s })
    end
  end

  def self.filter_ncp str
    if JSON.parse(str)['method'] == 'ncp'
      #p JSON.parse(str)['params']
      #pp *JSON.parse(str)['params'][1..-1]

      $log.info "Ncp == #{JSON.parse(str)['params']} #{file} #{target}"

      result = NCP.public_send JSON.parse(str)['params'].first, *JSON.parse(str)['params'][1..-1]

      return false, JSON.generate({ jsonrpc: "2.0", result: result, id: JSON.parse(str)['id'] })
    else
      return true, str
    end
  end
end
