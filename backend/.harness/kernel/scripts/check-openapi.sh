#!/usr/bin/env bash
set -euo pipefail

SPEC=""
for f in internal/adapter/rest/swagger/openapi.yaml internal/adapter/rest/docs/openapi.yaml docs/openapi.yaml openapi.yaml; do
  if [ -f "$f" ]; then SPEC="$f"; break; fi
done

if [ -z "$SPEC" ]; then
  echo "No OpenAPI spec found."
  exit 1
fi

echo "Found OpenAPI spec: $SPEC"
python3 - <<'PY' "$SPEC"
import pathlib
import re
import sys

try:
    import yaml
except Exception:
    raise SystemExit('PyYAML is required for the OpenAPI gate.')

path = pathlib.Path(sys.argv[1])
with path.open('r', encoding='utf-8') as f:
    data = yaml.safe_load(f)

if not isinstance(data, dict):
    raise SystemExit('OpenAPI YAML is not an object')
if data.get('openapi') != '3.1.0':
    raise SystemExit('openapi must be 3.1.0')
if data.get('info', {}).get('version') != '1.5.0':
    raise SystemExit('info.version must be 1.5.0')
if data.get('security') != [{'cookieAuth': []}]:
    raise SystemExit('cookieAuth must be the default security requirement')

paths = data.get('paths')
schemas = data.get('components', {}).get('schemas')
security_schemes = data.get('components', {}).get('securitySchemes', {})
if not isinstance(paths, dict) or not paths:
    raise SystemExit('paths must be a non-empty object')
if not isinstance(schemas, dict) or not schemas:
    raise SystemExit('components.schemas must be a non-empty object')
cookie = security_schemes.get('cookieAuth', {})
if cookie.get('type') != 'apiKey' or cookie.get('in') != 'cookie':
    raise SystemExit('cookieAuth must be an apiKey cookie security scheme')

def resolve_pointer(ref):
    if not isinstance(ref, str) or not ref.startswith('#/'):
        raise SystemExit(f'Only local refs are allowed: {ref!r}')
    node = data
    for raw in ref[2:].split('/'):
        key = raw.replace('~1', '/').replace('~0', '~')
        if not isinstance(node, dict) or key not in node:
            raise SystemExit(f'Unresolved ref: {ref}')
        node = node[key]
    return node

def walk(node):
    if isinstance(node, dict):
        if '$ref' in node:
            resolve_pointer(node['$ref'])
        for value in node.values():
            walk(value)
    elif isinstance(node, list):
        for value in node:
            walk(value)

walk(data)

methods = {'get', 'post', 'put', 'patch', 'delete', 'options', 'head', 'trace'}
operation_ids = set()
operations = []
for route, path_item in paths.items():
    if not isinstance(path_item, dict):
        raise SystemExit(f'Path item must be an object: {route}')
    route_params = {
        p.get('name') for p in path_item.get('parameters', []) if isinstance(p, dict)
    }
    placeholders = set(re.findall(r'{([^}]+)}', route))
    for method, operation in path_item.items():
        if method not in methods:
            continue
        if not isinstance(operation, dict):
            raise SystemExit(f'Operation must be an object: {method.upper()} {route}')
        operation_id = operation.get('operationId')
        if not operation_id:
            raise SystemExit(f'Missing operationId: {method.upper()} {route}')
        if operation_id in operation_ids:
            raise SystemExit(f'Duplicate operationId: {operation_id}')
        operation_ids.add(operation_id)
        operations.append((method, route, operation))
        operation_params = {
            p.get('name'): p for p in operation.get('parameters', []) if isinstance(p, dict)
        }
        available = route_params | set(operation_params)
        if placeholders != available.intersection(placeholders):
            missing = placeholders - available
            raise SystemExit(f'Missing path parameters {missing}: {method.upper()} {route}')
        for name in placeholders:
            param = operation_params.get(name)
            if param and (param.get('in') != 'path' or param.get('required') is not True):
                raise SystemExit(f'Invalid path parameter {name}: {method.upper()} {route}')
        if not operation.get('responses'):
            raise SystemExit(f'Missing responses: {method.upper()} {route}')

