// Toggle create form
document.addEventListener('DOMContentLoaded', function() {
	const createBtn = document.querySelector('button[x-data]');
	if (createBtn) {
		createBtn.addEventListener('click', function() {
			const form = document.querySelector('.create-form');
			if (form.style.display === 'none') {
				form.style.display = 'block';
				form.querySelector('input[name="description"]').focus();
			} else {
				form.style.display = 'none';
			}
		});
	}
});

// Handle HTMX events for smooth interactions
document.addEventListener('htmx:afterSwap', function(event) {
	// Re-process Alpine.js directives if needed
	if (window.Alpine) {
		Alpine.scan(event.detail.target);
	}
});

// Auto-hide success messages
document.addEventListener('htmx:afterSwap', function(event) {
	const message = event.detail.target.querySelector('.status-message');
	if (message) {
		setTimeout(() => {
			event.detail.target.remove();
		}, 2000);
	}
});

// Smooth scroll to new habit
document.addEventListener('htmx:afterSwap', function(event) {
	if (event.detail.target.id && event.detail.target.id.startsWith('habit-')) {
		event.detail.target.scrollIntoView({ behavior: 'smooth' });
	}
});

// Set contribution grid date range based on screen width (fluid transition)
function getMonthsToShow() {
	const width = window.innerWidth;
	const minWidth = 320;
	const maxWidth = 1536;
	const minMonths = 3;
	const maxMonths = 12;
	
	// Clamp width to range
	const clampedWidth = Math.max(minWidth, Math.min(width, maxWidth));
	
	// Linear interpolation between min and max months
	const months = minMonths + ((clampedWidth - minWidth) / (maxWidth - minWidth)) * (maxMonths - minMonths);
	
	return Math.round(months);
}

function getDateRange() {
	const today = new Date();
	const monthsBack = getMonthsToShow();
	const fromDate = new Date(today.getFullYear(), today.getMonth() - monthsBack, today.getDate());
	
	return {
		from: fromDate.toISOString().split('T')[0],
		to: today.toISOString().split('T')[0]
	};
}

function loadContributionGrids() {
	const dates = getDateRange();
	
	document.querySelectorAll('[data-habit-id]').forEach(el => {
		const habitID = el.getAttribute('data-habit-id');
		const url = `/api/habits/${habitID}/contribution?from=${dates.from}&to=${dates.to}`;
		
		el.setAttribute('hx-get', url);
		htmx.process(el);
		htmx.trigger(el, 'load');
	});
}

// Load on initial page load
document.addEventListener('DOMContentLoaded', function() {
	setTimeout(loadContributionGrids, 100);
});

// Reload on window resize with debounce
let resizeTimeout;
window.addEventListener('resize', function() {
	clearTimeout(resizeTimeout);
	resizeTimeout = setTimeout(loadContributionGrids, 500);
});
