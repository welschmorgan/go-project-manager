class View {
  constructor(options) {
    this.selector = (options && options.selector) || "body";
    this.container = document.querySelector(this.selector);
    this.component = null;
  }

  inject(component) {
    this.innerHTML = component.html;
    this.component = component;
    this.link();
    if (component && component.instance && "init" in component.instance) {
      component.instance.init();
    }
  }

  linkInstructions(html) {
    // {{ fn2() }}
    const instr = /\{\{\s*(.*?)\s*(\(.*?\)|)\}\}/g;
    const instrMatches = instr.exec(html);
    if (instrMatches) {
      const jsLines = [
        "try {",
        "\tconst component = app.cachedComponents['" + this.component.name + "'];",
        "\tconst expr = component.instance." + instrMatches[1] + ";",
        instrMatches[2],
        "} catch (e) {",
        "\tconsole.error('failed to execute '" + instrMatches[1] + "' on component '" + this.name + "', e);",
        "}"
      ];
      html = html.replace(instrMatches[0], jsLines.join(";\n"));
    }
    return html;
  }

  linkEvents(html) {
    // (click)="method()"
    const events = /\([\w]+\)\s*=\s*(["']?[\w]+["']?)/g;
    const eventMatches = events.exec(html);
    
    if (eventMatches) {
      html = html.replace(
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
    return html;    
  }

  link() {
    let html = this.component.html;
    html = this.linkEvents(html);
    html = this.linkInstructions(html);
    this.innerHTML = html;
  }

  set innerHTML(content) {
    this.container.innerHTML = content;
  }

  set innerText(content) {
    this.container.innerText = content;
  }

  get innerHTML() {
    return this.container.innerHTML;
  }

  get innerText() {
    return this.container.innerHTML;
  }
}
