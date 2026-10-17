// Bring a source line into view below the header.
if (location.hash.startsWith("#L")) {
  document.getElementById(location.hash.slice(1))?.scrollIntoView({ block: "center" });
}

// Pan and zoom the diagrams: drag to move, ctrl or cmd and the wheel, or
// a pinch, to zoom around the pointer, and buttons to zoom and fit.
for (const frame of document.querySelectorAll(".diagram")) {
  const svg = frame.querySelector("svg");
  if (!svg) continue;
  const whole = svg.viewBox.baseVal;
  const full = { x: whole.x, y: whole.y, w: whole.width, h: whole.height };
  let view = { ...full };
  svg.removeAttribute("width");
  svg.removeAttribute("height");
  svg.setAttribute("preserveAspectRatio", "xMidYMin meet");
  // natural shows the diagram at its own size, or smaller to fit the
  // frame's width, from the top: a tall one in a frame of its own height.
  const natural = () => {
    const scale = Math.min(1, frame.clientWidth / full.w);
    const height = Math.min(window.innerHeight * 0.75, full.h * scale);
    svg.style.height = height + "px";
    const w = frame.clientWidth / scale;
    view = { x: full.x - (w - full.w) / 2, y: full.y, w, h: height / scale };
  };
  const show = () => svg.setAttribute("viewBox", `${view.x} ${view.y} ${view.w} ${view.h}`);
  // point is a pointer's position in the diagram's own coordinates.
  const point = (event) => {
    const p = new DOMPoint(event.clientX, event.clientY).matrixTransform(svg.getScreenCTM().inverse());
    return { x: p.x, y: p.y };
  };
  const zoom = (factor, at) => {
    const w = Math.min(full.w * 4, Math.max(full.w / 20, view.w * factor));
    const scale = w / view.w;
    view = { x: at.x - (at.x - view.x) * scale, y: at.y - (at.y - view.y) * scale, w, h: view.h * scale };
    show();
  };
  const center = () => ({ x: view.x + view.w / 2, y: view.y + view.h / 2 });

  const tools = document.createElement("div");
  tools.className = "tools";
  for (const [label, title, act] of [
    ["+", "Zoom in", () => zoom(1 / 1.4, center())],
    ["−", "Zoom out", () => zoom(1.4, center())],
    ["Fit", "Show the whole diagram", () => { view = { ...full }; show(); }],
    ["1:1", "Show the diagram at its own size", () => { natural(); show(); }],
  ]) {
    const button = document.createElement("button");
    button.type = "button";
    button.textContent = label;
    button.title = title;
    button.addEventListener("click", act);
    tools.append(button);
  }
  const hint = document.createElement("span");
  hint.textContent = "drag to move · ctrl+wheel or pinch to zoom";
  tools.append(hint);
  frame.before(tools);

  frame.addEventListener("wheel", (event) => {
    if (!event.ctrlKey && !event.metaKey) return;
    event.preventDefault();
    zoom(Math.exp(Math.max(-0.4, Math.min(0.4, event.deltaY * 0.01))), point(event));
  }, { passive: false });

  let drag = null;
  frame.addEventListener("pointerdown", (event) => {
    if (event.button !== 0) return;
    drag = { from: point(event), moved: false, id: event.pointerId };
  });
  frame.addEventListener("pointermove", (event) => {
    if (!drag) return;
    const at = point(event);
    const dx = at.x - drag.from.x, dy = at.y - drag.from.y;
    if (!drag.moved && Math.hypot(dx, dy) < 3 * (view.w / frame.clientWidth)) return;
    if (!drag.moved) frame.setPointerCapture(drag.id);
    drag.moved = true;
    view.x -= dx;
    view.y -= dy;
    show();
  });
  const end = () => { if (drag?.moved) frame.dataset.dragged = "true"; drag = null; };
  frame.addEventListener("pointerup", end);
  frame.addEventListener("pointercancel", end);
  // A drag that ends on a link does not follow it.
  frame.addEventListener("click", (event) => {
    if (frame.dataset.dragged) {
      event.preventDefault();
      delete frame.dataset.dragged;
    }
  }, true);

  natural();
  window.addEventListener("resize", () => { natural(); show(); });
  show();
}
