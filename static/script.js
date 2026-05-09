const FACTS_23 = [
	{ text: "23 IS PRIME. SO ARE 2 AND 3.", cite: "NUMBER THEORY" },
	{ text: "EARTH'S AXIS IS TILTED ABOUT 23.5 DEGREES.", cite: "ASTRONOMY" },
	{ text: "HUMAN DNA IS ORGANIZED INTO 23 PAIRS OF CHROMOSOMES.", cite: "CELL BIOLOGY" },
	{ text: "BIRTHDAY PARADOX: WITH 23 PEOPLE IN A ROOM, THERE'S A >50% CHANCE TWO SHARE A BIRTHDAY.", cite: "PROBABILITY" },
	{ text: "JULIUS CAESAR WAS REPORTEDLY STABBED 23 TIMES.", cite: "SUETONIUS" },
	{ text: "MICHAEL JORDAN WORE 23.", cite: "BASKETBALL" },
	{ text: "23 IS THE SMALLEST PRIME WHOSE DIGITS ARE BOTH PRIME.", cite: "NUMBER THEORY" },
	{ text: "THE LATIN ALPHABET HAS 23 LETTERS — NO J, U, OR W.", cite: "ETYMOLOGY" },
];

const HISTORY_KEY = "numerology.history.v1";
const HISTORY_CAP = 50;
const DEFAULT_API_BASE = "/api";

const GROUP_VARS = [
	"g0", "g1", "g2", "g3", "g4", "g5", "g6", "g7",
	"g8", "g9", "g10", "g11", "g12", "g13", "g14", "g15",
];
function groupColor(idx) {
	const k = GROUP_VARS[idx % GROUP_VARS.length];
	return { fg: `var(--${k}-fg)`, bg: `var(--${k}-bg)` };
}

function loadHistory() {
	try {
		const raw = localStorage.getItem(HISTORY_KEY);
		if (!raw) return [];
		const v = JSON.parse(raw);
		return Array.isArray(v) ? v : [];
	} catch (_) { return []; }
}
function saveHistory(h) {
	try { localStorage.setItem(HISTORY_KEY, JSON.stringify(h.slice(0, HISTORY_CAP))); } catch (_) { }
}
function pad2(n) { return String(n).padStart(2, "0"); }
function nowStamp() {
	const d = new Date();
	return `${d.getFullYear()}.${pad2(d.getMonth() + 1)}.${pad2(d.getDate())} ${pad2(d.getHours())}:${pad2(d.getMinutes())}`;
}

// Locale-aware default date string. US → MMDDYYYY, else DDMMYYYY.
function defaultDateDigits() {
	try {
		const parts = new Intl.DateTimeFormat(undefined, {
			year: "numeric", month: "2-digit", day: "2-digit",
		}).formatToParts(new Date());
		const map = {};
		for (const p of parts) {
			if (p.type === "year" || p.type === "month" || p.type === "day") map[p.type] = p.value;
		}
		const order = parts
			.filter(p => ["year", "month", "day"].includes(p.type))
			.map(p => p.type);
		const mIdx = order.indexOf("month");
		const dIdx = order.indexOf("day");
		if (mIdx < dIdx) return `${map.month}${map.day}${map.year}`;
		return `${map.day}${map.month}${map.year}`;
	} catch (_) {
		const d = new Date();
		return `${pad2(d.getDate())}${pad2(d.getMonth() + 1)}${d.getFullYear()}`;
	}
}

async function callNumerology({ baseUrl, target, input }) {
	const url = `${baseUrl.replace(/\/+$/, "")}/${encodeURIComponent(target)}/${encodeURIComponent(input)}?format=json`;
	let res;
	try {
		res = await fetch(url, { method: "GET", headers: { Accept: "application/json" } });
	} catch (e) {
		throw new Error(`NETWORK BLACKOUT: ${e.message}`);
	}
	if (!res.ok) {
		let detail = "";
		try { detail = await res.text(); } catch (_) { }
		if (res.status === 404) throw new Error(`NO PROOF FOUND FOR THESE DIGITS. ${detail || ""}`.trim());
		if (res.status === 400) throw new Error(`BAD REQUEST. ${detail || ""}`.trim());
		throw new Error(`UNEXPECTED STATUS ${res.status}. ${detail || ""}`.trim());
	}
	return await res.json();
}

