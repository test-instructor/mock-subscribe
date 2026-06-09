import argparse
import hashlib
import json
import random
import string
import sys
import time
import uuid
import xml.etree.ElementTree as ET
from datetime import datetime
from pathlib import Path

import requests


SUCCESS_RETURN_CODE = 'SUCCESS'
SUCCESS_RESULT_CODE = 'SUCCESS'
ACTIVE_CONTRACT_STATUS = 'ACTIVE'
SUCCESS_TRADE_STATE = 'SUCCESS'
FAILED_TRADE_STATE = 'FAILED'
ORDER_ACCEPTED_STATE = 'ACCEPT'
ORDER_USERPAYING_STATE = 'USERPAYING'
ORDER_NOTPAY_STATE = 'NOTPAY'
ORDER_PAYERROR_STATE = 'PAYERROR'
ORDER_CLOSED_STATE = 'CLOSED'
ORDER_REFUND_STATE = 'REFUND'
ORDER_TERMINAL_STATES = {
    SUCCESS_TRADE_STATE,
    FAILED_TRADE_STATE,
    ORDER_PAYERROR_STATE,
    ORDER_CLOSED_STATE,
    ORDER_REFUND_STATE,
}
ORDER_PENDING_STATES = {
    '',
    ORDER_ACCEPTED_STATE,
    ORDER_USERPAYING_STATE,
    ORDER_NOTPAY_STATE,
}


def log_status(status: str, message: str, **fields) -> None:
    timestamp = datetime.now().strftime('%Y-%m-%d %H:%M:%S')
    extras = ' '.join(f'{key}={json.dumps(value, ensure_ascii=False)}' for key, value in fields.items())
    line = f'[{timestamp}] [{status}] {message}'
    if extras:
        line = f'{line} | {extras}'
    print(line, flush=True)


def load_config(config_path: str) -> dict:
    path = Path(config_path)
    if not path.exists():
        raise FileNotFoundError(f'config not found: {path}')
    return json.loads(path.read_text(encoding='utf-8'))


def random_nonce(length: int = 24) -> str:
    alphabet = string.ascii_letters + string.digits
    return ''.join(random.choice(alphabet) for _ in range(length))


def random_nonce_str(length: int = 32) -> str:
    return random_nonce(length)


def pick_first_non_empty(*values):
    for value in values:
        if str(value).strip() != '':
            return value
    return ''


def build_sign(params: dict, sign_key: str) -> str:
    keys = sorted(k for k, v in params.items() if k != 'sign' and str(v).strip() != '')
    raw = '&'.join(f'{key}={params[key]}' for key in keys)
    raw = f'{raw}&key={sign_key}'
    return hashlib.md5(raw.encode('utf-8')).hexdigest().upper()


def verify_sign(params: dict, sign_key: str) -> bool:
    actual_sign = (params.get('sign') or '').strip().upper()
    if not actual_sign:
        return False
    expected_sign = build_sign(params, sign_key)
    return expected_sign == actual_sign


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
    url = f"{base_url.rstrip('/')}{path}"
    log_status('START', '发起 XML 请求', method='POST', url=url)
    try:
        resp = requests.post(url, data=payload.encode('utf-8'), headers={'Content-Type': 'application/xml'}, timeout=60)
        resp.raise_for_status()
    except Exception as exc:
        log_status('FAILED', 'XML 请求失败', method='POST', url=url, error=str(exc))
        raise
    log_status('SUCCESS', 'XML 请求成功', method='POST', url=url, status_code=resp.status_code)
    return resp.text


def get_json(base_url: str, path: str, params: dict) -> dict:
    url = f"{base_url.rstrip('/')}{path}"
    log_status('START', '发起 JSON 请求', method='GET', url=url, params=params)
    try:
        resp = requests.get(url, params=params, timeout=60)
        resp.raise_for_status()
    except Exception as exc:
        log_status('FAILED', 'JSON 请求失败', method='GET', url=url, error=str(exc))
        raise
    log_status('SUCCESS', 'JSON 请求成功', method='GET', url=url, status_code=resp.status_code)
    return resp.json()


