# Batman

Batman is a scriptable HTTP client that provides way to run single and multiple HTTP calls and chain them to scripts. It allows scripting using JavaScript, check examples for more details.

Examples:

* `BASE_URL=jsonplaceholder.typicode.com ITEM=todos ID=42 batman run examples/rest_collection.yml "Get todo"`
* `batman run examples/rest_collection.yml -e examples/.env "Get todo" "Create todo"`
* `batman script rest_script.yml`

## JavaScript API

Following JS API is available for scripting.

* `setEnv('key', value)` - Set variable
* `getEnv('key')` - Get variable
* `stop` - Stop executing script
* `assert(boolean_value)` - Assert
* `goto('step name')` - Go to step