function tokenizeExpression(expr) {
	if (typeof expr !== "string") return [];
	const out = [];
	const s = expr
		.replace(/×/g, "*")
		.replace(/·/g, "*")
		.replace(/÷/g, "/")
		.replace(/−/g, "-");
	let i = 0;
	while (i < s.length) {
		const c = s[i];
		if (/\s/.test(c)) { i++; continue; }
		if (/[0-9]/.test(c)) {
			let j = i;
			while (j < s.length && /[0-9]/.test(s[j])) j++;
			const raw = s.slice(i, j);
			out.push({ type: "num", value: raw, raw, digits: raw.split("") });
			i = j;
			continue;
		}
		if ("+-*/()".includes(c)) {
			out.push({ type: "op", value: c, raw: c });
			i++;
			continue;
		}
		out.push({ type: "op", value: c, raw: c });
		i++;
	}
	return out;
}

function buildResultColoring(tokens, input) {
	const inputDigits = (input || "").split("").filter(c => /[0-9]/.test(c));
	const used = new Array(inputDigits.length).fill(false);
	let allMatched = true;

	const tagged = tokens.map(t => {
		if (t.type !== "num") return { ...t };
		const digitColors = t.digits.map(d => {
			for (let i = 0; i < inputDigits.length; i++) {
				if (!used[i] && inputDigits[i] === d) {
					used[i] = true;
					return i;
				}
			}
			allMatched = false;
			return null;
		});
		return { ...t, digitColors };
	});

	const aligned = allMatched && used.every(Boolean);
	const perChar = inputDigits.map((_, i) => i);
	return { tokens: tagged, perChar, aligned };
}

const state = {
	history: loadHistory(),
	factIdx: Math.floor(Math.random() * FACTS_23.length),
	result: null,
	error: null,
	busy: false,
	apiBase: DEFAULT_API_BASE,
};

const $ = (sel) => document.querySelector(sel);
function clear(node) { while (node.firstChild) node.removeChild(node.firstChild); }
function el(tag, attrs, ...children) {
	const n = document.createElement(tag);
	if (attrs) {
		for (const k in attrs) {
			const v = attrs[k];
			if (v == null || v === false) continue;
			if (k === "class") n.className = v;
			else if (k === "style" && typeof v === "object") Object.assign(n.style, v);
			else if (k.startsWith("on") && typeof v === "function") n.addEventListener(k.slice(2), v);
			else if (k in n) n[k] = v;
			else n.setAttribute(k, v);
		}
	}
	for (const c of children) {
		if (c == null || c === false) continue;
		n.appendChild(typeof c === "string" || typeof c === "number" ? document.createTextNode(String(c)) : c);
	}
	return n;
}

function renderFact() {
	const f = FACTS_23[state.factIdx];
	$("#fact-body").textContent = `“${f.text}”`;
	$("#fact-cite").textContent = `— ${f.cite}`;
}


function renderError() {
	const slot = $("#error-slot");
	clear(slot);
	if (!state.error) return;
	const card = el("div", { class: "card error", role: "alert" },
		el("p", { class: "err-title" }, "REQUEST FAILED. TRY DIFFERENT DIGITS."),
		el("div", { class: "err-msg" }, state.error),
	);
	slot.appendChild(card);
}

