require './lib/ncp'

module Chain
  def self.change_json str
    begin
      !!JSON.parse(str)

      return true, str
    rescue
      return true, JSON.generate({ method: str.split.first, params: str.split[1..] })
    end
  end

  def self.filter_ncp str
    puts JSON.parse(str)['method']
    if JSON.parse(str)['method'] == 'ncp'
      puts "NCP" + str

      # NCP...
      # NCP.download response[0]['name'].split[3], config['file'][response[0]['name'].split[2]]
      return false, str
    else
      return true, str
    end
  end
end
