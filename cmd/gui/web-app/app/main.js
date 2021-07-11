class App {
  constructor(options) {
    try {
      this.currentComponent = null;
      this.cachedComponents = {};
      this.baseTitle = "[GRLM:UI]";
      this.componentsUrlPrefix = "app/pages/";
      this.viewContainerSelector = options && options.viewContainerSelector;
      this.mainMenuSelector = options && options.mainMenuSelector;
      this.view = new View({ selector: this.viewContainerSelector });
      this.mainMenu = document.querySelector(this.mainMenuSelector);
      this.routes = [
        { route: "home", label: "Home" },
        { route: "projects", label: "Projects" },
        { route: "versions", label: "Versions" },
      ];
      this.createMenu();
      this.navigate("home");
    } catch (e) {
      console.error("failed to initialize app", e);
    }
  }

  async navigate(name) {
    const route = this.routes.find((r) => r.route === name);
    setTitle(`${this.baseTitle} ${route ? route.label : name}`);
    this.view.innerHTML = "";
    return this.createComponent(name);
  }

  async createComponent(name) {
    return new Promise((resolve, reject) => {
      try {
        if (this.isComponentCached(name)) {
          const component = this.getCachedComponent(name);
          this.view.inject(component);
          resolve(component);
        } else {
          let counter = 0;
          this.cachedComponents[name] = new Component(this, name, "");
          const checkAllFiles = (e) => {
            counter++;
            if (counter == 2) {
              this.view.inject(this.cachedComponents[name]);
              resolve(this.cachedComponents[name]);
            }
          };
          this.fetchComponentPart(name, "js").then(
            (response) =>
              this.onScriptFetched(this.cachedComponents[name], response, checkAllFiles),
            reject
          );
          this.fetchComponentPart(name, "html").then(
            (response) =>
              this.onViewFetched(this.cachedComponents[name], response, checkAllFiles),
            reject
          );
        }
      } catch (e) {
        console.error("failed to create component", e);
        reject(e);
      }
    });
  }

  onViewFetched = (component, response, counter) => {
    component.html = response || "";
    counter();
  };

  onScriptFetched = (component, response, counter) => {
    const prevScr = document.querySelector(`app-${component.name}`);
    if (!prevScr) {
      const scr = document.createElement("script");
      scr.type = "text/javascript";
      scr.text = response;
      scr.id = `app-${component.name}`;
      document.body.appendChild(scr);
    }
    // eval(response);
    component.factory = eval(
      component.name.charAt(0).toUpperCase() + component.name.substring(1)
    );
    component.instance = new component.factory();
    counter();
  };

  async fetchComponentPart(name, ext) {
    return fetch(this.componentsUrlPrefix + name + "/" + name + "." + ext).then(
      (response) => response.text()
    );
  }

  createMenuItem(route, label) {
    const li = document.createElement("li");
    const btn = document.createElement("button");
    btn.innerText = label;
    btn.onclick = () => this.navigate(route);
    li.appendChild(btn);
    this.mainMenu.appendChild(li);
  }

  createMenu() {
    for (const route of this.routes) {
      this.createMenuItem(route.route, route.label);
    }
  }

  isComponentCached(name) {
    return !!this.getCachedComponent(name);
  }

  getCachedComponent(name) {
    for (const key in this.cachedComponents) {
      if (key === name && this.cachedComponents[key].name == name) {
        return this.cachedComponents[key];
      }
    }
  }

  message(text) {
    const li = document.createElement("li");
    const p = document.createElement("p");
    p.innerHTML = text;
    li.appendChild(p);
    this.statusBar.appendChild(li);
  }
}
