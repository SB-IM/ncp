require 'mqtt'

class Mqtt
  def initialize uri, id
    #, token=''
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
    @server.subscribe("nodes/#{@id}/rpc/send" => 1)
  end

  def offline
    @server.publish(@api_status, @status_map[:offline].to_s, retain=true, qos=1)
    @server.disconnect
  end

  def heartbeat payload=''
    @server.publish(@api_heartbeat, payload.to_json)
  end

  def cloud_get
    #@server.get("nodes/#{@id}/rpc/send" => 1)
    @server.get
  end

  def cloud_put payload=''
    @server.publish("nodes/#{@id}/rpc/recv", payload, retain=true, qos=1)
  end

  def send_message payload=''
    $log.info "Pub == #{@api_heartbeat} #{payload}"
    @server.publish(@api_heartbeat, payload, retain=true, qos=1)
  end

end