function renderResult() {
	const slot = $("#result-slot");
	clear(slot);
	if (!state.result) return;
	const r = state.result;
	const { input, value, tokens, perCharGroup, aligned } = r;

	const inputLine = el("div", { class: "res-line input-line", "aria-label": "input digits, colored" },
		el("span", { class: "lbl" }, "INPUT »"),
	);
	const chars = (input || "").split("");
	for (let i = 0; i < chars.length; i++) {
		const g = perCharGroup ? perCharGroup[i] : i;
		const { fg, bg } = groupColor(g ?? 0);
		inputLine.appendChild(el("span", { class: "tok", style: { color: fg, background: bg } }, chars[i]));
	}

	const exprLine = el("div", { class: "res-line", "aria-label": "expression" },
		el("span", { class: "lbl" }, "PROOF »"),
	);
	for (const t of tokens) {
		if (t.type === "num") {
			for (let i = 0; i < t.digits.length; i++) {
				const idx = t.digitColors ? t.digitColors[i] : null;
				if (idx == null) {
					exprLine.appendChild(el("span", { class: "tok", style: { color: "#000", background: "#eee" } }, t.digits[i]));
				} else {
					const { fg, bg } = groupColor(idx);
					exprLine.appendChild(el("span", { class: "tok", style: { color: fg, background: bg } }, t.digits[i]));
				}
			}
		} else {
			const op = t.value === "*" ? "×" : t.value === "/" ? "÷" : t.value === "-" ? "−" : t.value;
			exprLine.appendChild(el("span", { class: "op" }, op));
		}
	}

	const eqLine = el("div", { class: "res-line" },
		el("span", { class: "equals blink" }, "="),
		el("span", { class: "final blink" }, String(value)),
	);

	const card = el("div", { class: "card evidence", role: "region", "aria-label": "proof result" },
		el("div", { class: "result-headline" }, "◢◤ proof complete ◥◣"),
		inputLine,
		exprLine,
		eqLine,
	);

	if (!aligned) {
		card.appendChild(el("div", { class: "align-warn" }, "(digit alignment imperfect — coloring by position.)"));
	}

	slot.appendChild(card);
}

function renderBook() {
	const histSlot = $("#history-slot");
	clear(histSlot);
	$("#burn-evidence").style.display = state.history.length > 0 ? "" : "none";

	if (state.history.length === 0) {
		histSlot.appendChild(el("div", { class: "empty-page" },
			el("span", { class: "quill", "aria-hidden": "true" }, "🪶"),
			"no proofs yet.",
		));
		return;
	}

	state.history.forEach((e, i) => {
		histSlot.appendChild(el("div", { class: "ev-row" },
			el("span", { class: "ev-num" }, `#${state.history.length - i}`),
			el("span", { class: "ev-expr" },
				el("span", { class: "ev-input" }, e.input),
				el("span", { class: "arrow" }, "→"),
				el("span", null, e.expression),
				el("span", { class: "arrow" }, "="),
				el("b", null, String(e.target)),
			),
			el("span", { class: "ev-meta" }, e.when),
		));
	});
}

async function submit() {
	const input = $("#digits").value;
	const targetRaw = $("#target").value;
	state.error = null;
	state.busy = true;
	$("#btn-reveal").disabled = true;
	$("#btn-reveal").textContent = "COMPUTING…";
	renderError();
	try {
		const trgNum = parseInt(targetRaw, 10);
		if (Number.isNaN(trgNum)) throw new Error("TARGET MUST BE AN INTEGER.");
		if (!input || input.length === 0) throw new Error("DIGITS FIELD IS EMPTY.");
		const data = await callNumerology({ baseUrl: state.apiBase, target: trgNum, input });
		const tokens = tokenizeExpression(data.expression || "");
		const colored = buildResultColoring(tokens, data.input || input);
		const r = {
			input: data.input || input,
			target: data.target ?? trgNum,
			expression: data.expression,
			value: data.result,
			tokens: colored.tokens,
			perCharGroup: colored.perChar,
			aligned: colored.aligned,
			when: nowStamp(),
			id: `${Date.now()}-${Math.random().toString(36).slice(2, 7)}`,
		};
		state.result = r;
		state.history = [{
			id: r.id, input: r.input, target: r.target,
			expression: r.expression, value: r.value, when: r.when,
		}, ...state.history].slice(0, HISTORY_CAP);
		saveHistory(state.history);

		renderResult();
		renderBook();
	} catch (e) {
		state.error = e.message || String(e);
		state.result = null;
		renderError();
		renderResult();
	} finally {
		state.busy = false;
		$("#btn-reveal").disabled = false;
		$("#btn-reveal").textContent = "CRUNCH";
	}
}

function clearHistory() {
	if (confirm("Clear all history? This cannot be undone.")) {
		state.history = [];
		saveHistory(state.history);
		renderBook();
	}
}

