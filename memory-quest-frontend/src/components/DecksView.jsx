import { Link } from "react-router-dom";
import React from "react";

export default function DecksView({
                                      token,
                                      authOk,
                                      decks,
                                      deckTitle,
                                      setDeckTitle,
                                      deckDescription,
                                      setDeckDescription,
                                      onCreateDeck,
                                      deckIsPublic,
                                      setDeckIsPublic
                                  }) {
    const card = { background: "#2a2a2a", padding: 14, borderRadius: 10 };
    const input = { width: "100%", padding: 12, boxSizing: "border-box" };
    const deckList = Array.isArray(decks) ? decks : [];


    return (
        <div style={{ display: "grid", gap: 12 }}>
            <div style={card}>
                Status: {token ? (authOk ? "✅ Authorized" : "⏳ Checking...") : "🔒 Login required"}
            </div>

            <div style={card}>
                <form onSubmit={onCreateDeck} style={{ display: "grid", gap: 10 }}>
                    <input value={deckTitle} onChange={e => setDeckTitle(e.target.value)} placeholder="Deck title" style={input} disabled={!token} />
                    <input value={deckDescription} onChange={e => setDeckDescription(e.target.value)} placeholder="Description" style={input} disabled={!token} />
                    <label style={{ display: "flex", alignItems: "center", gap: 8, opacity: token ? 1 : 0.6 }}>
                        <input type="checkbox" checked={!!deckIsPublic} onChange={(e) => setDeckIsPublic(e.target.checked)} disabled={!token} />
                        Public deck
                    </label>
                    <button disabled={!token}>Create</button>
                </form>
            </div>

            <div style={card}>
                {deckList.length === 0 ? (
                    "(No decks yet)"
                ) : (
                    <ul>
                        {deckList.map((d, i) => {
                            const deckId = d.id ?? d.ID ?? d.Id ?? i;
                            const title = d.title ?? d.Title ?? "(untitled)";
                            const desc = d.description ?? d.Description ?? "";
                            const isPublic = d.is_public ?? d.IsPublic ?? false;


                            return (
                                <li key={`${deckId}-${i}`}>
                                    <Link to={`/decks/${deckId}`} style={{ color: "inherit" }}>
                                        <b>{title}</b> {isPublic ? <span style={{ fontSize: 12, opacity: 0.75 }}>(public)</span> : <span style={{ fontSize: 12, opacity: 0.75 }}>(private)</span>}
                                    </Link>
                                    {desc ? ` — ${desc}` : ""}
                                </li>
                            );
                        })}
                    </ul>
                )}
            </div>
        </div>
    );
}