expected_paths = {
    '/health', '/readiness', '/auth/setup-status', '/auth/setup',
    '/auth/setup/verify', '/auth/login', '/auth/session', '/auth/recovery',
    '/auth/recovery/verify', '/system/status', '/system/diagnostics',
    '/overview', '/activity', '/operations/{operation_id}',
    '/integrations/github/accounts', '/integrations/github/accounts/{account_id}',
    '/integrations/github/accounts/{account_id}/connection-test',
    '/integrations/github/accounts/{account_id}/repositories',
    '/integrations/github/app-manifest/registrations',
    '/integrations/github/app-manifest/callback',
    '/integrations/github/app-installations/callback',
    '/integrations/dokploy/servers', '/integrations/dokploy/servers/{server_id}',
    '/integrations/dokploy/servers/{server_id}/connection-test',
    '/integrations/dokploy/servers/{server_id}/applications',
    '/integrations/qvac/configuration', '/integrations/qvac/connection-test',
    '/integrations/qvac/status', '/projects', '/projects/{project_id}',
    '/projects/{project_id}/monitoring-configuration', '/settings/automation',
    '/incidents', '/incidents/{incident_id}', '/incidents/{incident_id}/log-events',
    '/incidents/{incident_id}/timeline', '/incidents/{incident_id}/investigations',
    '/investigations/{investigation_id}', '/investigations/{investigation_id}/evidence',
    '/incidents/{incident_id}/remediation', '/remediations/{remediation_id}',
    '/remediations/{remediation_id}/validation-results',
    '/remediations/{remediation_id}/pull-request', '/pull-requests',
    '/pull-requests/{pull_request_id}',
}
missing_paths = sorted(expected_paths - set(paths))
if missing_paths:
    raise SystemExit(f'Missing required paths: {missing_paths}')

public_operations = {
    ('get', '/health'), ('get', '/readiness'), ('get', '/auth/setup-status'),
    ('post', '/auth/setup'), ('post', '/auth/setup/verify'),
    ('post', '/auth/login'), ('post', '/auth/recovery'),
    ('post', '/auth/recovery/verify'),
    ('get', '/integrations/github/app-manifest/callback'),
    ('get', '/integrations/github/app-installations/callback'),
}
for method, route, operation in operations:
    is_public = (method, route) in public_operations
    if is_public and operation.get('security') != []:
        raise SystemExit(f'Public operation must declare security: []: {method.upper()} {route}')
    if not is_public and operation.get('security') == []:
        raise SystemExit(f'Unexpected public operation: {method.upper()} {route}')

documented_statuses = {
    str(status)
    for _, _, operation in operations
    for status in operation.get('responses', {})
}
required_error_statuses = {'400', '401', '403', '404', '409', '429', '500'}
missing_error_statuses = sorted(required_error_statuses - documented_statuses)
if missing_error_statuses:
    raise SystemExit(f'Missing required error response coverage: {missing_error_statuses}')

monitoring = schemas.get('MonitoringConfiguration', {})
expected_monitoring = {
    'enabled', 'error_patterns', 'ignored_patterns', 'grouping_window',
    'context_before', 'context_after',
}
if set(monitoring.get('properties', {})) != expected_monitoring:
    raise SystemExit('MonitoringConfiguration must contain exactly the approved fields')
for field, expected in {
    'enabled': False,
    'error_patterns': [],
    'ignored_patterns': [],
    'grouping_window': 'PT30M',
    'context_before': 20,
    'context_after': 20,
}.items():
    if monitoring['properties'][field].get('default') != expected:
        raise SystemExit(f'Invalid MonitoringConfiguration default: {field}')

automation = schemas.get('AutomationPolicy', {}).get('properties', {})
for field in ('automatic_investigation', 'automatic_remediation', 'automatic_pull_request'):
    if automation.get(field, {}).get('default') is not True:
        raise SystemExit(f'Automation default must be true: {field}')

