class Page {
  constructor(name, html, factory) {
    this.name = name
    this.html = html
    this.factory = factory
    this.instance = undefined;
  }
}

class App {
  constructor() {
    this.view = document.querySelector('#view');
    this.mainMenu = document.querySelector('#main-menu');
    this.currentPage = null;
    this.cachedPages = [];
    this.routes = [
      {route: 'home', label: 'Home'}, 
      {route: 'projects', label: 'Projects'}, 
      {route: 'versions', label: 'Versions'},
    ];
    this.createMenu();
    this.navigate('home');
  }

  createMenuItem(route, label) {
    const li = document.createElement('li');
    const btn = document.createElement('button');
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

  isPageCached(name) {
    return !!this.getCachedPage(name)
  }

  getCachedPage(name) {
    for (const page of this.cachedPages) {
      if (page.name === name) {
        return page;
      }
    }
  }

  navigate(name) {
    this.view.innerHTML = '';
    if (this.isPageCached(name)) {
      const page = this.getCachedPage(name);
      this.view.innerHTML = page.html
      page.instance.init();
    } else {
      return new Promise((resolve, reject) => {
        let counter = 0;
        const page = new Page(name, '');
        const checkAllFiles = (e) => {
          counter++;
          if (counter == 2) {
            this.view.innerHTML = page.html;
            if ('init' in page.instance) {
              page.instance.init();
            }
            this.cachedPages.push(page);
            resolve(name)
          }
        };
        fetch('app/pages/' + name + '/' + name + '.js')
          .then((response) => response.text())
          .then((response) => {
            const scr = document.createElement('script');
            scr.type = "text/javascript";
            scr.text = response;
            document.body.appendChild(scr);
            // eval(response);
            page.factory = eval(name.charAt(0).toUpperCase() + name.substring(1));
            page.instance = new (page.factory)();
            checkAllFiles();
          }, alert);
        fetch('app/pages/' + name + '/' + name + '.html')
          .then((response) => response.text())
          .then((response) => {
            page.html = response;
            checkAllFiles()
          }, alert);
      });
    }
  }

}