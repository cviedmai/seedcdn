require 'set'
require 'time'
require 'geoip'

root = ARGV[0]

$all = Hash.new {|h,k| h[k] = {minutes: {}, hours: {}, total: 0, files: Set.new, size: 0, hits: 0}}
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
    group[:minutes][minute] = 0 unless group[:minutes].include?(minute)
    group[:minutes][minute] += size
    group[:hours][hour] = 0 unless group[:hours].include?(hour)
    group[:hours][hour] += size
    group[:total] += size
    group[:hits] += 1

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

$all.each do |key, group|
  sorted_by_time = group[:minutes].sort{|a, b| a[0] <=> b[0]}

  top = group[:minutes].sort{|a, b| b[1] <=> a[1]}.first
  group[:minutes] = top[1] / (60 * 131072.0)

  top = group[:hours].sort{|a, b| b[1] <=> a[1]}.first
  group[:hours] = top[1] / (3600 * 131072.0)

  earliest = sorted_by_time.first[0]
  latest = sorted_by_time.last[0]
  if earliest == latest
    group[:total] = 0
  else
    group[:total] = group[:total] / ((latest - earliest) * 131072.0)
  end
  group[:size] = group[:size] / 1073741824.0 #GB
end

$all.sort{|a, b| b[1][:total] <=> a[1][:total]}.each do |key, d|
  key = key + ' ' * (16 - key.length) if key.length < 16
  puts "#{key} \t %10.2f \t %10.2f \t %10.2f \t %10.2f \t %10d \t %10d" % [d[:minutes], d[:hours], d[:total], d[:size], d[:hits], d[:files].length]
end
