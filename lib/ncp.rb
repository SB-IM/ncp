require 'json'
require 'rest-client'

class Ncp
  def initialize api_host, id, token='', retry_time=1
    @api_host = api_host
    @id = id
    @token = token

    @retry_time = retry_time.to_i

    api_ver = "/api/v1"
    @api_heartbeat = "#{api_ver}/nodes/#{id}/status_lives/"
    @api_mission = "#{api_ver}/nodes/#{id}/mission_queues/"
  end

  def heartbeat
    connect_ncp
  end

  def get_mission
    connect_ncp :get, @api_mission
  end

  def finish_mission id
    connect_ncp :delete, @api_mission + id.to_s
  end

  private
    def connect_ncp rest=:get, api=@api_heartbeat
      begin
        JSON.parse(RestClient.send rest, @api_host + api)
      rescue => e
        puts "Err: #{e}"
        sleep @retry_time

        retry
      end
    end
end
