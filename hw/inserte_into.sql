insert into developers(id, name, department, geolocation, last_known_ip, is_available)
select	uuid_generate_v4(),
	    (array['James', 'Mary', 'John', 'Patricia', 'Robert'])[floor(random() * 5+1)] || ' ' || (array['Smith', 'Johnson', 'Williams', 'Brown', 'Jones'])[floor(random() * 4+1)], 
	    (array['backend', 'android', 'ios', 'frontend'])[floor(random() * 4+1)],
	    st_setsrid(st_makepoint(54.713505 + random() * 5, 20.454683 + random() * 5), 4326), 
	    (floor(random() * 255)::text || '.' || floor(random() * 255)::text || '.' || floor(random() * 255)::text ||'.' || floor(random() * 255)::text)::cidr,
	    (array[true, false])[floor(random() * 2+1)]
from generate_series(1, 100000)