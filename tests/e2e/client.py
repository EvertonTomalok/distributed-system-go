import os

import requests

url = os.getenv("BASE_URL", "http://0.0.0.0:8080/api")


def create_order(user_id: str, value: float, installment: int, method: str) -> dict:
    payload={
        'value': value,
        'user_id': user_id,
        'installment': installment,
        'method': method,
    }
    return requests.post(url + "/orders", data=payload)


def get_orders_from_user(user_id: str):
    return requests.get(url + f"/orders/{user_id}?limit=5").json()
    
    
def get_order_by_id(order_id: str):
    return requests.get(url + f"/order/{order_id}")


def delete_orders_from_user_test():
    return requests.delete(url + "/orders/user/test")
