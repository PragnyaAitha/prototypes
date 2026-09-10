package main

import "testing"

func TestShardPlanCreatesThreePartitions(t *testing.T) {
	plan := NewShardPlan()

	if len(plan.Partitions) != 3 {
		t.Fatalf("expected 3 partitions, got %d", len(plan.Partitions))
	}

	for _, shard := range plan.Partitions {
		if shard.Database == "" {
			t.Fatalf("expected each shard to have a database name")
		}
	}
}

func TestShardPlanContainsRequiredTables(t *testing.T) {
	plan := NewShardPlan()

	for _, shard := range plan.Partitions {
		if !contains(shard.Tables, "users") || !contains(shard.Tables, "profiles") || !contains(shard.Tables, "user_profiles") {
			t.Fatalf("shard %s should have users, profiles, and user_profiles tables", shard.Database)
		}
	}
}

func contains(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}
