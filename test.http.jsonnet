// Constructs a JSON-RPC request for a tool method with optional parameters.
// Expects the following external variables:
// - METHOD: The tool method (e.g., "list", "call")
// - NAME: The name of the tool (optional: required if METHOD=="call")
// - ARGUMENTS: A JSON string representing the arguments for the tool (optional)
// Output
{
  local method = std.extVar('METHOD'),
  local name = std.extVar('NAME'),

  jsonrpc: '2.0',
  id: 1,
  method: 'tools/%(method)s' % { method: method },
  // If the method is "list" or the name is empty (""), params is an empty object ({})
  params: if method == 'list' || name == '' then {} else {
    // If the arguments is empty (""), arguments is an empty object ({})
    // Otherwise parse the arguments as JSON
    local arguments =
      local raw_args = std.extVar('ARGUMENTS');
      if raw_args == '' then {} else std.parseJson(raw_args),

    name: name,
    arguments: arguments,
  },
}