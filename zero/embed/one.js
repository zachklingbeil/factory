const state = {
	currentFrame: 0,
	totalFrames: 1,
	baseURL: 'http://localhost:1001',
};

function updateBody(html) {
	let oneDiv = document.getElementById('one');
	if (!oneDiv) {
		oneDiv = document.createElement('div');
		oneDiv.id = 'one';
		document.body.appendChild(oneDiv);
	}
	oneDiv.innerHTML = html;
}
async function showFrame(idx) {
	const frameCount = state.totalFrames + 1;
	const safeIdx = ((idx % frameCount) + frameCount) % frameCount;
	const path = `/frame/${safeIdx}`;
	try {
		const response = await fetch(state.baseURL + path);
		if (!response.ok) throw new Error();
		const html = await response.text();
		updateBody(html);
		state.currentFrame = safeIdx;
	} catch {
		updateBody('');
	}
}

document.addEventListener('DOMContentLoaded', async () => {
	try {
		const response = await fetch(state.baseURL + '/frame/0');
		if (!response.ok) throw new Error();
		const html = await response.text();
		updateBody(html);

		const header = response.headers.get('X-Frames');
		if (header) {
			state.totalFrames = parseInt(header, 10);
		}
		state.currentFrame = 0;
	} catch (error) {
		console.error('Error loading initial frame:', error);
		updateBody('');
	}

	document.addEventListener('keydown', (e) => {
		if (e.key === 'q') showFrame(state.currentFrame - 1);
		if (e.key === 'e') showFrame(state.currentFrame + 1);
	});
});

// --- Keybinds ---
function setupContainerKeybinds(containerId, handler) {
	document.addEventListener('DOMContentLoaded', () => {
		const container = document.getElementById(containerId);
		if (!container) return;
		container.tabIndex = 0;
		container.addEventListener('keydown', (event) =>
			handler(event, container)
		);
	});
}

function scrollKeyHandler(event, container) {
	if (event.key === 'w')
		container.scrollBy({ top: -100, behavior: 'smooth' });
	if (event.key === 's') container.scrollBy({ top: 100, behavior: 'smooth' });
}
