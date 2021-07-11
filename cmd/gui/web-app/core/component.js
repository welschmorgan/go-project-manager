class Component {
  constructor(app, name, html, factory) {
    this.name = name;
    this.app = app;
    this.html = html;
    this.factory = factory;
    this.instance = undefined;
  }
}

