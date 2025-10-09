let sections = [];

window.onload = function () {
  sections = document.querySelectorAll(".program-pulje-section[id]");
  window.addEventListener("scroll", navHighlighter);
  console.log('Scroll highlighter script loaded');
}

function navHighlighter() {
  let scrollY = window.pageYOffset;

  sections.forEach(current => {
    const sectionHeight = current.offsetHeight;
    const sectionTop = current.offsetTop - 200;
    const sectionId = current.getAttribute("id");
    const elementQuery = ".scrollnav-button a[href*=" + sectionId + "]"

    if (sectionId === "FredagKveld") {
      console.log(sectionId, scrollY, sectionTop);
    }

    if (
      scrollY > sectionTop &&
      scrollY <= sectionTop + sectionHeight
    ) {
      document.querySelector(elementQuery).classList.add("is-active");
    } else {
      document.querySelector(elementQuery).classList.remove("is-active");
    }
  });
}