sensitive_names = {
    'password', 'new_password', 'bootstrap_token', 'personal_access_token',
    'api_credential', 'bearer_token', 'basic_password', 'private_key',
    'webhook_secret', 'session_token', 'totp_secret', 'manual_entry_secret',
    'totp_code',
}
for schema_name, schema in schemas.items():
    properties = schema.get('properties', {}) if isinstance(schema, dict) else {}
    for name, prop in properties.items():
        if name in sensitive_names and prop.get('writeOnly') is not True:
            raise SystemExit(f'Sensitive property must be writeOnly: {schema_name}.{name}')

safe_resource_schemas = {
    'Administrator', 'Session', 'GitHubAccount', 'DokployServer',
    'QvacConfiguration', 'Project', 'Incident', 'Evidence', 'Remediation',
    'PullRequest', 'SystemStatus',
}
forbidden_resource_properties = sensitive_names | {
    'otpauth_uri', 'manual_entry_key', 'state', 'installation_token',
}

def collect_property_names(schema, seen=None):
    seen = set() if seen is None else seen
    if not isinstance(schema, dict):
        return set()
    if '$ref' in schema:
        ref = schema['$ref']
        if ref in seen:
            return set()
        seen.add(ref)
        return collect_property_names(resolve_pointer(ref), seen)
    names = set(schema.get('properties', {}))
    for prop in schema.get('properties', {}).values():
        names |= collect_property_names(prop, seen)
    for key in ('allOf', 'anyOf', 'oneOf'):
        for child in schema.get(key, []):
            names |= collect_property_names(child, seen)
    if isinstance(schema.get('items'), dict):
        names |= collect_property_names(schema['items'], seen)
    return names

for schema_name in safe_resource_schemas:
    names = collect_property_names(schemas.get(schema_name, {}))
    forbidden = sorted(names & forbidden_resource_properties)
    if forbidden:
        raise SystemExit(f'Secret-like properties in safe resource {schema_name}: {forbidden}')

expected_enums = {
    'RootCauseStatus': ['identified', 'suspected', 'unknown'],
    'ResolutionStatus': ['fixable', 'requires_human'],
    'InvestigationStatus': ['pending', 'running', 'completed', 'failed'],
    'RemediationStatus': ['planned', 'in_progress', 'validated', 'failed', 'pull_request_created'],
    'IncidentPhase': ['detected', 'investigating', 'publishing_issue', 'remediating', 'completed', 'failed'],
    'TerminalOutcome': ['pull_request_created', 'requires_human', 'remediation_failed', 'investigation_failed', 'issue_publication_failed'],
}
for schema_name, expected in expected_enums.items():
    if schemas.get(schema_name, {}).get('enum') != expected:
        raise SystemExit(f'Unexpected enum values: {schema_name}')

idempotent_commands = {
    'startIncidentInvestigation', 'startIncidentRemediation',
    'createRemediationPullRequest',
}
for _, _, operation in operations:
    if operation['operationId'] not in idempotent_commands:
        continue
    refs = [p.get('$ref') for p in operation.get('parameters', []) if isinstance(p, dict)]
    if '#/components/parameters/IdempotencyKey' not in refs:
        raise SystemExit(f'Missing Idempotency-Key: {operation["operationId"]}')
    if '202' not in operation['responses']:
        raise SystemExit(f'Idempotent command must return 202: {operation["operationId"]}')

if 'post' in paths.get('/incidents', {}):
    raise SystemExit('Manual incident creation is outside the MVP')
if any('merge' in route or 'deploy' in route or 'team' in route for route in paths):
    raise SystemExit('Out-of-scope merge/deploy/team path detected')

required_example_schemas = {
    'ErrorResponse', 'SetupRequest', 'LoginRequest',
    'MonitoringConfiguration', 'IncidentSummary', 'Operation',
}
for schema_name in required_example_schemas:
    if not schemas.get(schema_name, {}).get('examples'):
        raise SystemExit(f'Missing representative example: {schema_name}')

print(f'OpenAPI check passed: {len(operation_ids)} operations, {len(schemas)} schemas.')
PY
