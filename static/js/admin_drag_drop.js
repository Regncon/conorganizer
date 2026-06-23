const dragType = "application/x-conorganizer-dnd";

function closestDragItem(target) {
  return target instanceof Element ? target.closest("[data-dnd-kind][draggable='true']") : null;
}

function closestDropTarget(target) {
  return target instanceof Element ? target.closest("[data-dnd-accept][data-dnd-drop-url-template]") : null;
}

function canDrop(item, target) {
  return Boolean(item && target && item.dataset.dndKind === target.dataset.dndAccept);
}

function dropURL(item, target) {
  const id = item.dataset.dndId || "";
  if (!id) {
    return "";
  }
  return (target.dataset.dndDropUrlTemplate || "").replace("{id}", encodeURIComponent(id));
}

function patchReturnedHTML(html) {
  const doc = new DOMParser().parseFromString(html, "text/html");
  for (const replacement of doc.body.children) {
    if (!replacement.id) {
      continue;
    }
    const current = document.getElementById(replacement.id);
    current?.replaceWith(replacement);
  }
}

document.addEventListener("dragstart", (evt) => {
  const item = closestDragItem(evt.target);
  if (!item || !evt.dataTransfer) {
    return;
  }

  evt.dataTransfer.effectAllowed = "move";
  evt.dataTransfer.setData(dragType, item.dataset.dndId || "");
  item.classList.add("is-dragging");
});

document.addEventListener("dragend", (evt) => {
  const item = closestDragItem(evt.target);
  item?.classList.remove("is-dragging");
  document.querySelectorAll(".is-drag-over").forEach((el) => el.classList.remove("is-drag-over"));
});

document.addEventListener("dragover", (evt) => {
  const item = document.querySelector(".is-dragging");
  const target = closestDropTarget(evt.target);
  if (!canDrop(item, target)) {
    return;
  }

  evt.preventDefault();
  evt.dataTransfer.dropEffect = "move";
  target.classList.add("is-drag-over");
});

document.addEventListener("dragleave", (evt) => {
  const target = closestDropTarget(evt.target);
  if (!target || target.contains(evt.relatedTarget)) {
    return;
  }
  target.classList.remove("is-drag-over");
});

document.addEventListener("drop", async (evt) => {
  const item = document.querySelector(".is-dragging");
  const target = closestDropTarget(evt.target);
  if (!canDrop(item, target)) {
    return;
  }

  evt.preventDefault();
  target.classList.remove("is-drag-over");

  const url = dropURL(item, target);
  if (!url) {
    return;
  }

  const response = await fetch(url, {
    method: "POST",
    headers: {
      Accept: "text/html, text/event-stream, application/json",
      "Conorganizer-Dnd": "true",
      "Datastar-Request": "true",
    },
  });

  if (!response.ok) {
    console.error(`Drag/drop update failed: ${response.status} ${response.statusText}`);
    return;
  }

  const contentType = response.headers.get("Content-Type") || "";
  if (contentType.includes("text/html")) {
    patchReturnedHTML(await response.text());
  }
});
