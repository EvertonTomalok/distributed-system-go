from behave import given

from client import delete_orders_from_user_test

@given("Delete all orders from user test")
def delete_all_orders_user(context):
    delete_orders_from_user_test()

