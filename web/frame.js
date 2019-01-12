/* global app $ PerfectScrollbar Arrive */

export var currentPage = {};
export const maxRenderTimeout = 10000;

var currentRenderCount = 0;
var renderTimeout;
var stopLoadingAnimation = false;
var scrollbars = [];

export function ShowPage($elem, name, renderCount) {
  if (renderCount != currentRenderCount) {
    return console.warn("ignoring late render for", name);
  }

  if (!app.flags.pageLoaded) app.loading.stopPageLoadingAnimation();

  app.loading.stopLoadingAnimation();

  if (name == "main")
    $("main")
      .find("header")
      .prependTo($elem);

  $("main")
    .eq(0)
    .replaceWith($elem);

  currentPage.$page = $elem;

  resetScrollbars();

  if (name && pageTitles[name]) document.title = pageTitles[name]();
  else document.title = "Ascension Exchange";

  clearTimeout(renderTimeout);

  app.flags.pageLoaded = true;
}

export function Page(name, config) {
  let result;

  console.log(`Page(${name}, ${JSON.stringify(config)})`);
  console.log(` (old) currentRenderCount = ${currentRenderCount}`);
  console.log(` (old) renderTimeout[active] = ${!!renderTimeout}`);

  if (renderTimeout && currentPage.name === name) {
    return console.warn("ignoring re-render while page is still rendering");
  }

  try {
    result = app.page[name].Render(config);
  } catch (err) {
    return console.warn(`failed to switch to page '${name}' because: ${err}`);
  }
  if (!result) return console.warn(`failed to switch to page '${name}'`);

  currentPage.name = name;
  currentPage.opts = config;
  currentPage.subscriptions = [];

  let renderCount = ++currentRenderCount;
  if (renderTimeout) {
    clearTimeout(renderTimeout);
    renderTimeout = 0;
  }
  renderTimeout = setTimeout(() => {
    currentRenderCount++;
    app.flags.renderInProgress = false;
    renderTimeout = undefined;
    alert(`Unable to load ${name} right now. Unknown error encountered.`);
    if (app.flags.pageLoaded) app.router.goBack();
  }, maxRenderTimeout);

  app.flags.renderInProgress = true;

  if (typeof result.then === "function") {
    app.loading.beginLoadingAnimation();
    try {
      result.then($page => {
        app.loading.stopLoadingAnimation();
        ShowPage($page, name, renderCount);
      });
    } catch (err) {
      app.loading.stopLoadingAnimation();
      app.flags.renderInProgress = false;
      return console.warn(
        "failed to (async) switch to page '${name}' because:",
        err
      );
    }
  } else {
    ShowPage(result, name, renderCount);
  }
  return true;
}

export function EnsurePage(name, config) {
  if (name === currentPage.name) return false;
  Page(name, config);
  return true;
}

export function stopLoadingAnimationAfterPageRender() {
  stopLoadingAnimation = true;
}

export function Unsubscribe() {
  console.log(
    "frame: unsubscribing page",
    currentPage.name,
    "from all URIs",
    currentPage.subscriptions
  );
  app.state.unsubscribe(currentPage.subscriptions);
  currentPage.subscriptions = [];
}

export function Reload() {
  Page(currentPage.name, currentPage.opts);
}

export function View(name) {
  if (app.page[currentPage.name].SelectView) {
    console.log("frame: switching page", currentPage.name, "to view", name);
    app.page[currentPage.name].SelectView(currentPage.$page, name).then(() => {
      console.log("frame: view switch complete");
      app.viewbar.ActiveView(currentPage.$page, name);
    });
  } else console.warn(`NYI: app.page.${currentPage.name}.SwitchView(...)`);
}

let pageTitles = {
  main: () => "Ascension",
  about: () => "About · Ascension"
};

let doneWithScrollbarMutationObserver;

function resetScrollbars() {
  Arrive.unbindAllArrive();
  scrollbars.forEach(ps => ps.destroy());
  scrollbars.splice(0, scrollbars.length);

  // create a scrollbar for all .scrollable elements on the page;
  //   if elements are added later, or are modified to be .scrollable,
  //     create the scrollbar whenever that happens
  currentPage.$page[0].arrive(
    ".scrollable",
    { fireOnAttributesModification: true, existing: true },
    function() {
      createScrollbar(this);
    }
  );

  app.utils.onResize(
    currentPage.$page,
    () => {
      scrollbars.forEach(ps => ps.update());
    },
    true
  );
}

function createScrollbar(scrollable) {
  if ($(scrollable).data("scrollbarEnabled")) return;
  console.log("enabling perfect-scrollbar for elem", scrollable);
  scrollbars.push(
    new PerfectScrollbar(scrollable, {
      suppressScrollX: true
    })
  );
  $(scrollable).data("scrollbarEnabled", true);
}
