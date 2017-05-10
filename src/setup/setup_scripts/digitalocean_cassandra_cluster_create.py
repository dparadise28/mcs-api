import requests
import json
import sys

'''
current interface:

$ python setup/setup_scripts/digitalocean_cassandra_cluster_create.py setup/confs/digitalocean_create_instance_types.json cass_node
checking for required inputs:



        Enter the region you would like to deploy the server in (nyc3 suggested): nyc3

        Enter The desired ram of your server instance (https://www.digitalocean.com/pricing/#droplet; eg, 2gb): 1gb

        Please enter your api token: sample token
{
    "headers": {
        "Content-Type": "application/json",
        "Authorization": "Bearer sample token"
    },
    "data": {
        "region": "nyc3",
        "size": "1gb",
        "user_data": "#cloud-config\nruncmd:\n  - wget https://apt.puppetlabs.com/puppetlabs-release-trusty.deb",
        "image": "cassandra",
        "name": "cass"
    }
}
'''

def create_instance(instance_type, conf_path):
	with open(conf_path, 'rb') as do_conf:
		conf_str = do_conf.read()
		conf = json.loads(conf_str)[instance_type]
		print "checking for required inputs: \n\n"
		required_user_inputs = {
			req_input: raw_input(prompt)
				for
			req_input, prompt
				in
			conf['required_data'].iteritems()
		}
		
		formatted_request = json.dumps(conf)
		for k, v in required_user_inputs.iteritems():
			formatted_request = formatted_request.replace('{' + k + '}', v)
		formatted_request = json.loads(formatted_request)
		formatted_request['data']['user_data'] = '\n'.join(formatted_request['data']['user_data'])
		del formatted_request['required_data']
		
		print json.dumps(formatted_request, indent=4)
		confirmed = raw_input("Please confirm the following request is correct before attempting build (y/n): ")
		if confirmed == 'y':
			print 'attempting server build'
			response = requests.post(
				'https://api.digitalocean.com/v2/droplets',
				headers=formatted_request['headers'],
				data=json.dumps(formatted_request['data']),
			)
			print json.dumps(response.json(), indent=4)
		else:
			print "canceling server build"

if __name__ == "__main__":
	_, digital_ocean_conf_path, instance_type = sys.argv
	create_instance(instance_type, digital_ocean_conf_path)
