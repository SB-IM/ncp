require 'json'
require 'rest-client'

class Ncp
  def initialize api_host, id, token='', retry_time=1
    @id = id
    @token = token

    @retry_time = retry_time.to_i

    api_ver = "/api/v1"
    @api_heartbeat = "#{api_ver}/nodes/#{id}/status_lives/0"
    @api_mission = "#{api_ver}/nodes/#{id}/mission_queues/"

    @server = RestClient::Resource.new(api_host)
  end

  def heartbeat payload=''
    connect_ncp :patch, @api_heartbeat, payload
  end

  def get_mission
    connect_ncp :get, @api_mission
  end

  def finish_mission id
    connect_ncp :delete, @api_mission + id.to_s
  end

  private
    def connect_ncp rest=:get, api=@api_heartbeat, payload=''
      begin
        JSON.parse(@server[api].send rest, { payload: payload })
      rescue Exception => e
        puts @api
        puts "Err: #{e}"
        sleep @retry_time

        retry
      end
    end
end
