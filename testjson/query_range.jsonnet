local query = std.extVar('QUERY');
local start = std.extVar('START');
local end = std.extVar('END');
local step = std.extVar('STEP');

// Outputs JSON-RPC params.arguments
// Used by test.http.sh and test.http.jsonnet
{
    query: query,
    start: start,
    end: end,
    step: step,
}