function renderVisitorCounter() {
	const node = $("#visitor-counter");
	if (!node) return;
	const multiplier = 100 + Math.floor(Math.random() * 49900);
	const count = 23 * multiplier;
	const digits = String(count).padStart(8, "0").split("");
	clear(node);
	for (const d of digits) {
		node.appendChild(el("span", { class: "digit7" }, d));
	}
}

function rand(min, max) { return min + Math.random() * (max - min); }

function spawnOneDolphin(host) {
	const vw = window.innerWidth;
	const vh = window.innerHeight;
	const size = rand(28, 56);
	const offscreen = size * 3;

	function pointOnEdge(edge) {
		switch (edge) {
			case "left": return { x: -offscreen, y: rand(-0.1 * vh, 1.0 * vh) };
			case "right": return { x: vw + offscreen, y: rand(-0.1 * vh, 1.0 * vh) };
			case "top": return { x: rand(-0.1 * vw, 1.0 * vw), y: -offscreen };
			case "bottom": return { x: rand(-0.1 * vw, 1.0 * vw), y: vh + offscreen };
		}
	}
	const edges = ["left", "right", "top", "bottom"];
	const startEdge = edges[Math.floor(Math.random() * edges.length)];
	const exitCandidates = edges.filter(e => e !== startEdge);
	const endEdge = exitCandidates[Math.floor(Math.random() * exitCandidates.length)];
	const { x: startX, y: startY } = pointOnEdge(startEdge);
	const { x: endX, y: endY } = pointOnEdge(endEdge);
	const duration = rand(5500, 11000);

	const travelAngleDeg = Math.atan2(endY - startY, endX - startX) * 180 / Math.PI;
	const rightward = endX >= startX;
	const tiltDeg = rightward ? travelAngleDeg : travelAngleDeg - 180;
	const flip = rightward ? -1 : 1;

	const d = document.createElement("div");
	d.className = "dolphin";
	d.textContent = "🐬";
	d.style.fontSize = `${size}px`;
	const transform = (x, y) => `translate(${x}px, ${y}px) rotate(${tiltDeg}deg) scaleX(${flip})`;
	d.style.transform = transform(startX, startY);
	host.appendChild(d);

	const anim = d.animate(
		[
			{ transform: transform(startX, startY) },
			{ transform: transform(endX, endY) },
		],
		{ duration, easing: "linear", fill: "forwards" }
	);
	anim.onfinish = () => d.remove();
}

const DOLPHIN_CAP = 23;
let dolphinsPerSpawn = 0;

function spawnDolphin() {
	const host = document.getElementById("dolphins");
	if (!host) return;
	dolphinsPerSpawn = Math.min(dolphinsPerSpawn + 1, DOLPHIN_CAP);
	for (let i = 0; i < dolphinsPerSpawn; i++) {
		setTimeout(() => spawnOneDolphin(host), i * rand(80, 260));
	}
}

function scheduleDolphin() {
	const delay = rand(90, 28000);
	setTimeout(() => {
		if (!document.hidden) spawnDolphin();
		scheduleDolphin();
	}, delay);
}

document.addEventListener("DOMContentLoaded", () => {
	$("#digits").value = defaultDateDigits();

	renderFact();
	renderBook();
	renderVisitorCounter();

	let dolphinsStarted = false;
	const startDolphinsOnce = () => {
		if (dolphinsStarted) return;
		if (matchMedia("(prefers-reduced-motion: reduce)").matches) return;
		dolphinsStarted = true;
		spawnDolphin();
		scheduleDolphin();
		scheduleDolphin();
	};

	$("#proof-form").addEventListener("submit", (e) => {
		e.preventDefault();
		startDolphinsOnce();
		submit();
	});
	$("#digits").addEventListener("input", (e) => {
		e.target.value = e.target.value.replace(/\D+/g, "");
	});
	$("#target").addEventListener("input", (e) => {
		e.target.value = e.target.value.replace(/[^0-9-]/g, "");
	});
	$("#burn-evidence").addEventListener("click", (e) => { e.preventDefault(); clearHistory(); });
});
