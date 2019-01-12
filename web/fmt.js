// @flow

// global: app Hogan $

export var Raw = Hogan.Lambda(function(tmpl, render, data) {
  return this[tmpl.trim()];
});

export var Currency = Hogan.Lambda(function(tmpl, render, data) {
  return app.currency.formatCurrency(+this[tmpl.trim()] || 0);
});

export var Say = Hogan.Lambda(function(tmpl, render, data) {
  return this;
});
