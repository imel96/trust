Feature: Financing Activities

  Scenario: Add capital

    Given there's a trust in the database
      And the "cash" account is 0
      And the "capital" account is 0
    When the trust raises capital by 100 Dr
    Then the "cash" account balance becomes 100 Dr
      And the "capital" account balance becomes 100 Cr
      And the total "Financing" cash flow will be 100 Dr

  Scenario: Borrow loan

    Given there's a trust in the database
      And the "cash" account is 0
      And the "loan" account is 0
    When the trust borrow loan of 100
    Then the "cash" account becomes 100 Dr
      And the "loan" account becomes 100 Cr
      And the total "Financing" cash flow will be 100 Dr
