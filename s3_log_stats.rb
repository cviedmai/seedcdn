require 'set'
require 'time'
require 'geoip'

root = ARGV[0]

$all = Hash.new {|h,k| h[k] = {minutes: {}, hours: {}, total: [], files: Set.new, size: 0}}
$geo = GeoIP.new(root + 'GeoIP.dat')

def process(file)
  File.readlines(file).each do |line|
    parts = line.split(' ', 17)
    size = parts[14].to_i
    next if size == 0
    date_part = parts[2][1..-1]
    minute = Time.strptime(date_part, '%d/%b/%Y:%H:%M').to_i
    hour = Time.strptime(date_part, '%d/%b/%Y:%H').to_i

    ip = parts[4]
    key = $geo.country(ip).country_code2

    group = $all[key]
    group[:minutes][minute] = [] unless group[:minutes].include?(minute)
    group[:minutes][minute] << size
    group[:hours][hour] = [] unless group[:hours].include?(hour)
    group[:hours][hour] << size
    group[:total] << size

    url = parts[8]
    unless group[:files].include?(url)
      full_size = parts[15].to_i
      unless full_size == 0
        group[:files] << url
        group[:size] += full_size
      end
    end
  end
end

Dir[root + 'logs/*'].each do |f|
  process(f)
end

data = Hash.new {|h,k| h[k] = {minutes: [], hours: []}}
$all.each do |key, d|
  group = data[key]
  earliest = Time.now.to_i
  latest = 0
  d[:minutes].each do |time, sizes|
    if time < earliest
      earliest = time
    elsif time > latest
      latest = time
    end
    group[:minutes] << sizes.inject { |sum, x| sum += x }
  end
  group[:minutes].sort!

  d[:hours].each do |time, sizes|
    group[:hours] << sizes.inject { |sum, x| sum += x }
  end
  group[:hours].sort!

  group[:total] = d[:total].inject { |sum, x| sum += x } / ((latest - earliest) * 131072.0)
  group[:size] = d[:size] / 1073741824.0 #GB
end


data.sort{|a, b| b[1][:total] <=> a[1][:total]}.each do |key, d|
  key = key + ' ' * (16 - key.length) if key.length < 16
  puts "#{key} \t %10.2f \t %10.2f \t %10.2f \t %10.2f" % [d[:minutes][-1] / (60 * 131072.0), d[:hours][-1] / (3600 * 131072.0), d[:total], d[:size]]
end
