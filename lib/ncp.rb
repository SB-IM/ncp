require 'json'
#require 'open-uri'
require 'rest-client'

#module Ncp
class Ncp
#  attr_accessor :api_host, :id, :token

  def initialize api_host, id, token=''
    @api_host = api_host
    @id = id
    @token = token

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
      end
    end

end
