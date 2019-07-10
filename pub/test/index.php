<?php

$ips = explode(',', getenv('MEMCACHE_SERVERS'));

$servers = [];
foreach ($ips as $ip) {
    $memcache = new Memcached();
    $memcache->addServer($ip, 11211);
    $servers[$ip] = $memcache;
}

foreach ($servers as $ip => $server) {
    $server->set('00', '00');
}

$results = [];
foreach ($servers as $ip => $server) {
    $results[$ip] = $server->get('00') === '00';
}

header('Content-Type: application/json;charset=utf8');
printf("%s\n", json_encode($results));