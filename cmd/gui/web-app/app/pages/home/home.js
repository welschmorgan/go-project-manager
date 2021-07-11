class Home {
  constructor() {
  }

  init() {
    this.view = document.querySelector('#workspace-infos');
    this.workspace = undefined;
    this.view.innerHTML = '';

    fetch('/api/workspace')
      .then(res => res.json())
      .then(workspace => {
        this.workspace = workspace;
        this.createField('Name');
        this.createField('Path');
        this.createField('Author', workspace.Author.name);
      });
  }

  createField(name, value) {
    const fieldset = document.createElement("fieldset");
    
    const label = document.createElement("label");
    label.textContent = name;
    
    const input = document.createElement("input");
    input.type = 'text';
    input.readOnly = true;
    input.value = value || this.workspace[name];
    
    fieldset.appendChild(label);
    fieldset.appendChild(input);
    this.view.appendChild(fieldset);
    return fieldset;
  };
}