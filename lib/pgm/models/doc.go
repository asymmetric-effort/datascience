//go:generate python3 ../../../tests/python/generate_fixtures.py --package models --output ../../../tests/fixtures
//go:generate python3 ../../../tests/python/generate_fixtures.py --package markov_network --output ../../../tests/fixtures
//go:generate python3 ../../../tests/python/generate_fixtures.py --package dseparation --output ../../../tests/fixtures

// Package models provides graphical model structures including
// Bayesian networks, Markov networks, and factor graphs.
package models
