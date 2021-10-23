class Loader {
  constructor(container) {
    // store parent container
    this.container = container;

    // create loading element
    this.loadingElement = document.createElement('div');
    this.loadingElement.id = this.container.id + '-loader';
    this.loadingElement.style.display = 'none';
    this.loadingElement.innerHTML = '<p>Loading ...</p>';
    this.container.appendChild(this.loadingElement);

    this.show();
  }

  show() {
    this.loadingElement.style.display = 'initial';
    for (let child of this.loadingElement.children) {
      child.style.display = 'none';
    }
  }

  hide() {
    this.loadingElement.style.display = 'none';
    for (let child of this.loadingElement.children) {
      child.style.display = 'initial';
    }
  }
}