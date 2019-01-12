// @flow

// globals: app Marquee3k

export function Render() {
  let $page = app.render.Template("main", {
    TitleTop: "Ascension opens a",
    TitleMiddle: "new and exciting",
    TitleBottom: "world of investments",
    Subtitle:
      "Trade on the fortunes of Elon Musk or the next Disney blockbuster.",
    Ticker1: { Marquee: app.ticker.getTicker1() },
    Ticker2: { Marquee: app.ticker.getTicker2() }
  });

  $page.find("#whitepaper button").click(function() {
    location.href = "/assets/whitepaper.html";
  });

  app.ticker.Attach($page);

  setTimeout(() => app.logo.Attach($page));

  return $page;
}
