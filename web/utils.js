// @flow

// global: app CBOR $ nacl smoothScroll ResizeObserver

export var uriRegExp = new RegExp("^[a-z-]+:[^:]+(:[^:]+)?$");

export function chain(...fns: Array<Function>): () => mixed {
  return () => {
    for (var i = 0; i < fns.length; i++) {
      fns[i]();
    }
  };
}

export function makeLinksOpenNewTabs($elem: JQuery) {
  $($elem)
    .find("a")
    .attr("target", "_blank");
}

export function strToNum(str: string): number {
  let n = +str;
  return isNaN(n) ? 0 : n;
}

function _cached_copy<T>(cache: Array<Object>, o: T): T | void {
  for (let i = 0; i < cache.length; i++) {
    if (cache[i].object === o) {
      return cache[i].copy;
    }
  }
}

export function copy<T>(o: T, cache: Array<Object> | void): T {
  if (!o) return o;
  if (!cache) cache = [];

  if (Array.isArray(o)) {
    let c = _cached_copy(cache, o);
    if (c) return c;
    c = [];
    cache.push({ object: o, copy: c });
    for (let i = 0; i < o.length; i++) c.push(copy(o[i], cache));
    return (c: any);
  }
  if (typeof o === "object") {
    let c = _cached_copy(cache, o);
    if (c) return c;
    c = {};
    cache.push({ object: o, copy: c });
    for (var key in o) {
      c[key] = copy(o[key], cache);
    }
    return (c: any);
  } else {
    return o;
  }
}

export function merge(o: Object, s: Object): Object {
  if (s) {
    for (var key in s) {
      o[key] = s[key];
    }
  }
  return o;
}

export function fill(o: Object, s: Object): Object {
  if (s) {
    for (var key in s) {
      if (!(key in o)) o[key] = s[key];
    }
  }
  return o;
}

export function safeCall(fn: () => mixed) {
  try {
    return fn();
  } catch (err) {
    console.error(err);
  }
}

export async function sleep(ms: number) {
  await new Promise((resolve, reject) => setTimeout(resolve, ms));
}

export function scrollTo(elem_selector: any) {
  let $scrollTo = $(elem_selector);
  let $container = $scrollTo.parents(".scrollable").eq(0);

  smoothScroll($scrollTo[0], 500, () => true, $container[0]);
}

// trigger changes when the element indicated by <selector> changes in size
export function onResize(
  selector: string | JQuery,
  callback: Function,
  invokeOnce?: boolean
) {
  // this code is based on github.com/developit/simple-element-resize-detector

  let element = $(selector)[0];

  if (
    $(element)
      .css("position")
      .toString() != "relative"
  ) {
    console.warn(
      "setting 'position: relative;' on element, to observe resize events"
    );
    $(element).css("position", "relative");
  }

  const CSS =
    "position:absolute;left:0;top:-100%;width:100%;height:100%;margin:1px 0 0;border:none;opacity:0;visibility:hidden;pointer-events:none;";

  let frame = document.createElement("iframe");
  frame.style.cssText = CSS;
  frame.onload = () => {
    frame.contentWindow.onresize = () => {
      callback(element);
    };
  };
  element.appendChild(frame);

  if (invokeOnce) callback(element);

  return frame;
}

// function escapeHtml(unsafe) {
//   return unsafe
//     .replace(/&/g, "&amp;")
//     .replace(/</g, "&lt;")
//     .replace(/>/g, "&gt;")
//     .replace(/"/g, "&quot;")
//     .replace(/'/g, "&#039;");
// }
