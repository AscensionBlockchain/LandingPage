// @flow

export var unit = {
  cubit: 1.0,
  cube: 1.0 * 1000,
  doublecube: 1.0 * 1000 * 1000,
  triplecube: 1.0 * 1000 * 1000 * 1000,
  tetracube: 1.0 * 1000 * 1000 * 1000 * 1000,
  pentacube: 1.0 * 1000 * 1000 * 1000 * 1000 * 1000
};

export function formatCurrency(price: number, precision?: number): string {
  precision = precision || 2;
  let denomination = "cubit";
  let approximatePrice = "0";
  if (price > unit.pentacube) {
    approximatePrice = (price / unit.pentacube).toFixed(precision);
    denomination = "pentacube";
  } else if (price > unit.tetracube) {
    approximatePrice = (price / unit.tetracube).toFixed(precision);
    denomination = "tetracube";
  } else if (price > unit.triplecube) {
    approximatePrice = (price / unit.triplecube).toFixed(precision);
    denomination = "triplecube";
  } else if (price > unit.doublecube) {
    approximatePrice = (price / unit.doublecube).toFixed(precision);
    denomination = "doublecube";
  } else if (price > unit.cube) {
    approximatePrice = (price / unit.cube).toFixed(precision);
    denomination = "cube";
  }
  return `${approximatePrice} ${denomination}${price == 1 ? "" : "s"}`;
}

export function estimateTokenPrice(tokens: number, priceLevel: number): number {
  let p = priceLevel;
  let n = tokens;
  const c = 0.2;
  let total = c * Math.log(n) * (Math.pow(n, 2) / 2) + n * p;
  return isNaN(total) ? 0 : Math.floor(total);
}
