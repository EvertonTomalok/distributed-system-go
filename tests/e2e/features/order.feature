Feature: Orders Tests

Scenario: Insert Valid ordem
    Given A set of Orders
        | value   | method      | installments | user_id |
        | 1000.00 | credit_card |      10      | user    |
        | 1000.00 | debit_card  |      1       | user    |

    When create order
    Then check status is PROCESSING
