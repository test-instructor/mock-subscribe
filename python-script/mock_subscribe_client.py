import argparse
import hashlib
import json
import random
import string
import sys
import time
import uuid
import xml.etree.ElementTree as ET
from pathlib import Path

import requests


def load_config(config_path: str) -> dict:
    path = Path(config_path)
    if not path.exists():
        raise FileNotFoundError(f'config not found: {path}')
    return json.loads(path.read_text(encoding='utf-8'))


def random_nonce(length: int = 24) -> str:
    alphabet = string.ascii_letters + string.digits
    return ''.join(random.choice(alphabet) for _ in range(length))


def build_sign(params: dict, sign_key: str) -> str:
    keys = sorted(k for k, v in params.items() if k != 'sign' and str(v).strip() != '')
    raw = '&'.join(f'{key}={params[key]}' for key in keys)
    raw = f'{raw}&key={sign_key}'
    return hashlib.md5(raw.encode('utf-8')).hexdigest().upper()


def dict_to_xml(data: dict) -> str:
    root = ET.Element('xml')
    for key, value in data.items():
        child = ET.SubElement(root, key)
        child.text = str(value)
    return ET.tostring(root, encoding='utf-8', xml_declaration=True).decode('utf-8')


def xml_to_dict(xml_text: str) -> dict:
    root = ET.fromstring(xml_text)
    return {child.tag: (child.text or '') for child in root}


def post_xml(base_url: str, path: str, payload: str) -> str:
    resp = requests.post(f"{base_url.rstrip('/')}{path}", data=payload.encode('utf-8'), headers={'Content-Type': 'application/xml'}, timeout=60)
    resp.raise_for_status()
    return resp.text


def get_json(base_url: str, path: str, params: dict) -> dict:
    resp = requests.get(f"{base_url.rstrip('/')}{path}", params=params, timeout=60)
    resp.raise_for_status()
    return resp.json()


def make_sign_request(config: dict, out_contract_code: str, openid: str) -> dict:
    payload = {
        'appid': config['app_id'],
        'mch_id': config['mch_id'],
        'plan_id': config['plan_id'],
        'out_contract_code': out_contract_code,
        'outer_openid': openid,
        'contract_display_account': config['contract_display_account'],
        'notify_url': config['notify_url'],
        'sign_type': config['sign_type'],
        'version': '1.0',
        'timestamp': str(int(time.time())),
        'nonce': random_nonce()
    }
    payload['sign'] = build_sign(payload, config['sign_key'])
    return payload


def make_query_request(config: dict, contract_id: str, out_contract_code: str) -> dict:
    payload = {
        'appid': config['app_id'],
        'mch_id': config['mch_id'],
        'contract_id': contract_id,
        'out_contract_code': out_contract_code,
        'sign_type': config['sign_type'],
        'timestamp': str(int(time.time())),
        'nonce': random_nonce()
    }
    payload['sign'] = build_sign(payload, config['sign_key'])
    return payload


def sign_flow(config: dict, out_contract_code: str | None = None, openid: str | None = None) -> dict:
    actual_out_contract_code = out_contract_code or f"{config['out_contract_code_prefix']}-{uuid.uuid4().hex[:12]}"
    actual_openid = openid or f"{config['openid_prefix']}_{uuid.uuid4().hex[:10]}"
    req = make_sign_request(config, actual_out_contract_code, actual_openid)
    xml_payload = dict_to_xml(req)
    resp_xml = post_xml(config['base_url'], config['sign_path'], xml_payload)
    resp_data = xml_to_dict(resp_xml)
    result = {
        'request': req,
        'request_xml': xml_payload,
        'response_xml': resp_xml,
        'response': resp_data,
        'out_contract_code': actual_out_contract_code,
        'openid': actual_openid,
        'contract_id': resp_data.get('contract_id', '')
    }
    return result


def wait_callback_flow(config: dict, out_contract_code: str) -> dict:
    deadline = time.time() + int(config.get('poll_timeout_seconds', 60))
    interval = int(config.get('poll_interval_seconds', 2))
    latest = None
    while time.time() < deadline:
        latest = get_json(config['base_url'], config['callback_list_path'], {
            'page': 1,
            'pageSize': 20,
            'outContractCode': out_contract_code
        })
        if latest.get('code') == 0:
            records = latest.get('data', {}).get('list', []) or []
            if records:
                return latest
        time.sleep(interval)
    raise TimeoutError(f'callback not received in {config.get("poll_timeout_seconds", 60)} seconds')


def query_flow(config: dict, contract_id: str, out_contract_code: str) -> dict:
    req = make_query_request(config, contract_id, out_contract_code)
    xml_payload = dict_to_xml(req)
    resp_xml = post_xml(config['base_url'], config['query_path'], xml_payload)
    return {
        'request': req,
        'request_xml': xml_payload,
        'response_xml': resp_xml,
        'response': xml_to_dict(resp_xml)
    }


def e2e_flow(config: dict) -> dict:
    sign_result = sign_flow(config)
    callback_result = wait_callback_flow(config, sign_result['out_contract_code'])
    query_result = query_flow(config, sign_result['contract_id'], sign_result['out_contract_code'])
    return {
        'sign': sign_result,
        'callback': callback_result,
        'query': query_result
    }


def main() -> int:
    parser = argparse.ArgumentParser(description='WeChat mock subscribe client')
    parser.add_argument('--config', default='config.example.json', help='config json path')
    sub = parser.add_subparsers(dest='command', required=False)

    sign_cmd = sub.add_parser('sign', help='send contract sign request')
    sign_cmd.add_argument('--out-contract-code', default='')
    sign_cmd.add_argument('--openid', default='')

    wait_cmd = sub.add_parser('wait-callback', help='poll callback records')
    wait_cmd.add_argument('--out-contract-code', required=True)

    query_cmd = sub.add_parser('query', help='query contract status')
    query_cmd.add_argument('--contract-id', required=True)
    query_cmd.add_argument('--out-contract-code', required=True)

    sub.add_parser('e2e', help='run sign -> wait-callback -> query')

    args = parser.parse_args()
    config = load_config(args.config)
    command = args.command or 'e2e'

    if command == 'sign':
        result = sign_flow(config, args.out_contract_code or None, args.openid or None)
    elif command == 'wait-callback':
        result = wait_callback_flow(config, args.out_contract_code)
    elif command == 'query':
        result = query_flow(config, args.contract_id, args.out_contract_code)
    else:
        result = e2e_flow(config)

    print(json.dumps(result, ensure_ascii=False, indent=2))
    return 0


if __name__ == '__main__':
    try:
        raise SystemExit(main())
    except Exception as exc:
        print(json.dumps({'error': str(exc)}, ensure_ascii=False, indent=2), file=sys.stderr)
        raise
