import { useEffect, useState } from "react";
import React from "react";
import TopBar from "./components/TopBar";
import DecksView from "./components/DecksView";
import AuthModal from "./components/AuthModal";
import { loginUser, registerUser, fetchDecks, fetchPublicDecks, createDeck } from "./api";
import DeckDetailView from "./components/DeckDetailView";
import { Routes, Route, Navigate } from "react-router-dom";

export default function App() {
    const [token, setToken] = useState(localStorage.getItem("token") || "");
    const [authOk, setAuthOk] = useState(false);
    const [decks, setDecks] = useState([]);
    const [publicDecks, setPublicDecks] = useState([]);
    const [msg, setMsg] = useState("");

    const [authOpen, setAuthOpen] = useState(false);
    const [mode, setMode] = useState("login");

    const [username, setUsername] = useState("");
    const [email, setEmail] = useState("v@v.com");
    const [password, setPassword] = useState("123456");
    const [showPassword, setShowPassword] = useState(false);

    const [deckTitle, setDeckTitle] = useState("");
    const [deckDescription, setDeckDescription] = useState("");
    const [guest, setGuest] = useState(false);
    const [userId, setUserId] = useState(undefined);


    useEffect(() => {
        async function loadProfile() {
            if (!token) {
                setUserId(undefined);
                return;
            }

            try {
                const res = await fetch("http://localhost:8080/api/profile", {
                    headers: { Authorization: `Bearer ${token}` },
                });
                const data = await res.json();
                setUserId(data.user_id);
            } catch {
                setUserId(undefined);
            }
        }

        loadProfile();
    }, [token]);


    function continueAsGuest() {
        setGuest(true);
        setAuthOpen(false);
        setMsg("");
    }

    async function loadDecks(t) {
        try {
            const d = await fetchDecks(t);
            setDecks(Array.isArray(d) ? d : []);
            setAuthOk(true);
        } catch (e) {
            setAuthOk(false);
            setDecks([]);
            setMsg("Deck load failed: " + (e?.message || String(e)));
        }
    }

    async function loadPublicDecks() {
        try {
            const d = await fetchPublicDecks();
            setPublicDecks(Array.isArray(d) ? d : []);
        } catch (e) {
            // Public decks are non-critical; keep UI alive
            setPublicDecks([]);
        }
    }

    useEffect(() => {
        if (!token) return;
        setMsg("");
        loadDecks(token);
    }, [token]);

    // Load public decks even if user isn't logged in
    useEffect(() => {
        loadPublicDecks();
    }, []);

    async function login(e) {
        e.preventDefault();
        setMsg("");
        try {
            const data = await loginUser(email, password);
            const t = data.token;
            if (!t) throw new Error("No token returned");
            localStorage.setItem("token", t);
            setToken(t);
            setAuthOpen(false);
            setAuthOk(false);
            setGuest(false);
        } catch (e) {
            setMsg("Login failed: " + (e?.message || String(e)));
        }
    }

    async function register(e) {
        e.preventDefault();
        setMsg("");
        try {
            await registerUser(username, email, password);
            setMode("login");
            setMsg("✅ Registered. Now login.");
        } catch (e) {
            setMsg("Register failed: " + (e?.message || String(e)));
        }
    }

    async function addDeck(e) {
        e.preventDefault();
        if (!token) return;

        setMsg("");
        try {
            await createDeck(token, deckTitle, deckDescription);
            setDeckTitle("");
            setDeckDescription("");
            await loadDecks(token);
            setMsg("✅ Deck created");
        } catch (e) {
            setMsg("Create deck failed: " + (e?.message || String(e)));
        }
    }

    return (
        <div style={{ maxWidth: 900, margin: "40px auto" }}>
            <TopBar
                token={token}
                onOpenAuth={() => setAuthOpen(true)}
                onLogout={() => {
                    localStorage.removeItem("token");
                    setToken("");
                    setGuest(true); // optional: go back to guest browsing
                }}
            />
            <Routes>
                <Route path="/" element={<Navigate to="/decks" replace />} />
                <Route
                    path="/decks"
                    element={
                        <DecksView
                            token={token}
                            authOk={authOk}
                            decks={decks}
                            publicDecks={publicDecks}
                            deckTitle={deckTitle}
                            setDeckTitle={setDeckTitle}
                            deckDescription={deckDescription}
                            setDeckDescription={setDeckDescription}
                            onCreateDeck={addDeck}
                        />
                    }
                />
                <Route path="/decks/:deckId" element={<DeckDetailView token={token} userId={userId} />} />
            </Routes>

            <AuthModal
                open={authOpen}
                token={token}
                mode={mode}
                setMode={setMode}
                username={username}
                setUsername={setUsername}
                email={email}
                setEmail={setEmail}
                password={password}
                setPassword={setPassword}
                showPassword={showPassword}
                setShowPassword={setShowPassword}
                onLogin={login}
                onRegister={register}
                onClose={() => setAuthOpen(false)}
                onGuest={continueAsGuest}
            />

            {msg && (
                <pre style={{ marginTop: 12, background: "#000", color: "#fff", padding: 12, borderRadius: 8 }}>
          {msg}
        </pre>
            )}
        </div>
    );
}