def post_json(base_url: str, path: str, payload: dict) -> tuple[dict, int]:
    url = f"{base_url.rstrip('/')}{path}"
    log_status('START', '发起 JSON 请求', method='POST', url=url, payload=payload)
    try:
        resp = requests.post(url, json=payload, timeout=60)
        resp.raise_for_status()
    except Exception as exc:
        log_status('FAILED', 'JSON 请求失败', method='POST', url=url, error=str(exc))
        raise
    log_status('SUCCESS', 'JSON 请求成功', method='POST', url=url, status_code=resp.status_code)
    if resp.text.strip():
        return resp.json(), resp.status_code
    return {}, resp.status_code


def is_xml_api_success(resp_data: dict) -> bool:
    return resp_data.get('return_code') == SUCCESS_RETURN_CODE and resp_data.get('result_code') == SUCCESS_RESULT_CODE


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
        'nonce': random_nonce(),
    }
    payload['sign'] = build_sign(payload, config['sign_key'])
    return payload


def make_contract_query_request(config: dict, contract_id: str, out_contract_code: str) -> dict:
    payload = {
        'appid': config['app_id'],
        'mch_id': config['mch_id'],
        'contract_id': contract_id,
        'out_contract_code': out_contract_code,
        'sign_type': config['sign_type'],
        'timestamp': str(int(time.time())),
        'nonce': random_nonce(),
    }
    payload['sign'] = build_sign(payload, config['sign_key'])
    return payload


def make_deduct_apply_request(config: dict, contract_id: str, out_trade_no: str, total_amount: int, transaction_id: str = '') -> dict:
    payload = {
        'appid': config['app_id'],
        'mch_id': config['mch_id'],
        'out_trade_no': out_trade_no,
        'contract_id': contract_id,
        'transaction_id': transaction_id,
        'body': config.get('deduct_body', config.get('body', '委托代扣')),
        'detail': config.get('deduct_detail', ''),
        'attach': config.get('deduct_attach', ''),
        'total_fee': total_amount,
        'total_amount': total_amount,
        'fee_type': config.get('fee_type', 'CNY'),
        'notify_url': config.get('deduct_notify_url', config['notify_url']),
        'trade_type': config.get('trade_type', 'PAP'),
        'device_info': config.get('device_info', ''),
        'nonce_str': random_nonce_str(),
        'sign_type': config['sign_type'],
        'timestamp': str(int(time.time())),
        'nonce': random_nonce(),
    }
    payload['sign'] = build_sign(payload, config['sign_key'])
    return payload


def make_order_query_request(config: dict, out_trade_no: str, transaction_id: str) -> dict:
    payload = {
        'appid': config['app_id'],
        'mch_id': config['mch_id'],
        'out_trade_no': out_trade_no,
        'transaction_id': transaction_id,
        'nonce_str': random_nonce_str(),
        'sign_type': config['sign_type'],
        'timestamp': str(int(time.time())),
        'nonce': random_nonce(),
    }
    payload['sign'] = build_sign(payload, config['sign_key'])
    return payload


def make_pre_notify_request(config: dict, total_amount: int) -> dict:
    return {
        'mchid': config['mch_id'],
        'appid': config['app_id'],
        'deduct_duration': {
            'count': int(config.get('deduct_duration_count', 1)),
            'unit': config.get('deduct_duration_unit', 'DAY'),
        },
        'estimated_amount': {
            'amount': int(config.get('estimated_amount', total_amount)),
            'currency': config.get('estimated_currency', 'CNY'),
        },
    }


