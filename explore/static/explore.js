// Bring a source line into view below the header.
if (location.hash.startsWith("#L")) {
  document.getElementById(location.hash.slice(1))?.scrollIntoView({ block: "center" });
}
