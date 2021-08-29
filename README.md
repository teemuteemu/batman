# Batman

Batman > Postman

Batman is a scriptable HTTP client that provides way to store and run single and multiple HTTP calls, and chain them to scripts. It allows scripting using JavaScript. Currently it focuses to REST APIs, so it might not cover all the alternative ways of dealing with HTTP APIs.

Batman can:
 * Store all your API calls into a single YAML file
 * Execute 0-n of the calls from the provided collections
 * Use environment variables (or `.env` files) for templating the requests (URLs, headers, body)
 * Run scripts over multiple of API calls to execute complex call flows
 * Format the calls to `curl` commands

## Examples:

### Run a single request using env vars
`BASE_URL=jsonplaceholder.typicode.com ITEM=todos ID=42 batman run examples/rest_collection.yml "Get todo"`

### Run multiple call using .env file
`batman run examples/rest_collection.yml -e examples/.env "Get todo" "Create todo"`

### Run a scipt that uses a collection
`batman script rest_script.yml`

### Output a call as a curl command
`batman print rest_collection.yml examples/rest_collection.yml "Get todo"`

## JavaScript API

Following JS API is available for scripting.

* `setEnv('key', value)` - Set variable
* `getEnv('key')` - Get variable
* `stop` - Stop executing script
* `assert(boolean_value)` - Assert
* `goto('step name')` - Go to step

For more details check the `examples` directory.
