class Component {
  constructor(app, name, html, factory) {
    this.name = name;
    this.app = app;
    this.html = html;
    this.factory = factory;
    this.instance = undefined;
  }

  linkTemplateVars() {
    const events = /\([\w]+\)\s*=\s*(["']?[\w]+["']?)/g;
    const instr = /\{\{\s*(.*?)\s*(\(.*?\)|)\}\}/g;
    const instrMatches = instr.exec(this.html);
    const eventMatches = events.exec(this.html);

    if (instrMatches) {
      const jsLines = [
        "try {",
        "\tconst component = app.cachedComponents['" + this.name + "'];",
        "\tconst expr = component.instance." + instrMatches[1] + ";",
        instrMatches[2],
        "} catch (e) {",
        "\tconsole.error('failed to execute '" + instrMatches[1] + "' on component '" + this.name + "', e);",
        "}"
      ];
      this.html = this.html.replace(instrMatches[0], jsLines.join(";\n"));
    }
    if (eventMatches) {
      this.html = this.html.replace(
        eventMatches[0],
        "on" +
          eventMatches[1] +
          "=\"app.cachedComponents['" +
          this.name +
          "'].instance." +
          eventMatches[2] +
          '"'
      );
    }
  }
}

