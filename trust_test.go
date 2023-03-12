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

type personCtxKey struct{}
type trustCtxKey struct{}
type brokerCtxKey struct{}

func hash_string(text string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(text)))
}

func execute_query(query string) string {
	return execute_queries([]string{query})
}

func execute_queries(queries []string) string {
	var sqls string

	for _, sql := range queries {
		sqls = sqls + sql
	}
	var out bytes.Buffer
	cmd := exec.Command("sqlite3", path, sqls)
	cmd.Stdout = &out
	cmd.Run()
	return out.String()
}

func insert_party(name string, party_type string, fields string) {
	execute_queries([]string{
		"INSERT INTO party_hub values('" + hash_string(name) + "', '" + name + "', CURRENT_TIMESTAMP);",
		"INSERT INTO party_sat values('" + hash_string(fields) + "', '" + hash_string(name) + "', '" + party_type + "', CURRENT_TIMESTAMP);",
	})
}

func insert_broker_trust(broker string, trust string) {
	execute_query("insert into broker_trust_link values('" + hash_string(hash_string(broker)+"|"+hash_string(trust)) + "', '" + hash_string(broker) + "', '" + hash_string(trust) + "', CURRENT_TIMESTAMP);")
}

func insert_person_trust(person string, trust string) {
	execute_query("insert into person_trust_link values('" + hash_string(hash_string(person)+"|"+hash_string(trust)) + "', '" + hash_string(person) + "', '" + hash_string(trust) + "', CURRENT_TIMESTAMP);")
}

func thePersonCanBeSeenInTheTrustMemberList(ctx context.Context) error {
	out := execute_query("select person.name || \"|\" || trust.name from party_hub person join person_trust_link on (person_hash = person.party_hash) join party_hub trust on (trust_hash = trust.party_hash);")

	if strings.TrimRight(out, "\n") != "melby|shamsuddins" {
		return fmt.Errorf("expected response to be: %s, but actual is: %s", "melbys|shamsuddins", out)
	}
	return nil
}

func thePersonIsAddedToTheTrust(ctx context.Context) error {
	insert_person_trust(
		fmt.Sprintf("%v", ctx.Value(personCtxKey{})),
		fmt.Sprintf("%v", ctx.Value(trustCtxKey{})),
	)
	return nil
}

func theresAPersonInTheDatabase(ctx context.Context) (context.Context, error) {
	insert_party("melby", "trustee", "melby|trustee")
	return context.WithValue(ctx, personCtxKey{}, "melby"), nil
}

func theresATrustInTheDatabase(ctx context.Context) (context.Context, error) {
	insert_party("shamsuddins", "trust", "shamsuddins|trust")
	return context.WithValue(ctx, trustCtxKey{}, "shamsuddins"), nil
}

func theBrokerCanBeSeenInTheTrustBrokerList(ctx context.Context) error {
	out := execute_query("select broker.name || \"|\" || trust.name from party_hub broker join broker_trust_link on (broker_hash = broker.party_hash) join party_hub trust on (trust_hash = trust.party_hash);")

	if strings.TrimRight(out, "\n") != "cmcmarkets|shamsuddins" {
		return fmt.Errorf("expected response to be: %s, but actual is: %s", "cmcmarkets|shamsuddins", out)
	}
	return nil
}

func theBrokerIsAddedToTheTrust(ctx context.Context) error {
	insert_broker_trust(
		fmt.Sprintf("%v", ctx.Value(brokerCtxKey{})),
		fmt.Sprintf("%v", ctx.Value(trustCtxKey{})),
	)
	return nil
}

func theresABrokerInTheDatabase(ctx context.Context) (context.Context, error) {
	insert_party("cmcmarkets", "cmcmarkets", "cmcmarkets|broker")
	return context.WithValue(ctx, brokerCtxKey{}, "cmcmarkets"), nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		execute_queries([]string{
			"DELETE FROM party_hub;",
			"DELETE FROM party_sat;",
			"DELETE FROM person_trust_link;",
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
}

func FeatureContext(s *godog.TestSuiteContext) {
	s.BeforeSuite(func() {
		execute_queries([]string{
			"CREATE TABLE IF NOT EXISTS party_hub(party_hash text unique, name text, load_time timestamp);",
			"CREATE TABLE IF NOT EXISTS party_sat(hash text unique, party_hash text, party_type text, load_time timestamp);",
			"CREATE TABLE IF NOT EXISTS person_trust_link(person_trust_hash text unique, person_hash text, trust_hash text, load_time timestamp);",
			"CREATE TABLE IF NOT EXISTS broker_trust_link(broker_trust_hash text unique, broker_hash text, trust_hash text, load_time timestamp);",
		})
	})
}
