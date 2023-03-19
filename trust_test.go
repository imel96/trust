package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/cucumber/godog"
	"os/exec"
	"strings"
)

const path = "/Users/melbysjamsuddin/Documents/trust_test"

type brokerCtxKey struct{}
type personCtxKey struct{}
type trustCashAccountCtxKey struct{}
type trustCostAccountCtxKey struct{}
type trustCtxKey struct{}

func hashString(text string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(text)))
}

func executeQuery(query string) string {
	return executeQueries([]string{query})
}

func executeQueries(queries []string) string {
	var sqls string

	for _, sql := range queries {
		sqls = sqls + sql
		//fmt.Printf("%s\n", sql)
	}
	var out bytes.Buffer
	cmd := exec.Command("sqlite3", path, sqls)
	cmd.Stdout = &out
	cmd.Run()
	return out.String()
}

func insertAccount(account_number string, name string, balance int) {
	executeQueries([]string{
		"INSERT OR IGNORE INTO account_hub(account_hash, account_number) values('" + hashString(account_number) + "', '" + account_number + "');",
		"INSERT INTO account_sat(hash, account_hash, name, balance) values('" + hashString(account_number+"|"+name+"|"+fmt.Sprintf("%d", balance)) + "', '" + hashString(account_number) + "', '" + name + "', '" + fmt.Sprintf("%d", balance) + "');",
	})
}

func insertCompany(code string, name string, gics string) {
	executeQueries([]string{
		"INSERT INTO company_hub(company_hash, code) values('" + hashString(code) + "', '" + code + "');",
		"INSERT INTO company_sat values('" + hashString(code+"|"+name+"|"+gics) + "', '" + hashString(code) + "', '" + name + "', '" + gics + "');",
	})
}

func insertParty(name string, party_type string) {
	executeQueries([]string{
		"INSERT INTO party_hub(party_hash, name) values('" + hashString(name) + "', '" + name + "');",
		"INSERT INTO party_sat(hash, party_hash, party_type) values('" + hashString(name+"|"+party_type) + "', '" + hashString(name) + "', '" + party_type + "');",
	})
}

func insertLinkWithTwoHubs(link string, a string, b string) {
	executeQuery("INSERT INTO " + link + " values('" + hashString(a+"|"+b) + "', '" + hashString(a) + "', '" + hashString(b) + "', STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW'));")
}

func insertTrustCompanySat(trust string, company string, shareQuantity int) {
	executeQuery("INSERT INTO trust_company_sat(hash, company_hash, name, gics) values('" + hashString(trust+"|"+company) + "', " + fmt.Sprintf("%d", shareQuantity) + ");")
}

func theCompanySharesShowUpInTheTrustPortfolio() error {
	out := executeQuery("select trust.name || \"|\" || company_hub.code from party_hub trust join trust_company_link on (trust_hash = trust.party_hash) JOIN company_hub USING (company_hash);")

	if strings.TrimRight(out, "\n") != "shamsuddins|AZJ" {
		return fmt.Errorf("expected response to be: %s, but actual is: %s", "shamsuddins|AZJ", out)
	}
	return nil
}

func theTrustHasSomeMoneyInTheCashAccount() error {
	insertAccount("1", "shamsuddins cash", 100)
	insertAccount("5", "shamsuddins cash", 0)
	insertLinkWithTwoHubs("trust_account_link", "shamsuddins", "1")
	insertLinkWithTwoHubs("trust_account_link", "shamsuddins", "5")
	return nil
}

func theTrustsCashAccountBecomeEmpty(ctx context.Context) error {
	out := executeQuery("SELECT account_sat.name || \"|\" || account_sat.balance FROM trust_account_link JOIN account_hub USING(account_hash) JOIN account_sat USING(account_hash) WHERE account_hash = '" + hashString(fmt.Sprintf("%v", ctx.Value(trustCashAccountCtxKey{}))) + "' ORDER BY account_sat.load_time DESC LIMIT 1;")

	if strings.TrimRight(out, "\n") != "shamsuddins cash|0" {
		return fmt.Errorf("expected response to be: %s, but actual is: %s", "shamsuddins cash|0", out)
	}
	return nil
}

