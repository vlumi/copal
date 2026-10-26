// While the processor opens or runs, disable the form's button and say so.
document.addEventListener("submit", (event) => {
  const form = event.target;
  const button = form.querySelector("button");
  if (form.dataset.busy && button) {
    button.disabled = true;
    button.textContent = form.dataset.busy;
  }
});
