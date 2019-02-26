require './lib/chain'

module Help
  def is_json? str
    begin
      !!JSON.parse(str)
    rescue
      false
    end
  end

  def is_json_rpc? str
    if is_json? str
      JSON.parse(str).has_key? "jsonrpc"
    else
      false
    end
  end

  def chain str, chain_lists
    bool, out_str = Chain.public_send(chain_lists.first, str)
    bool, out_str = chain(out_str, chain_lists[1..-1]) if bool && chain_lists.length > 1
    return bool, out_str
  end
end