func theTrustsCostAccountIncreased(ctx context.Context) error {
	out := executeQuery("SELECT account_sat.name || \"|\" || account_sat.balance FROM trust_account_link JOIN account_hub USING(account_hash) JOIN account_sat USING(account_hash) WHERE account_hash = '" + hashString(fmt.Sprintf("%v", ctx.Value(trustCostAccountCtxKey{}))) + "' ORDER BY account_sat.load_time DESC LIMIT 1;")

	if strings.TrimRight(out, "\n") != "shamsuddins cost|1" {
		return fmt.Errorf("expected response to be: %s, but actual is: %s", "shamsuddins cost|1", out)
	}
	return nil
}

func theTrustBuysSomeSharesOfTheCompanyThroughTheBroker(ctx context.Context) (context.Context, error) {
	insertLinkWithTwoHubs("trust_company_link", "shamsuddins", "AZJ")
	insertTrustCompanySat("shamsuddins", "AZJ", 10)
	insertAccount("1", "shamsuddins cash", 0)
	insertAccount("5", "shamsuddins cost", 1)
	ctx = context.WithValue(ctx, trustCostAccountCtxKey{}, "5")
	return context.WithValue(ctx, trustCashAccountCtxKey{}, "1"), nil
}

func theresACompanyInTheDatabase() error {
	insertCompany("AZJ", "Aurizon Holdings", "Transportation")
	return nil
}

func thePersonCanBeSeenInTheTrustMemberList(ctx context.Context) error {
	out := executeQuery("select person.name || \"|\" || trust.name from party_hub person join person_trust_link on (person_hash = person.party_hash) join party_hub trust on (trust_hash = trust.party_hash);")

	if strings.TrimRight(out, "\n") != "melby|shamsuddins" {
		return fmt.Errorf("expected response to be: %s, but actual is: %s", "melbys|shamsuddins", out)
	}
	return nil
}

func thePersonIsAddedToTheTrust(ctx context.Context) error {
	insertLinkWithTwoHubs(
		"person_trust_link",
		fmt.Sprintf("%v", ctx.Value(personCtxKey{})),
		fmt.Sprintf("%v", ctx.Value(trustCtxKey{})),
	)
	return nil
}

func theresAPersonInTheDatabase(ctx context.Context) (context.Context, error) {
	insertParty("melby", "trustee")
	return context.WithValue(ctx, personCtxKey{}, "melby"), nil
}

func theresATrustInTheDatabase(ctx context.Context) (context.Context, error) {
	insertParty("shamsuddins", "trust")
	return context.WithValue(ctx, trustCtxKey{}, "shamsuddins"), nil
}

func theBrokerCanBeSeenInTheTrustBrokerList(ctx context.Context) error {
	out := executeQuery("select broker.name || \"|\" || trust.name from party_hub broker JOIN broker_trust_link on (broker_hash = broker.party_hash) JOIN party_hub trust on (trust_hash = trust.party_hash);")

	if strings.TrimRight(out, "\n") != "cmcmarkets|shamsuddins" {
		return fmt.Errorf("expected response to be: %s, but actual is: %s", "cmcmarkets|shamsuddins", out)
	}
	return nil
}

func theBrokerIsAddedToTheTrust(ctx context.Context) error {
	insertLinkWithTwoHubs(
		"broker_trust_link",
		fmt.Sprintf("%v", ctx.Value(brokerCtxKey{})),
		fmt.Sprintf("%v", ctx.Value(trustCtxKey{})),
	)
	return nil
}

