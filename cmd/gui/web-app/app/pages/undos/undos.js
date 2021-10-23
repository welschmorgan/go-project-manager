class Undos {
  constructor() {
    this.columns = [
      {label: 'Name', field: 'name'}, 
      {label: 'Date', field: 'date'}
    ];
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
      .catch((reason) => console.error(reason));
  }

  renderUndos() {
    this.view.innerHTML = 'Got ' + this.undos.length + ' undos';
    console.log(this.undos);
    for (const undo of this.undos) {
      const row = document.createElement('tr');
      for (const col of this.columns) {
        const cell = document.createElement('td');
        cell.textContent = undo[col.field];
        row.appendChild(cell);
      }
      this.view.appendChild(row);
    }
  }
}