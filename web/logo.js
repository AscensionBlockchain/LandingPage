// @flow

export function Attach($page: JQuery) {
  let $path = $page.find("header .logo-wrapper #svg-logo-path");
  let stdFill = $path.css("fill").toString();
  $page
    .find("header .logo-wrapper")
    .on("mouseenter", () => {
      $path.css("fill", "#ca731a");
    })
    .on("mouseleave", () => {
      $path.css("fill", stdFill);
    });
}
