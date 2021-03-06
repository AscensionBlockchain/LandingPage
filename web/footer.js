const discordURL = "https://discord.gg/EUdryhF";
const githubURL = "https://github.com/AscensionBlockchain/Ascension";
const telegramURL = "https://t.me/ascensionxch";

export async function Attach($body) {
  let $footer = app.render.Template("footer");

  $footer.find("[data-target-telegram]").click(() => {
    window.open(telegramURL, "_blank");
  });
  $footer.find("[data-target-discord]").click(() => {
    window.open(discordURL, "_blank");
  });
  $footer.find("[data-target-github]").click(() => {
    window.open(githubURL, "_blank");
  });

  $body.find("footer").replaceWith($footer);
}
