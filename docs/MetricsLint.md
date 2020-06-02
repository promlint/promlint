
## Common Rules
- A metric should contains `help` text.
- A metric's unit should be one of the `Metric Standard Unit`.
- metric name should not include type, such as `COUNTER`, `GAUGE`, `SUMMARY`, `UNTYPED`, `HISTOGRAM`.
- metric name should not contain ':'.
- metric name should be written in 'snake_case' not 'camelCase'.
- label name should be written in 'snake_case' not 'camelCase'.
- metric name should not contain abbreviated units.

## Rules For Counter
- A counter metric should have `_total` suffix.
- A non-counter metric should not have `_total` suffix.

## Rules For Histogram
- non-histogram metrics should not have "_bucket" suffix`.
- non-histogram and non-summary metrics should not have "_count" suffix
- non-histogram and non-summary metrics should not have "_sum" suffix

## Metric Standard Unit

### Base Units
- amperes
- bytes
- celsius
- grams
- joules
- metres
- seconds
- volts

### Time Units
- seconds

### Temperature Units
- celsius

### Length Units
- meters

### Bytes Units
- bytes

### Energy Units
- joules

### Mass Units
- grams