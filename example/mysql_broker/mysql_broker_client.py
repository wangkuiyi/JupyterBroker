import urllib

f = {'url' : 'select * from mysql.time_zone_name limit 10;'}
print(urllib.urlencode(f))

