Feature: Shares management
  # assummed there's only one broker

  Scenario: Buy shares

    Given there's a trust in the database
      And the trust has some money in the cash account
      And there's a broker in the database
      And there's a company in the database
    When the trust buys some shares of the company through the broker
    Then the company shares show up in the trust portfolio
      And the trust's cash account become empty
      And the trust's cost account increased
