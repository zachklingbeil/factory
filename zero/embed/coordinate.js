function updateBody(html) {
	let oneDiv = document.getElementById('one');
	if (!oneDiv) {
		oneDiv = document.createElement('div');
		oneDiv.id = 'one';
		document.body.appendChild(oneDiv);
	}
	oneDiv.innerHTML = html;
}

class CoordinatePlane {
	constructor() {
		this.fetchAndRender();
	}

	async fetchAndRender() {
		try {
			const response = await fetch('/api/test');
			if (!response.ok) throw new Error();
			const data = await response.json();
			this.nRows = Math.max(...data.map((coord) => coord.y)) + 1;
			const html = this.render(data);
			updateBody(html);
		} catch (err) {
			console.error('Failed to load test.json:', err);
			updateBody('');
		}
	}

	createCoordinate({ x, y, z }) {
		const axisType = x < 0 ? 'negative' : x > 0 ? 'positive' : 'label';
		return `
            <div class="coordinate ${axisType}">
                <div>${z.peer}</div>
                <div>${z.time}</div>
                <div>${z.value}</div>
            </div>
        `;
	}

	createRow(rowIndex, coordinates) {
		const negativeCoords = coordinates
			.filter((coord) => coord.y === rowIndex && coord.x < 0)
			.map((coord) => this.createCoordinate(coord))
			.join('');
		const positiveCoords = coordinates
			.filter((coord) => coord.y === rowIndex && coord.x > 0)
			.map((coord) => this.createCoordinate(coord))
			.join('');
		const labelCoordinate = coordinates.find(
			(coord) => coord.y === rowIndex && coord.x === 0
		);
		const label = labelCoordinate ? rowIndex : rowIndex;

		return `
            <div class="row">
                <div class="axis left">
                    <div class="coordinate-grid">${negativeCoords}</div>
                </div>
                <div class="label">${label}</div>
                <div class="axis right">
                    <div class="coordinate-grid">${positiveCoords}</div>
                </div>
            </div>
        `;
	}

	render(coordinates) {
		let html = `<div class="coordinate-plane" id="coordinate-plane">`;
		for (let row = 0; row < this.nRows; row++) {
			html += this.createRow(row, coordinates);
		}
		html += `</div>`;
		return html;
	}
}

document.addEventListener('DOMContentLoaded', () => {
	new CoordinatePlane();
});
