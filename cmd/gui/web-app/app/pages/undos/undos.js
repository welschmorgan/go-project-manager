class Undos {
  constructor() {
    this.columns = ['Name', 'Date']
  }

  init() {
    this.view = document.querySelector('#undo-table-body');
    this.view.innerHTML = 'Loading ...';
    fetch("/api/undos")
      .then(res => res.json())
      .then(undos => {
        this.undos = undos;
        this.renderUndos();
      })
  }

  renderUndos() {
    this.view.innerHTML = '';
    for (const undo of undos) {
      const row = document.createElement('tr');
      for (const col of this.columns) {
        const cell = document.createElement('td');
        cell.textContent = undo[col];
        row.appendChild(cell);
      }
      this.view.appendChild(row);
    }
  }
}