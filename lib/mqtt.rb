require 'mqtt'

class Mqtt
  def initialize uri, id, token=''
    @id = id
    #@token = token

    @api_heartbeat = "nodes/#{id}/message"
    @api_status = "nodes/#{id}/status"

    @status_map = {
      online: 0,
      offline: 1,
      neterror: 2
    }

    @server = MQTT::Client.connect(
      uri,
      :will_topic => @api_status,
      :will_payload => @status_map[:neterror].to_s,
      :will_qos => 1,
      :will_retain => true)

    @server.publish(@api_status, @status_map[:online].to_s, retain=true, qos=1)
  end

  def heartbeat payload=''
    @server.publish(@api_heartbeat, payload.to_json)
  end

  def get_mission
    @server.get("nodes/#{@id}/ctl")
  end

  def send_message payload=''
    $log.info "Pub == #{@api_heartbeat} #{payload}"
    @server.publish(@api_heartbeat, payload)
  end

end
