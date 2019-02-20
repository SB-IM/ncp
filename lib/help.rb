require './lib/chain'

module Help
  def is_json? str
    begin
      !!JSON.parse(str)
    rescue
      false
    end
  end

  def chain str, chain_lists
    bool, out_str = Chain.public_send(chain_lists.first, str)
    chain(out_str, chain_lists[1...]) if bool && chain_lists.length > 1
    out_str
  end
end

