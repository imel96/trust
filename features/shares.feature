Feature: Shares management
  # assummed there's only one broker

  Scenario: Buy shares

    Given there's a trust in the database
      And the trust has 100 in "cash" account
      And the trust has 0 in "cost" account
      And there's a broker in the database
      And there's a company in the database
    When the trust buys some shares of the company through the broker
    Then the company shares show up in the trust portfolio
      And the trust's "cash" account become "0"
      And the trust's "cost" account become "1"

  Scenario: Sell shares

    Given there's a broker in the database
      And there's a trust in the database
      And the trust has 0 in "cash" account
      And the trust has 0 in "cost" account
      And the trust has some shares in the trust portfolio
      And there's a company in the database
    When the trust sells all shares of the company through the broker
    Then the trust portfolio shares of the company become "0"
      And the trust's "cash" account become "100"
      And the trust's "cost" account become "1"
