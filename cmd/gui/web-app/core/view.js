class View {
  constructor(options) {
    this.selector = (options && options.selector) || "body";
    this.container = document.querySelector(this.selector);
    this.component = null;
  }

  inject(component) {
    this.innerHTML = component.html;
    this.component = component;
    if (component && component.instance && "init" in component.instance) {
      component.instance.init();
    }
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
