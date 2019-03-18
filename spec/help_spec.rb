require './lib/help'
include Help

describe "#help" do
  it "not json method" do
    expect(JSON.parse(change_json 'aa bb')['method']).to eq('aa')
  end

  it "not json params" do
    expect(JSON.parse(change_json 'aa bb cc')['params']).to eq(['bb', 'cc'])
  end

  it "is json" do
    expect(is_json?('{"jsonrpc":"2.0","method":"aa","params":["bb"],"id":"1552893057"}')).to eq(true)
  end

  it "not json" do
    expect(is_json?('aa bb cc')).to eq(false)
  end

  it "is rpc" do
    expect(is_json_rpc?('{"jsonrpc":"2.0","method":"aa","params":["bb"],"id":"1"}')).to eq(true)
  end

  it "is rpc" do
    expect(is_json_rpc?('{"method":"aa","params":["bb"],"id":"1"}')).to eq(false)
  end

  it "not rpc" do
    expect(is_json?('aa bb cc')).to eq(false)
  end
end