def sign_flow(config: dict, out_contract_code: str | None = None, openid: str | None = None) -> dict:
    actual_out_contract_code = out_contract_code or f"{config['out_contract_code_prefix']}-{uuid.uuid4().hex[:12]}"
    actual_openid = openid or f"{config['openid_prefix']}_{uuid.uuid4().hex[:10]}"
    log_status('START', '开始签约请求', out_contract_code=actual_out_contract_code, openid=actual_openid)
    req = make_sign_request(config, actual_out_contract_code, actual_openid)
    xml_payload = dict_to_xml(req)
    resp_xml = post_xml(config['base_url'], config['sign_path'], xml_payload)
    resp_data = xml_to_dict(resp_xml)
    sign_valid = verify_sign(resp_data, config['sign_key']) if resp_data.get('sign') else None
    result = {
        'request': req,
        'request_xml': xml_payload,
        'response_xml': resp_xml,
        'response': resp_data,
        'response_sign_valid': sign_valid,
        'out_contract_code': actual_out_contract_code,
        'openid': actual_openid,
        'contract_id': resp_data.get('contract_id', ''),
    }
    log_status(
        'SUCCESS' if is_xml_api_success(resp_data) else 'FAILED',
        '签约请求完成',
        out_contract_code=actual_out_contract_code,
        contract_id=result['contract_id'],
        return_code=resp_data.get('return_code', ''),
        result_code=resp_data.get('result_code', ''),
        response_sign_valid=sign_valid,
    )
    return result


def wait_callback_flow(config: dict, out_contract_code: str, callback_kind: str = 'contract') -> dict:
    timeout_seconds = int(config.get('poll_timeout_seconds', 60))
    deadline = time.time() + timeout_seconds
    interval = int(config.get('poll_interval_seconds', 2))
    latest = None
    attempt = 0
    callback_list_path = config['callback_list_path']
    query_params = {
        'page': 1,
        'pageSize': 20,
        'outContractCode': out_contract_code,
    }
    wait_message = '开始等待回调'
    polling_message = '轮询回调中'
    success_message = '已收到回调'
    pending_message = '暂未收到回调，继续等待'
    timeout_message = '等待回调超时'
    if callback_kind == 'deduct':
        callback_list_path = config.get('deduct_callback_list_path', config['callback_list_path'])
        query_params = {
            'page': 1,
            'pageSize': 20,
            'outTradeNo': out_contract_code,
        }
        wait_message = '开始等待代扣回调'
        polling_message = '轮询代扣回调中'
        success_message = '已收到代扣回调'
        pending_message = '暂未收到代扣回调，继续等待'
        timeout_message = '等待代扣回调超时'
    log_status('START', wait_message, out_contract_code=out_contract_code, timeout_seconds=timeout_seconds, interval_seconds=interval)
    while time.time() < deadline:
        attempt += 1
        log_status('WAITING', polling_message, out_contract_code=out_contract_code, attempt=attempt)
        latest = get_json(config['base_url'], callback_list_path, query_params)
        if latest.get('code') == 0:
            records = latest.get('data', {}).get('list', []) or []
            if records:
                log_status('SUCCESS', success_message, out_contract_code=out_contract_code, attempt=attempt, record_count=len(records))
                return latest
        remaining_seconds = max(0, int(deadline - time.time()))
        log_status('WAITING', pending_message, out_contract_code=out_contract_code, attempt=attempt, remaining_seconds=remaining_seconds)
        time.sleep(interval)
    log_status('FAILED', timeout_message, out_contract_code=out_contract_code, timeout_seconds=timeout_seconds)
    raise TimeoutError(f'callback not received in {timeout_seconds} seconds')


def query_contract_until_active_flow(config: dict, contract_id: str, out_contract_code: str) -> dict:
    timeout_seconds = int(config.get('poll_timeout_seconds', 60))
    deadline = time.time() + timeout_seconds
    interval = int(config.get('poll_interval_seconds', 2))
    latest = None
    attempt = 0
    current_contract_id = contract_id
    log_status(
        'START',
        '开始主动查询签约关系',
        contract_id=contract_id,
        out_contract_code=out_contract_code,
        timeout_seconds=timeout_seconds,
        interval_seconds=interval,
    )
    while time.time() < deadline:
        attempt += 1
        log_status('WAITING', '轮询签约关系中', contract_id=current_contract_id, out_contract_code=out_contract_code, attempt=attempt)
        latest = query_flow(config, current_contract_id, out_contract_code)
        response = latest['response']
        current_contract_id = latest.get('contract_id') or current_contract_id
        if is_xml_api_success(response):
            contract_status = latest.get('contract_status', '')
            if contract_status == ACTIVE_CONTRACT_STATUS:
                log_status(
                    'SUCCESS',
                    '签约关系已生效',
                    contract_id=current_contract_id,
                    out_contract_code=out_contract_code,
                    attempt=attempt,
                    contract_status=contract_status,
                )
                return latest
        remaining_seconds = max(0, int(deadline - time.time()))
        log_status(
            'WAITING',
            '签约关系未就绪，继续查询',
            contract_id=current_contract_id,
            out_contract_code=out_contract_code,
            attempt=attempt,
            contract_status=(latest or {}).get('contract_status', ''),
            remaining_seconds=remaining_seconds,
        )
        time.sleep(interval)
    log_status('FAILED', '主动查询签约关系超时', contract_id=current_contract_id, out_contract_code=out_contract_code, timeout_seconds=timeout_seconds)
    raise TimeoutError(f'contract not active in {timeout_seconds} seconds')


