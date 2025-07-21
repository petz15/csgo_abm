# Distribution Analysis Summary: CSFNormalDistribution_std_custom_skew

## Key Parameter Effects (with minOutput=0, maxOutput=5)

### 1. **stdDevFactor** - Controls Distribution Spread
- **stdDevFactor = 2.0**: Very wide spread (stdDev=1.807), values cluster at extremes (0.0 and 5.0)
- **stdDevFactor = 4.0**: Moderate spread (stdDev=1.200), normal bell curve
- **stdDevFactor = 8.0**: Tight spread (stdDev=0.632), values concentrated around mean

### 2. **CSF Competition Parameter (r)** - Affects Determinism
- **r = 0.5**: Less deterministic outcomes, similar distribution shapes
- **r = 1.0**: Standard competition level
- **r = 2.0**: More deterministic, stronger teams have more predictable advantages

### 3. **Team Strength Ratios (x:y)** - Shifts Mean Position
- **Equal (100:100)**: Mean ≈ 2.5 (center of 0-5 range)
- **2:1 advantage (200:100)**: Mean ≈ 3.3 (CSF prob = 0.667)
- **3:1 advantage (300:100)**: Mean ≈ 3.8 (CSF prob = 0.750)
- **1:2 disadvantage (100:200)**: Mean ≈ 1.7 (CSF prob = 0.333)

### 4. **Skewness** - Dramatically Changes Shape
- **skew = 0.0**: Normal bell curve distribution
- **skew = -0.5**: Moderate left skew, 54% of values in [0.0-0.5] range
- **skew = +0.5**: Moderate right skew, 54% of values in [4.5-5.0] range
- **skew = -1.0**: Strong left skew, 75% of values in [0.0-0.5] range
- **skew = +1.0**: Strong right skew, 75% of values in [4.5-5.0] range

### 5. **Combined Effects** - Real-World Scenarios
- **Strong team + right skew**: Amplifies advantage (mean 4.5, 70% in top bin)
- **Weak team + left skew**: Amplifies disadvantage (mean 0.5, 70% in bottom bin)

## Practical Implications for CS:GO Simulation

### Equipment/Weapon Selection (0-5 scale)
- **Economy balanced teams**: Use stdDevFactor=4, skew=0 for realistic variety
- **Eco vs Full-buy rounds**: Use higher team ratios (3:1 or 4:1) with slight skew
- **Pistol rounds**: Use equal ratios with moderate positive skew (slight bias to better equipment)

### Round Outcome Probability (0-5 scale could represent round win chance)
- **Conservative modeling**: Use stdDevFactor=6-8 for tighter distributions
- **Volatile matches**: Use stdDevFactor=2-3 for more extreme outcomes
- **Momentum effects**: Use positive skew for winning teams, negative for losing teams

### Recommended Parameter Combinations
1. **Standard balanced**: x=100, y=100, r=1.0, stdDevFactor=4.0, skew=0.0
2. **Economic advantage**: x=200, y=100, r=1.0, stdDevFactor=4.0, skew=0.2
3. **Momentum swing**: x=150, y=100, r=1.0, stdDevFactor=3.0, skew=0.5
4. **Tight conservative**: x=100, y=100, r=2.0, stdDevFactor=6.0, skew=0.0
5. **High variance upset potential**: x=100, y=100, r=0.5, stdDevFactor=2.0, skew=0.0

The skewness parameter is particularly powerful - even small values (±0.5) create noticeable asymmetry, while values of ±1.0 create extreme clustering at the boundaries.
