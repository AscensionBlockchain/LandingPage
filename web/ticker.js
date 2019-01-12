export function Attach($page) {
  $page
    .find(".marquee")
    .eq(0)
    .css("top", 0);

  $page
    .find(".marquee")
    .eq(1)
    .css({
      bottom: 0,
      "border-bottom": "",
      "border-top": "1px solid black"
    })
    .find(".marquee3k")
    .attr("data-reverse", true);

  setTimeout(() => Marquee3k.init({ selector: "marquee3k", speed: 0.5 }));
}

export function getTicker1() {
  return [
    {
      Name: "Elon Musk",
      Change: 23,
      Up: true
    },
    {
      Name: "Annihilation",
      Change: 15,
      Up: false
    },
    {
      Name: "Sandra Bullock",
      Change: 2,
      Up: false
    },
    {
      Name: "Taj Mahal",
      Change: 15,
      Up: false
    },
    {
      Name: "Russell Crowe",
      Change: 76,
      Up: true
    },
    {
      Name: "Yosemite Park",
      Change: 5,
      Up: false
    },
    {
      Name: "Arnold Schwarzenegger",
      Change: 4,
      Up: true
    }
  ];
}

export function getTicker2() {
  return [
    {
      Name: "Aquaman",
      Change: 13,
      Up: true
    },
    {
      Name: "Mission Impossible: 6",
      Change: 4,
      Up: true
    },
    {
      Name: "Bruce Willis",
      Change: 20,
      Up: false
    },
    {
      Name: "Jack Nicholson",
      Change: 121,
      Up: true
    },
    {
      Name: "Bird Box",
      Change: 20,
      Up: true
    },
    {
      Name: "Rowan Atkinson",
      Change: 10,
      Up: false
    },
    {
      Name: "Xiaomi Mi Series",
      Change: 20,
      Up: true
    },
    {
      Name: "Acapulco",
      Change: 121,
      Up: true
    },
    {
      Name: "Alexandria Ocasio-Cortez",
      Change: 55,
      Up: true
    },
    {
      Name: "The Washington Monument",
      Change: 12,
      Up: true
    }
  ];
}