def query_flow(config: dict, contract_id: str, out_contract_code: str) -> dict:
    log_status('START', '开始查询签约状态', contract_id=contract_id, out_contract_code=out_contract_code)
    req = make_contract_query_request(config, contract_id, out_contract_code)
    xml_payload = dict_to_xml(req)
    resp_xml = post_xml(config['base_url'], config['query_path'], xml_payload)
    resp_data = xml_to_dict(resp_xml)
    sign_valid = verify_sign(resp_data, config['sign_key']) if resp_data.get('sign') else None
    result = {
        'request': req,
        'request_xml': xml_payload,
        'response_xml': resp_xml,
        'response': resp_data,
        'response_sign_valid': sign_valid,
        'contract_id': resp_data.get('contract_id', contract_id),
        'out_contract_code': out_contract_code,
        'contract_status': resp_data.get('contract_status', ''),
    }
    log_status(
        'SUCCESS' if is_xml_api_success(resp_data) else 'FAILED',
        '签约状态查询完成',
        contract_id=result['contract_id'],
        out_contract_code=out_contract_code,
        contract_status=result['contract_status'],
        return_code=resp_data.get('return_code', ''),
        result_code=resp_data.get('result_code', ''),
        response_sign_valid=sign_valid,
    )
    return result


def apply_deduct_flow(config: dict, contract_id: str, total_amount: int, out_trade_no: str | None = None, transaction_id: str = '') -> dict:
    actual_out_trade_no = out_trade_no or f"{config.get('out_trade_no_prefix', 'MOCK-ORDER')}-{uuid.uuid4().hex[:12]}"
    log_status('START', '开始申请扣款', contract_id=contract_id, out_trade_no=actual_out_trade_no, total_amount=total_amount)
    req = make_deduct_apply_request(config, contract_id, actual_out_trade_no, total_amount, transaction_id)
    xml_payload = dict_to_xml(req)
    resp_xml = post_xml(config['base_url'], config['deduct_apply_path'], xml_payload)
    resp_data = xml_to_dict(resp_xml)
    sign_valid = verify_sign(resp_data, config['sign_key']) if resp_data.get('sign') else None
    result = {
        'request': req,
        'request_xml': xml_payload,
        'response_xml': resp_xml,
        'response': resp_data,
        'response_sign_valid': sign_valid,
        'out_trade_no': actual_out_trade_no,
        'transaction_id': resp_data.get('transaction_id', transaction_id),
        'deduct_status': SUCCESS_TRADE_STATE if is_xml_api_success(resp_data) else FAILED_TRADE_STATE,
    }
    log_status(
        'SUCCESS' if is_xml_api_success(resp_data) else 'FAILED',
        '申请扣款完成',
        contract_id=contract_id,
        out_trade_no=actual_out_trade_no,
        transaction_id=result['transaction_id'],
        deduct_status=result['deduct_status'],
        return_code=resp_data.get('return_code', ''),
        result_code=resp_data.get('result_code', ''),
        response_sign_valid=sign_valid,
    )
    return result


