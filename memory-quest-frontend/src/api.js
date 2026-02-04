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

export async function createDeck(token, title, description, is_public) {
    const res = await fetch(`${API}/api/decks`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ title, description, is_public: !!is_public }),
    });
    if (!res.ok) throw new Error(await res.text());
    return res.json();
}

export async function updateDeck(token, deckId, title, description, is_public) {
    const res = await fetch(`${API}/api/decks/${deckId}`, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ title, description, is_public: !!is_public }),
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
