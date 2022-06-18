from behave import *
from client import create_order


@given("A set of Orders")
def step_impl(context):
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


@when("create order")
def step_impl(context):
    context.responses = []
    for order in context.orders:
        context.responses.append(
            create_order(
                order["user_id"],
                order["value"],
                order["installments"],
                order["user_id"],  
            )
        )
        
        
@then('check status is "{status}"')
def step_impl(context, status):
    for response in context.responses:
        assert response["status"] == status