def pre_notify_flow(config: dict, contract_id: str, total_amount: int) -> dict:
    log_status('START', '开始预扣费通知', contract_id=contract_id, total_amount=total_amount)
    req = make_pre_notify_request(config, total_amount)
    path_template = config['pre_notify_path']
    path = path_template.format(contract_id=contract_id)
    resp_data, status_code = post_json(config['base_url'], path, req)
    result = {
        'request': req,
        'response': resp_data,
        'status_code': status_code,
        'success': resp_data.get('return_code') == SUCCESS_RETURN_CODE and resp_data.get('result_code') == SUCCESS_RESULT_CODE,
    }
    log_status(
        'SUCCESS' if result['success'] else 'FAILED',
        '预扣费通知完成',
        contract_id=contract_id,
        status_code=status_code,
        return_code=resp_data.get('return_code', ''),
        result_code=resp_data.get('result_code', ''),
        err_code=resp_data.get('err_code', ''),
    )
    return result


def query_order_flow(config: dict, out_trade_no: str, transaction_id: str = '') -> dict:
    log_status('START', '开始查询订单', out_trade_no=out_trade_no, transaction_id=transaction_id)
    req = make_order_query_request(config, out_trade_no, transaction_id)
    xml_payload = dict_to_xml(req)
    resp_xml = post_xml(config['base_url'], config['query_order_path'], xml_payload)
    resp_data = xml_to_dict(resp_xml)
    sign_valid = verify_sign(resp_data, config['sign_key']) if resp_data.get('sign') else None
    trade_state = resp_data.get('trade_state', '')
    result = {
        'request': req,
        'request_xml': xml_payload,
        'response_xml': resp_xml,
        'response': resp_data,
        'response_sign_valid': sign_valid,
        'out_trade_no': resp_data.get('out_trade_no', out_trade_no),
        'transaction_id': resp_data.get('transaction_id', transaction_id),
        'trade_state': trade_state,
    }
    log_status(
        'SUCCESS' if is_xml_api_success(resp_data) else 'FAILED',
        '订单查询完成',
        out_trade_no=result['out_trade_no'],
        transaction_id=result['transaction_id'],
        trade_state=trade_state,
        return_code=resp_data.get('return_code', ''),
        result_code=resp_data.get('result_code', ''),
        response_sign_valid=sign_valid,
    )
    if trade_state in ORDER_TERMINAL_STATES:
        log_status(
            'SUCCESS' if trade_state == SUCCESS_TRADE_STATE else 'FAILED',
            '订单状态已落定',
            out_trade_no=result['out_trade_no'],
            transaction_id=result['transaction_id'],
            trade_state=trade_state,
        )
    return result


