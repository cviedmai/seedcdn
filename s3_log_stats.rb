require 'time'
require 'geoip'

root = ARGV[0]

$all = Hash.new {|h,k| h[k] = {minutes: {}, hours: {}, total: []}}
$geo = GeoIP.new(root + 'GeoIP.dat')

def process(file)
  File.readlines(file).each do |line|
    parts = line.split(' ', 16)
    size = parts[14].to_i
    next if size == 0
    date_part = parts[2][1..-1]
    minute = Time.strptime(date_part, '%d/%b/%Y:%H:%M').to_i
    hour = Time.strptime(date_part, '%d/%b/%Y:%H').to_i

    ip = parts[4]
    url = parts[8]

    key = $geo.country(ip).country_code2

    $all[key][:minutes][minute] = [] unless $all[key][:minutes].include?(minute)
    $all[key][:minutes][minute] << size

    $all[key][:hours][hour] = [] unless $all[key][:hours].include?(hour)
    $all[key][:hours][hour] << size

    $all[key][:total] << size
  end
end

Dir[root + 'logs/*'].each do |f|
  process(f)
end

data = Hash.new {|h,k| h[k] = {minutes: [], hours: []}}
$all.each do |key, d|
  earliest = Time.now.to_i
  latest = 0
  d[:minutes].each do |time, sizes|
    if time < earliest
      earliest = time
    elsif time > latest
      latest = time
    end
    data[key][:minutes] << sizes.inject { |sum, x| sum += x }
  end

  data[key][:minutes].sort!

  d[:hours].each do |time, sizes|
    data[key][:hours] << sizes.inject { |sum, x| sum += x }
  end
  data[key][:hours].sort!

  data[key][:total] = d[:total].inject { |sum, x| sum += x } / ((latest - earliest) * 131072.0)
end


data.sort{|a, b| b[1][:total] <=> a[1][:total]}.each do |key, d|
  key = key + ' ' * (16 - key.length) if key.length < 16
  puts "#{key} \t %10.2f \t %10.2f \t %10.2f" % [d[:minutes][-1] / (60 * 131072.0), d[:hours][-1] / (3600 * 131072.0), d[:total]]
end
