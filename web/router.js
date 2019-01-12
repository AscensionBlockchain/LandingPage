/* global app StateMan $ */

export function SetupRoutes() {
  console.log("router: setting up routes");
  app.router.stateman = new StateMan();
  app.router.stateman
    .state("/", {
      enter: () => app.frame.Page("main")
      // leave: () => console.log("")
    })
    .state("feature", {
      url: "^/feature/:itemURI(.+)",
      canEnter: option => {
        return app.frame.canNavigateTo("feature", option.param);
      },
      enter: option => {
        console.log("loading feature-page with uri", option.param);
        option.param.itemURI = app.utils.createURI(
          "feature",
          option.param.itemURI
        );
        return app.frame.Page("feature", option.param);
      },
      leave: app.frame.Leave
    })
    // .state("**", { enter: app.page.stack.Render })

    .start({
      html5: true
    });
}

var checkpoints = [];

export function go(target, options) {
  if (options) {
    if (typeof options === "string") {
      throw new Error(
        "please supply an options object, not a string like: " + options
      );
    }
    if (options.checkpoint) {
      addCheckpoint(true);
    }
  }
  app.router.stateman.go(target, options);
}

export function nav(target) {
  app.router.stateman.nav(target);
}

export function goBack(toCheckpoint) {
  if (checkpoints.length && toCheckpoint) {
    let lastCheckpoint = checkpoints[checkpoints.length - 1];
    console.log("router.goBack: returning to checkpoint", lastCheckpoint);
    app.router.stateman.go(lastCheckpoint);
    lastCheckpoint = null;
    return true;
  } else if (app.router.stateman.previous) {
    console.log(
      "router.goBack: returning to immediately preceding view:",
      app.router.stateman.previous
    );
    app.router.stateman.go(app.router.stateman.previous.name);
    return true;
  }
  return false;
}

export function mustGoBack(toCheckpoint) {
  return goBack(toCheckpoint) || nav("/");
}

export function switchRouteBack() {
  return history.replaceState(null, "", app.router.stateman.previous.path);
}

export function addCheckpoint(beforeNav) {
  if (beforeNav) {
    checkpoints.push(app.router.stateman.current.name);
  } else {
    checkpoints.push(app.router.stateman.previous.name);
  }
}