def query_order_until_terminal_flow(config: dict, out_trade_no: str, transaction_id: str = '') -> dict:
    timeout_seconds = int(config.get('poll_timeout_seconds', 60))
    deadline = time.time() + timeout_seconds
    interval = int(config.get('poll_interval_seconds', 2))
    latest = None
    attempt = 0
    current_transaction_id = transaction_id
    log_status(
        'START',
        '开始轮询查询订单',
        out_trade_no=out_trade_no,
        transaction_id=transaction_id,
        timeout_seconds=timeout_seconds,
        interval_seconds=interval,
    )
    while time.time() < deadline:
        attempt += 1
        log_status(
            'WAITING',
            '轮询订单状态中',
            out_trade_no=out_trade_no,
            transaction_id=current_transaction_id,
            attempt=attempt,
        )
        latest = query_order_flow(config, out_trade_no, current_transaction_id)
        response = latest['response']
        trade_state = latest.get('trade_state', '')
        current_transaction_id = latest.get('transaction_id') or current_transaction_id
        if is_xml_api_success(response) and trade_state in ORDER_TERMINAL_STATES:
            log_status(
                'SUCCESS' if trade_state == SUCCESS_TRADE_STATE else 'FAILED',
                '订单轮询结束',
                out_trade_no=latest['out_trade_no'],
                transaction_id=current_transaction_id,
                trade_state=trade_state,
                attempt=attempt,
            )
            return latest
        if not is_xml_api_success(response) and response.get('err_code') == 'ORDERNOTEXIST':
            remaining_seconds = max(0, int(deadline - time.time()))
            log_status(
                'WAITING',
                '订单暂不存在，继续查询',
                out_trade_no=out_trade_no,
                transaction_id=current_transaction_id,
                attempt=attempt,
                remaining_seconds=remaining_seconds,
                err_code=response.get('err_code', ''),
                err_code_des=response.get('err_code_des', ''),
            )
            time.sleep(interval)
            continue
        if trade_state in ORDER_PENDING_STATES:
            remaining_seconds = max(0, int(deadline - time.time()))
            log_status(
                'WAITING',
                '订单未到终态，继续查询',
                out_trade_no=latest['out_trade_no'],
                transaction_id=current_transaction_id,
                trade_state=trade_state,
                attempt=attempt,
                remaining_seconds=remaining_seconds,
            )
            time.sleep(interval)
            continue
        if not is_xml_api_success(response):
            log_status(
                'FAILED',
                '订单查询失败，结束轮询',
                out_trade_no=latest['out_trade_no'],
                transaction_id=current_transaction_id,
                attempt=attempt,
                err_code=response.get('err_code', ''),
                err_code_des=response.get('err_code_des', ''),
                return_code=response.get('return_code', ''),
                result_code=response.get('result_code', ''),
            )
            return latest
        remaining_seconds = max(0, int(deadline - time.time()))
        log_status(
            'WAITING',
            '订单状态未识别为终态，继续查询',
            out_trade_no=latest['out_trade_no'],
            transaction_id=current_transaction_id,
            trade_state=trade_state,
            attempt=attempt,
            remaining_seconds=remaining_seconds,
        )
        time.sleep(interval)
    log_status(
        'FAILED',
        '订单查询超时',
        out_trade_no=out_trade_no,
        transaction_id=current_transaction_id,
        timeout_seconds=timeout_seconds,
    )
    raise TimeoutError(f'order not settled in {timeout_seconds} seconds')


def deduct_after_query_flow(config: dict, contract_id: str, out_contract_code: str, total_amount: int, out_trade_no: str | None = None) -> dict:
    log_status('START', '开始执行查约后扣款流程', contract_id=contract_id, out_contract_code=out_contract_code, total_amount=total_amount)
    query_result = query_contract_until_active_flow(config, contract_id, out_contract_code)
    query_response = query_result['response']
    if not is_xml_api_success(query_response):
        log_status('FAILED', '签约查询失败，终止扣款流程', contract_id=contract_id, out_contract_code=out_contract_code)
        return {
            'query': query_result,
            'deduct': None,
            'order': None,
            'pre_notify': None,
            'next_deduct': None,
            'next_order': None,
        }

    contract_status = query_response.get('contract_status', '')
    if contract_status != ACTIVE_CONTRACT_STATUS:
        log_status('FAILED', '签约状态不是已签约，终止扣款流程', contract_id=contract_id, out_contract_code=out_contract_code, contract_status=contract_status)
        return {
            'query': query_result,
            'deduct': None,
            'order': None,
            'pre_notify': None,
            'next_deduct': None,
            'next_order': None,
        }

    log_status('SUCCESS', '签约状态为已签约，开始申请扣款', contract_id=query_result['contract_id'], out_contract_code=out_contract_code, contract_status=contract_status)
    deduct_result = apply_deduct_flow(config, query_result['contract_id'], total_amount, out_trade_no=out_trade_no)
    deduct_response = deduct_result['response']
    if not is_xml_api_success(deduct_response):
        log_status('FAILED', '扣款请求失败，不再查询订单', contract_id=query_result['contract_id'], out_trade_no=deduct_result['out_trade_no'])
        return {
            'query': query_result,
            'deduct': deduct_result,
            'order': None,
            'pre_notify': None,
            'next_deduct': None,
            'next_order': None,
        }

    order_result = query_order_until_terminal_flow(config, deduct_result['out_trade_no'], deduct_result['transaction_id'])
    pre_notify_result = None
    next_deduct_result = None
    next_order_result = None

    if order_result.get('trade_state') == SUCCESS_TRADE_STATE:
        pre_notify_result = pre_notify_flow(config, query_result['contract_id'], total_amount)
        if pre_notify_result.get('success'):
            wait_seconds = int(config.get('pre_notify_wait_seconds', 60))
            log_status('WAITING', '等待下一次扣费窗口', wait_seconds=wait_seconds)
            time.sleep(wait_seconds)
            next_deduct_result = apply_deduct_flow(config, query_result['contract_id'], total_amount)
            next_deduct_response = next_deduct_result['response']
            if is_xml_api_success(next_deduct_response):
                next_order_result = query_order_until_terminal_flow(config, next_deduct_result['out_trade_no'], next_deduct_result['transaction_id'])

    return {
        'query': query_result,
        'deduct': deduct_result,
        'order': order_result,
        'pre_notify': pre_notify_result,
        'next_deduct': next_deduct_result,
        'next_order': next_order_result,
    }


