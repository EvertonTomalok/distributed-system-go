import requests

url = "http://0.0.0.0:5000/api"


def create_order(user_id: str, value: float, installment: int, method: str) -> dict:
    payload={
        'value': '2000.00',
        'user_id': 'uuuu-aaaa-bbbb-invalid',
        'installment': '1',
        'method': 'credit_card',
    }

    return requests.request("POST", url + "/orders", data=payload).json()