func theresABrokerInTheDatabase(ctx context.Context) (context.Context, error) {
	insertParty("cmcmarkets", "cmcmarkets")
	return context.WithValue(ctx, brokerCtxKey{}, "cmcmarkets"), nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		executeQueries([]string{
			"DELETE FROM account_hub;",
			"DELETE FROM account_sat;",
			"DELETE FROM broker_trust_link;",
			"DELETE FROM company_hub;",
			"DELETE FROM company_sat;",
			"DELETE FROM party_hub;",
			"DELETE FROM party_sat;",
			"DELETE FROM person_trust_link;",
			"DELETE FROM trust_account_link;",
			"DELETE FROM trust_company_link;",
			"DELETE FROM trust_company_sat;",
		})
		return ctx, nil
	})
	ctx.Step(`^the person can be seen in the trust member list$`, thePersonCanBeSeenInTheTrustMemberList)
	ctx.Step(`^the person is added to the trust$`, thePersonIsAddedToTheTrust)
	ctx.Step(`^there\'s a person in the database$`, theresAPersonInTheDatabase)
	ctx.Step(`^there\'s a trust in the database$`, theresATrustInTheDatabase)
	ctx.Step(`^the broker can be seen in the trust broker list$`, theBrokerCanBeSeenInTheTrustBrokerList)
	ctx.Step(`^the broker is added to the trust$`, theBrokerIsAddedToTheTrust)
	ctx.Step(`^there\'s a broker in the database$`, theresABrokerInTheDatabase)
	ctx.Step(`^the trust buys some shares of the company through the broker$`, theTrustBuysSomeSharesOfTheCompanyThroughTheBroker)
	ctx.Step(`^there\'s a company in the database$`, theresACompanyInTheDatabase)
	ctx.Step(`^the company shares show up in the trust portfolio$`, theCompanySharesShowUpInTheTrustPortfolio)
	ctx.Step(`^the trust has some money in the cash account$`, theTrustHasSomeMoneyInTheCashAccount)
	ctx.Step(`^the trust\'s cash account become empty$`, theTrustsCashAccountBecomeEmpty)
	ctx.Step(`^the trust\'s cost account increased$`, theTrustsCostAccountIncreased)
}

func FeatureContext(s *godog.TestSuiteContext) {
	s.BeforeSuite(func() {
		executeQueries([]string{
			"CREATE TABLE IF NOT EXISTS account_hub(account_hash text unique, account_number text, load_time timestamp DEFAULT(STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')));",
			"CREATE TABLE IF NOT EXISTS account_sat(hash text unique, account_hash text, name text, balance int, load_time timestamp DEFAULT(STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')));",
			"CREATE TABLE IF NOT EXISTS broker_trust_link(broker_trust_hash text unique, broker_hash text, trust_hash text, load_time timestamp);",
			"CREATE TABLE IF NOT EXISTS company_hub(company_hash text unique, code text, load_time timestamp DEFAULT(STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')));",
			"CREATE TABLE IF NOT EXISTS company_sat(hash text unique, company_hash text, name text, gics text, load_time timestamp DEFAULT(STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')));",
			"CREATE TABLE IF NOT EXISTS party_hub(party_hash text unique, name text, load_time timestamp DEFAULT(STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')) NOT NULL);",
			"CREATE TABLE IF NOT EXISTS party_sat(hash text unique, party_hash text, party_type text, load_time timestamp DEFAULT(STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')));",
			"CREATE TABLE IF NOT EXISTS person_trust_link(person_trust_hash text unique, person_hash text, trust_hash text, load_time timestamp);",
			"CREATE TABLE IF NOT EXISTS trust_account_link(trust_account_hash text unique, trust_hash text, account_hash text, load_time timestamp);",
			"CREATE TABLE IF NOT EXISTS trust_company_link(trust_company_hash text unique, trust_hash text, company_hash text, load_time timestamp);",
			"CREATE TABLE IF NOT EXISTS trust_company_sat(trust_company_hash text unique, shareQuantity int, load_time timestamp DEFAULT(STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')));",
		})
	})
}