def e2e_flow(config: dict, total_amount: int | None = None) -> dict:
    amount = total_amount if total_amount is not None else int(config.get('total_amount', 100))
    log_status('START', '开始执行完整流程', total_amount=amount)
    sign_result = sign_flow(config)
    flow_result = deduct_after_query_flow(config, sign_result['contract_id'], sign_result['out_contract_code'], amount)
    log_status('SUCCESS', '完整流程执行完成', out_contract_code=sign_result['out_contract_code'], contract_id=flow_result['query']['contract_id'])
    return {
        'sign': sign_result,
        'query': flow_result['query'],
        'deduct': flow_result['deduct'],
        'order': flow_result['order'],
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
    wait_cmd.add_argument('--kind', choices=['contract', 'deduct'], default='contract', help='callback record type to poll')

    query_cmd = sub.add_parser('query', help='query contract status')
    query_cmd.add_argument('--contract-id', required=True)
    query_cmd.add_argument('--out-contract-code', required=True)

    deduct_cmd = sub.add_parser('deduct', help='query contract until active, then apply deduct and query order')
    deduct_cmd.add_argument('--contract-id', required=True)
    deduct_cmd.add_argument('--out-contract-code', required=True)
    deduct_cmd.add_argument('--total-amount', type=int, required=True)
    deduct_cmd.add_argument('--out-trade-no', default='')

    order_cmd = sub.add_parser('query-order', help='query deduct order')
    order_cmd.add_argument('--out-trade-no', required=True)
    order_cmd.add_argument('--transaction-id', default='')
    order_cmd.add_argument('--wait-until-terminal', action='store_true', help='poll order status until terminal state')

    e2e_cmd = sub.add_parser('e2e', help='run sign -> query-contract-until-active -> deduct -> query-order')
    e2e_cmd.add_argument('--total-amount', type=int, default=0)

    args = parser.parse_args()
    config = load_config(args.config)
    command = args.command or 'e2e'
    log_status('START', '脚本开始执行', command=command, config=args.config)

    if command == 'sign':
        result = sign_flow(config, args.out_contract_code or None, args.openid or None)
    elif command == 'wait-callback':
        result = wait_callback_flow(config, args.out_contract_code, args.kind)
    elif command == 'query':
        result = query_flow(config, args.contract_id, args.out_contract_code)
    elif command == 'deduct':
        result = deduct_after_query_flow(config, args.contract_id, args.out_contract_code, args.total_amount, args.out_trade_no or None)
    elif command == 'query-order':
        if args.wait_until_terminal:
            result = query_order_until_terminal_flow(config, args.out_trade_no, args.transaction_id)
        else:
            result = query_order_flow(config, args.out_trade_no, args.transaction_id)
    else:
        total_amount = getattr(args, 'total_amount', 0) or None
        result = e2e_flow(config, total_amount)

    log_status('SUCCESS', '脚本执行完成', command=command)
    print(json.dumps(result, ensure_ascii=False, indent=2))
    return 0


if __name__ == '__main__':
    try:
        raise SystemExit(main())
    except Exception as exc:
        log_status('FAILED', '脚本执行失败', error=str(exc))
        print(json.dumps({'error': str(exc)}, ensure_ascii=False, indent=2), file=sys.stderr)
        raise
