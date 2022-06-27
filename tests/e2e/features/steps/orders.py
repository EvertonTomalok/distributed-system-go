from behave import *
from time import sleep

from client import create_order, get_order_by_id, get_orders_from_user


@given("A set of Orders")
def step_impl(context):
    context.user_id = context.table[0]["user_id"]
    context.orders = []
    for row in context.table:
        context.orders.append(
            {
                "value": row["value"],
                "method": row["method"],
                "installments": row["installments"],
                "user_id": row["user_id"],
            }
        )


@when("create orders")
def step_impl(context):
    context.responses = []
    for order in context.orders:
        context.responses.append(
            create_order(
                order["user_id"],
                order["value"],
                order["installments"],
                order["method"],
            )
        )


@then('check status is "{status}"')
def step_impl(context, status):
    for response in context.responses:
        r = response.json()
        assert r["status"] == status


@then("check user have orders")
def step_impl(context):
    assert get_orders_from_user(context.user_id)["total"] > 0


@then("check response status_code {status_code:d}")
def step_impl(context, status_code):
    for response in context.responses:
        assert response.status_code == status_code


@then('check after {sec:g} second(s) status orders are "{status}"')
def step_impl(context, sec, status):
    sleep(sec)

    for order_response in context.responses:
        order = get_order_by_id(order_response.json()["id"])
        assert order.status_code == 200
        assert order.json()["status"] == status
