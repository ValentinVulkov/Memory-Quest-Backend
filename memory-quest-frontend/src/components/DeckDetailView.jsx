import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { fetchDeck, fetchCards, createCard, updateDeck, fetchPublicDeck, fetchPublicCards } from "../api";

import CardViewerModal from "./CardViewModal";
import React from "react";

export default function DeckDetailView({ token, userId }) {
    const { deckId } = useParams();
    //const isGuest = !token;
    const [deck, setDeck] = useState(null);
    const [cards, setCards] = useState([]);
    const [msg, setMsg] = useState("");

    const [question, setQuestion] = useState("");
    const [answer, setAnswer] = useState("");

    // Viewing mode (read cards one-by-one)
    const [viewerOpen, setViewerOpen] = useState(false);
    const [viewerStartIndex, setViewerStartIndex] = useState(0);

    async function load() {
        const deckData = await fetchDeck(token, deckId);
        const cardsData = await fetchCards(token, deckId);
        setDeck(deckData);
        setCards(cardsData);
    }

    useEffect(() => {
        async function load() {
            try {
                setMsg("");
                if (token) {
                    const d = await fetchDeck(token, deckId);
                    const c = await fetchCards(token, deckId);
                    setDeck(d);
                    setCards(Array.isArray(c) ? c : []);
                } else {
                    const d = await fetchPublicDeck(deckId);
                    const c = await fetchPublicCards(deckId);
                    setDeck(d);
                    setCards(Array.isArray(c) ? c : []);
                }
            } catch (e) {
                setMsg(e?.message || String(e));
            }
        }
        load();
    }, [token, deckId]);

    async function onCreate(e) {
        e.preventDefault();

        const q = question.trim();
        const a = answer.trim();
        if (!q || !a) {
            setMsg("Question and answer are required.");
            return;
        }

        setMsg("Creating card...");

        try {
            await createCard(token, deckId, q, a);
            setQuestion("");
            setAnswer("");
            setMsg("✅ Card created");
            await load();
        } catch (e) {
            setMsg("Create card failed: " + (e?.message || String(e)));
        }
    }
    const isGuest = !token;
    const ownerId = deck?.user_id;
    // simple dark styles
    const card = { background: "#1e1e1e", padding: 14, borderRadius: 10, border: "1px solid #2a2a2a" };
    const input = { width: "100%", padding: 12, boxSizing: "border-box", borderRadius: 8, border: "1px solid #333", background: "#111", color: "#fff" };
    const btn = { padding: "10px 12px", cursor: "pointer", background: "#2a2a2a", color: "#fff", border: "1px solid #333", borderRadius: 8 };

    const deckTitle = deck ? (deck.title ?? deck.Title ?? `Deck #${deckId}`) : `Deck #${deckId}`;

    const deckUserId = deck?.user_id ?? deck?.UserID ?? null;


    const isOwner =
        !!token &&
        deck?.user_id != null &&
        userId != null &&
        Number(deck.user_id) === Number(userId);
    const isPublic = !!(deck?.is_public ?? deck?.IsPublic ?? false);


    return (

        <div style={{ display: "grid", gap: 12 }}>
            <div style={{fontSize: 12, opacity: 0.7}}>
                guest={String(isGuest)} ownerId={String(ownerId)} userId={String(userId)} isOwner={String(isOwner)}
            </div>
            <div style={card}>
                <Link to="/decks" style={{ color: "#fff" }}>← Back to decks</Link>
                <div style={{ marginTop: 10, fontWeight: 800 }}>{deckTitle}</div>

                <div style={{ marginTop: 8, display: "flex", alignItems: "center", gap: 10, flexWrap: "wrap" }}>
                    <span style={{ fontSize: 13, opacity: 0.85 }}>
                        Visibility: <b>{isPublic ? "Public" : "Private"}</b>
                    </span>

                    {isOwner && (
                        <label style={{ display: "flex", alignItems: "center", gap: 8, fontSize: 13, opacity: 0.95 }}>
                            <input
                                type="checkbox"
                                checked={isPublic}
                                onChange={async (e) => {
                                    const next = e.target.checked;

                                    try {
                                        const title = deck?.title ?? deck?.Title ?? "";
                                        const description = deck?.description ?? deck?.Description ?? "";

                                        await updateDeck(token, deck.id ?? deck.ID, {
                                            title,
                                            description,
                                            is_public: next,
                                        });

                                        // update UI locally
                                        setDeck({
                                            ...deck,
                                            is_public: next,
                                            IsPublic: next,
                                        });

                                        setMsg(next ? "✅ Deck is now public" : "✅ Deck is now private");
                                    } catch (err) {
                                        setMsg("Update visibility failed: " + (err?.message || String(err)));
                                    }
                                }}

                            />
                            Make public
                        </label>
                    )}
                </div>

                <div style={{ marginTop: 12, display: "flex", gap: 10, flexWrap: "wrap", marginBottom: -15 }}>
                    <button
                        type="button"
                        onClick={() => { setViewerStartIndex(0); setViewerOpen(true); }}
                        disabled={!cards.length}
                        style={{
                            padding: "24px 56px",
                            fontSize: "18px",
                            fontWeight: "bold",
                            borderRadius: "10px",
                            background: "#2563eb",
                            color: "white",
                            border: "none",
                            cursor: "pointer",
                            boxShadow: "0 4px 10px rgba(0,0,0,0.15)",
                            marginBottom: "15px"
                        }}
                    >
                        View
                    </button>

                </div>
            </div>
            {!isGuest && (
            <div style={card}>
                <div style={{ fontWeight: 800, marginBottom: 10 }}>Create card</div>
                <form onSubmit={onCreate} style={{ display: "grid", gap: 10 }}>
                    <input
                        value={question}
                        onChange={(e) => setQuestion(e.target.value)}
                        placeholder="Question"
                        style={input}
                    />
                    <input
                        value={answer}
                        onChange={(e) => setAnswer(e.target.value)}
                        placeholder="Answer"
                        style={input}
                    />
                    <button type="submit" style={btn}>Create</button>
                </form>
            </div>
                )}
            <div style={card}>
                <div style={{ fontWeight: 800, marginBottom: 10 }}>Cards</div>

                {cards.length === 0 ? (
                    <div style={{ opacity: 0.7 }}>(No cards yet)</div>
                ) : (
                    <ul style={{ margin: 0, paddingLeft: 18 }}>
                        {cards.map((c, i) => {
                            const id = c.id ?? c.ID ?? i;
                            const q = c.question ?? c.Question ?? "";
                            return (
                                <li
                                    key={`${id}-${i}`}
                                    style={{ marginBottom: 10, cursor: "pointer" }}
                                    onClick={() => {
                                        setViewerStartIndex(i);
                                        setViewerOpen(true);
                                    }}
                                    title="Click to open reading mode from here"
                                >
                                    <div><b>Q:</b> {q}</div>
                                    <div style={{ opacity: 0.9, marginBottom: 40 }}></div>
                                </li>
                            );
                        })}
                    </ul>
                )}
            </div>

            <CardViewerModal
                open={viewerOpen}
                cards={cards}
                startIndex={viewerStartIndex}
                onClose={() => setViewerOpen(false)}
            />

            {msg && <pre style={{ ...card, whiteSpace: "pre-wrap", margin: 0 }}>{msg}</pre>}
        </div>
    );
}
