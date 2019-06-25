require 'rest-client'
require 'api-auth'

class RestHttp
  def initialize id, secret_key=''
    @id = id.to_s
    @secret_key = secret_key
  end

  def request params = {}
    #ApiAuth.sign!(RestClient::Request.new(params), @id, @secret_key).execute.body
    tt = ApiAuth.sign!(RestClient::Request.new(params), @id, @secret_key)
    pp tt.headers
    tt.execute.body
  end
end
