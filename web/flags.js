// @flow

// global: app

var callbacks = {};

export var debug = false;

export var pageLoaded = false;

export var renderInProgress = false;

export function updateFlag(name: string, value: any) {
  if (app.flags[name] === value) return;
  app.flags[name] = value;
  if (callbacks[name]) callbacks[name].forEach(callback => callback(value));
}

export function watchFlag(name: string, callback: Function) {
  if (callbacks[name]) callbacks[name].push(callback);
  else callbacks[name] = [callback];

  callback(app.flags[name]);
}
