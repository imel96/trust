Feature: Trust management

Scenario: Add trust member
  Given there's a person in the database
  And there's a trust in the database
  When the person is added to the trust
  Then the person can be seen in the trust member list

Scenario: Add broker to trust
  Given there's a trust in the database
  And there's a broker in the database
  When the broker is added to the trust
  Then the broker can be seen in the trust broker list
