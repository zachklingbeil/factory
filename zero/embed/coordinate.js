class CoordinatePlane {
	constructor(container) {
		this.container = container;
	}

	initFromJson(coordinates) {
		this.nRows = Math.max(...coordinates.map((coord) => coord.y)) + 1;
		this.render(coordinates);
	}

	createCoordinate({ x, y, z }) {
		const axisType = x < 0 ? 'negative' : x > 0 ? 'positive' : 'label';
		const coordinate = document.createElement('div');
		coordinate.className = `coordinate ${axisType}`;
		coordinate.innerHTML = `
            <div>${z.peer}</div>
            <div>${z.time}</div>
            <div>${z.value}</div>
        `;
		return coordinate;
	}

	createRow(rowIndex, coordinates) {
		const row = document.createElement('div');
		row.className = 'row';

		const negativeAxis = document.createElement('div');
		negativeAxis.className = 'axis left';
		const negativeGrid = document.createElement('div');
		negativeGrid.className = 'coordinate-grid';
		coordinates
			.filter((coord) => coord.y === rowIndex && coord.x < 0)
			.forEach((coord) =>
				negativeGrid.appendChild(this.createCoordinate(coord))
			);
		negativeAxis.appendChild(negativeGrid);

		const labelAxis = document.createElement('div');
		labelAxis.className = 'label';
		const labelCoordinate = coordinates.find(
			(coord) => coord.y === rowIndex && coord.x === 0
		);
		labelAxis.textContent = labelCoordinate ? rowIndex : rowIndex;

		const positiveAxis = document.createElement('div');
		positiveAxis.className = 'axis right';
		const positiveGrid = document.createElement('div');
		positiveGrid.className = 'coordinate-grid';
		coordinates
			.filter((coord) => coord.y === rowIndex && coord.x > 0)
			.forEach((coord) =>
				positiveGrid.appendChild(this.createCoordinate(coord))
			);
		positiveAxis.appendChild(positiveGrid);

		row.appendChild(negativeAxis);
		row.appendChild(labelAxis);
		row.appendChild(positiveAxis);

		return row;
	}

	render(coordinates) {
		this.container.innerHTML = '';
		for (let row = 0; row < this.nRows; row++) {
			this.container.appendChild(this.createRow(row, coordinates));
		}
	}
}

document.addEventListener('DOMContentLoaded', () => {
	fetch('/api/test')
		.then((r) => r.json())
		.then((data) => {
			const plane = new CoordinatePlane(document.getElementById('one'));
			plane.initFromJson(data);
		})
		.catch((err) => console.error('Failed to load test.json:', err));
});
