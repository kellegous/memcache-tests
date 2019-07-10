import json
import os
import subprocess

class Node(object):
    def __init__(self, data):
        self.data = data

    def id(self):
        return self.data.get('Id')

    def labels(self):
        cfg = self.data.get('Config', {})
        return cfg.get('Labels', {})

    def role(self):
        return self.labels().get('role')

    def app(self):
        return self.labels().get('app')

    def name(self):
        name = self.data.get('Name')
        if len(name) > 0 and name[0] == '/':
            name = name[1:]
        return name
    
    def ip_address(self):
        network = self.data.get('NetworkSettings', {})
        return network.get('IPAddress')

    def stop(self):
        with open('/dev/null', 'w') as w:
            return subprocess.call(
                ['docker', 'stop', self.id()],
                stdout=w) == 0

def LoadNode(id):
    p = subprocess.Popen([
        'docker', 'inspect', id, '--format={{json .}}'
    ], stdout=subprocess.PIPE)
    out, _ = p.communicate()
    return Node(json.loads(out))

def LoadNodes():
    p = subprocess.Popen([
        'docker', 'ps', '-q', '--filter=label=app=memcache'
    ], stdout=subprocess.PIPE)
    out, _ = p.communicate()
    return [LoadNode(x) for x in out.split("\n") if x.strip() != ""]

def StartMemcached(name):
    p = subprocess.Popen([
        'docker', 'run', '-ti', '--rm', '-d',
        '--label=app=memcache',
        '--label=role=memcached',
        '--name={}'.format(name),
        'memcached'
    ], stdout=subprocess.PIPE)
    out, _ = p.communicate()
    return out.strip()

def StartApp(name, root):
    p = subprocess.Popen([
        'docker', 'run', '--rm', '-d',
        '--label=app=memcache',
        '--label=role=app',
        '--name={}'.format(name),
        '-v', '{}:/app'.format(os.path.join(root, 'pub')),
        '-v', '{}:/vendor'.format(os.path.join(root, 'vendor')),
        '-v', '{}:/templates'.format(os.path.join(root, 'templates')),
        '-p', '8080:80',
        'webdevops/php-nginx'
    ], stdout=subprocess.PIPE)