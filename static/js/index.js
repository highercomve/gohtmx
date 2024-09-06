function Load() {
	const elements = document.querySelectorAll(".custom-select label");
	elements.forEach((elm) => {
		elm.addEventListener("click", closeDetails);
	});
}

function closeDetails(event) {
	const details = event.target.closest("details");
	details.removeAttribute("open");
}

window.addEventListener("load", () => {
	Load();
	document.body.addEventListener("htmx:load", function (evt) {
		Load();
	});
});
