//go:build unit

package scigo

import (
	"math"
	"strings"
	"testing"
)

// --- Input validation: panics on invalid inputs ---

func TestBlackScholesCallPanicsOnInvalidInputs(t *testing.T) {
	cases := []struct {
		name  string
		S, K  float64
		sigma float64
		want  string
	}{
		{"S=0", 0, 100, 0.2, "S (underlying price) must be positive"},
		{"S<0", -1, 100, 0.2, "S (underlying price) must be positive"},
		{"K=0", 100, 0, 0.2, "K (strike price) must be positive"},
		{"K<0", 100, -1, 0.2, "K (strike price) must be positive"},
		{"sigma<0", 100, 100, -0.1, "sigma (volatility) must be non-negative"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if r == nil {
					t.Fatal("expected panic")
				}
				msg := r.(string)
				if !strings.Contains(msg, tc.want) {
					t.Errorf("panic = %q, want substring %q", msg, tc.want)
				}
			}()
			BlackScholesCall(tc.S, tc.K, 1.0, 0.05, tc.sigma)
		})
	}
}

func TestBlackScholesPutPanicsOnInvalidInputs(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for S=0")
		}
	}()
	BlackScholesPut(0, 100, 1.0, 0.05, 0.2)
}

func TestBlackScholesGreeksPanicsOnInvalidInputs(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for sigma<0")
		}
	}()
	BlackScholesGreeks(100, 100, 1.0, 0.05, -0.1)
}

// --- Zero volatility degenerate cases ---

func TestBlackScholesCallZeroVol(t *testing.T) {
	S, K, T, r := 100.0, 90.0, 1.0, 0.05

	// Forward = S * exp(r*T) = 100 * exp(0.05) ≈ 105.13
	// Call = max(forward - K, 0) * exp(-r*T) = (105.13 - 90) * exp(-0.05) ≈ 14.39
	price := BlackScholesCall(S, K, T, r, 0)
	forward := S * math.Exp(r*T)
	expected := math.Max(forward-K, 0) * math.Exp(-r*T)

	if math.Abs(price-expected) > 1e-10 {
		t.Errorf("Call(sigma=0) = %f, want %f", price, expected)
	}
}

func TestBlackScholesPutZeroVol(t *testing.T) {
	S, K, T, r := 100.0, 110.0, 1.0, 0.05

	// Forward ≈ 105.13, K=110 → put is ITM
	price := BlackScholesPut(S, K, T, r, 0)
	forward := S * math.Exp(r*T)
	expected := math.Max(K-forward, 0) * math.Exp(-r*T)

	if math.Abs(price-expected) > 1e-10 {
		t.Errorf("Put(sigma=0) = %f, want %f", price, expected)
	}
}

func TestBlackScholesCallZeroVolOTM(t *testing.T) {
	// Forward < K → call is worthless
	price := BlackScholesCall(100, 200, 1.0, 0.05, 0)
	if price != 0 {
		t.Errorf("OTM Call(sigma=0) = %f, want 0", price)
	}
}

func TestBlackScholesPutZeroVolOTM(t *testing.T) {
	// Forward > K → put is worthless
	price := BlackScholesPut(100, 50, 1.0, 0.05, 0)
	if price != 0 {
		t.Errorf("OTM Put(sigma=0) = %f, want 0", price)
	}
}

func TestBlackScholesGreeksZeroVol(t *testing.T) {
	// ITM call: delta should be 1
	g := BlackScholesGreeks(100, 90, 1.0, 0.05, 0)
	if g.Delta != 1.0 {
		t.Errorf("ITM Greeks(sigma=0).Delta = %f, want 1.0", g.Delta)
	}
	if g.Gamma != 0 {
		t.Errorf("Greeks(sigma=0).Gamma = %f, want 0", g.Gamma)
	}
	if g.Vega != 0 {
		t.Errorf("Greeks(sigma=0).Vega = %f, want 0", g.Vega)
	}

	// OTM call: delta should be 0
	g2 := BlackScholesGreeks(100, 200, 1.0, 0.05, 0)
	if g2.Delta != 0.0 {
		t.Errorf("OTM Greeks(sigma=0).Delta = %f, want 0.0", g2.Delta)
	}
}

func TestBlackScholesZeroVolPutCallParity(t *testing.T) {
	S, K, T, r := 100.0, 105.0, 1.0, 0.05
	call := BlackScholesCall(S, K, T, r, 0)
	put := BlackScholesPut(S, K, T, r, 0)

	// Put-call parity: C - P = S - K*exp(-rT)
	lhs := call - put
	rhs := S - K*math.Exp(-r*T)

	if math.Abs(lhs-rhs) > 1e-10 {
		t.Errorf("Put-call parity violated at sigma=0: C-P=%f, S-Ke^{-rT}=%f", lhs, rhs)
	}
}

// --- No NaN/Inf propagation ---

func TestBlackScholesNoNaNOrInf(t *testing.T) {
	// Normal case should produce finite values
	price := BlackScholesCall(100, 100, 1, 0.05, 0.2)
	if math.IsNaN(price) || math.IsInf(price, 0) {
		t.Errorf("Call produced non-finite value: %f", price)
	}

	g := BlackScholesGreeks(100, 100, 1, 0.05, 0.2)
	for _, v := range []float64{g.Delta, g.Gamma, g.Theta, g.Vega, g.Rho} {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			t.Errorf("Greeks produced non-finite value: %+v", g)
			break
		}
	}
}

// --- ImpliedVolatility validation ---

func TestImpliedVolatilityRejectsInvalidSK(t *testing.T) {
	_, err := ImpliedVolatility(10, 0, 100, 1, 0.05, "call")
	if err == nil {
		t.Fatal("expected error for S=0")
	}
	_, err = ImpliedVolatility(10, 100, 0, 1, 0.05, "call")
	if err == nil {
		t.Fatal("expected error for K=0")
	}
}
