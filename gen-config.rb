#!/usr/bin/env ruby
require 'erb'
[1, 2, 3].each do |broker_id|
  port = 9092 + broker_id*100
  erb = ERB.new(File.read "./kafka-x.properties.erb")
  open("kafka-#{broker_id}.properties", "w") do |f|
    f.puts erb.result(binding)
  end
end
