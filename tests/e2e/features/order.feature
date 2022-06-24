Feature: Orders Tests

    Scenario: Insert orders
        Given A set of Orders
            | value   | method      | installments | user_id   |
            | 1000.00 | credit_card |      10      | user_test |
            | 1000.00 | debit_card  |      1       | user_test |

        When create orders
        Then check status is "PROCESSING"
        Then check user has orders


    Scenario: Invalid method order
        Given A set of Orders
            | value   | method      | installments | user_id   |
            | 1000.00 | not_exist   |      50      | user_test |

        When create orders
        Then check response status_code 400

    @insert-approved-orders
    Scenario: Insert APPROVED orders
        Given A set of Orders
            | value   | method      | installments | user_id   |
            | 1500.00 | credit_card |      10      | user_test |
            | 150.00  | debit_card  |      1       | user_test |


        When create orders
        Then check after 0.5 second(s) status orders are "APPROVED"


    Scenario: Insert CANCELED orders
        Given A set of Orders
            | value       | method      | installments | user_id           |
            | 1500.00     | credit_card |      3       | user_test_invalid |
            | 1000000.00  | debit_card  |      1       | user_test         |


        When create orders
        Then check after 0.5 second(s) status orders are "CANCELED"

    Scenario: Delete All orders from user test
        Given Delete all orders from user test