const API = "http://localhost:8080";

export async function loginUser(email, password) {
    const res = await fetch(`${API}/api/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.error || "Login failed");
    return data;
}

export async function registerUser(username, email, password) {
    const res = await fetch(`${API}/api/register`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, email, password }),
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.error || "Register failed");
    return data;
}

export async function fetchDecks(token) {
    const res = await fetch(`${API}/api/decks`, {
        headers: { Authorization: `Bearer ${token}` },
    });
    if (!res.ok) throw new Error(await res.text());
    return res.json();
}

// Public decks (no auth)
export async function fetchPublicDecks() {
    const res = await fetch(`${API}/api/decks/public`);
    if (!res.ok) throw new Error(await res.text());
    return res.json();
}

export async function fetchDeck(token, id) {
    const res = await fetch(`${API}/api/decks/${id}`, {
        headers: { Authorization: `Bearer ${token}` },
    });
    if (!res.ok) throw new Error(await res.text());
    return res.json();
}

export async function fetchCards(token, deckId) {
    const res = await fetch(`${API}/api/decks/${deckId}/cards`, {
        headers: { Authorization: `Bearer ${token}` },
    });
    if (!res.ok) throw new Error(await res.text());
    return res.json();
}

export async function createDeck(token, title, description) {
    const res = await fetch(`${API}/api/decks`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ title, description }),
    });
    if (!res.ok) throw new Error(await res.text());
    return res.json();
}

export async function createCard(token, deckId, question, answer) {
    const res = await fetch(`${API}/api/decks/${deckId}/cards`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ question, answer }),
    });
    if (!res.ok) throw new Error(await res.text());
    return res.json();
}

export async function updateDeck(token, deckId, payload) {
    const res = await fetch(`http://localhost:8080/api/decks/${deckId}`, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(payload), // ✅ must be an object
    });

    const data = await res.json().catch(() => ({}));
    if (!res.ok) throw new Error(`updateDeck failed: ${res.status} ${JSON.stringify(data)}`);
    return data;
}


export async function fetchPublicDeck(deckId) {
    const res = await fetch(`${API}/api/decks/public/${deckId}`);
    const data = await res.json().catch(() => ({}));
    if (!res.ok) throw new Error(data.error || "Failed to load public deck");
    return data;
}

export async function fetchPublicCards(deckId) {
    const res = await fetch(`${API}/api/decks/public/${deckId}/cards`);
    const data = await res.json().catch(() => ([]));
    if (!res.ok) throw new Error((data && data.error) || "Failed to load public cards");
    return data;
}